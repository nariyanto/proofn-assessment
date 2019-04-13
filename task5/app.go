package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"proofn/task5/client"
	"proofn/task5/config"
	"proofn/task5/dao"
	"proofn/task5/service"

	"github.com/dimiro1/health"
	"github.com/dimiro1/health/db"
	"github.com/dimiro1/health/url"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
)

var orderService = service.Order{}

func AllOrdersEndpoint(c echo.Context) error {
	// User ID from path `users/:id`
	orders, err := orderService.GetOrders()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, orders)
}

func main() {
	log.Println("Starting server initialization")

	//Get our config from the file
	var config = config.Config{}
	config.Read()

	//Server params
	var credential = client.Credential{
		Token:          config.Vault.Credential.Token,
		RoleID:         config.Vault.Credential.RoleID,
		SecretID:       config.Vault.Credential.SecretID,
		ServiceAccount: config.Vault.Credential.ServiceAccount,
	}

	var vault = client.Vault{
		Host:           config.Vault.Host,
		Port:           config.Vault.Port,
		Scheme:         config.Vault.Scheme,
		Authentication: config.Vault.Authentication,
		Role:           config.Vault.Role,
		Mount:          config.Vault.Mount,
		Namespace:      config.Vault.Namespace,
		Credential:     credential,
	}

	//Init it
	log.Println("Starting vault initialization")
	err := vault.Initialize()
	if err != nil {
		log.Fatal(err)
	}

	//Make sure we got a DB role
	log.Println("Starting DB initialization")
	if len(config.Vault.Database.Role) == 0 {
		log.Fatal("Could not get DB role from config.")
	}

	//See if we need to go get dyanmic DB creds
	if len(config.Database.Username) == 0 && len(config.Database.Password) == 0 {
		log.Printf("DB role: %s", config.Vault.Database.Role)
		secret, err := vault.GetSecret(fmt.Sprintf("%s/creds/%s", config.Vault.Database.Mount, config.Vault.Database.Role))
		if err != nil {
			log.Fatal(err)
		}
		//Update our configuration with the dynamic creds
		config.Database.Username = secret.Data["username"].(string)
		config.Database.Password = secret.Data["password"].(string)
		//Start our Goroutine Renewal for the DB creds
		go vault.RenewSecret(secret)
	}

	//DAO config
	var orderDao = dao.Order{
		Host:     config.Database.Host,
		Port:     config.Database.Port,
		Database: config.Database.Name,
		User:     config.Database.Username,
		Password: config.Database.Password,
	}

	//Check our DAO connection
	err = orderDao.Connect()
	if err != nil {
		log.Fatal(err)
	}

	//Get our TLS cert from Vault
	cert, err := vault.GetCertificate(fmt.Sprintf("%s/issue/%s", config.Vault.Pki.Mount, config.Vault.Pki.Role), config.Vault.Pki.CN)
	if err != nil {
		log.Fatal(err)
	}

	public := cert.Data["certificate"].(string)
	private := cert.Data["private_key"].(string)
	publicBytes := []byte(public)
	privateBytes := []byte(private)

	x509, err := tls.X509KeyPair(publicBytes, privateBytes)
	if err != nil {
		log.Fatal(err)
	}

	//Create service
	orderService.Vault = &vault
	orderService.Dao = &orderDao
	orderService.Encyrption.Key = config.Vault.Transit.Key
	orderService.Encyrption.Mount = config.Vault.Transit.Mount

	//Router
	r := echo.New()

	//API Routes
	r.GET("/api/orders", AllOrdersEndpoint)
	// r.POST("/api/orders", CreateOrderEndpoint)
	// r.DELETE("/api/orders", DeleteOrdersEndpoint)

	//Health Check Routes
	h := health.NewHandler()
	conn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", config.Database.Username, config.Database.Password, config.Database.Name, config.Database.Host)
	database, _ := sql.Open("postgres", conn)
	pg := db.NewPostgreSQLChecker(database)
	h.AddChecker("Postgres", pg)
	h.AddChecker("Vault", url.NewChecker(fmt.Sprintf("%s://%s:%s/v1/sys/health?perfstandbyok=true", config.Vault.Scheme, config.Vault.Host, config.Vault.Port)))

	r.Logger.Fatal(r.Start(":" + config.Server.Port))

	//CORS
	r.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	//Server config - https
	go func() {
		tlsConfig := &tls.Config{Certificates: []tls.Certificate{x509}}
		s := &http.Server{
			Addr:      ":8443",
			Handler:   r,
			TLSConfig: tlsConfig,
		}
		log.Println(fmt.Sprintf("Server is now accepting https requests on port 8443"))
		if err := s.ListenAndServeTLS("", ""); err != nil {
			log.Fatal(err)
		}
	}()

	//Catch SIGINT AND SIGTERM to gracefully tear down tokens and secrets
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	sig := <-gracefulStop
	fmt.Printf("caught sig: %+v", sig)
	vault.Close()
	os.Exit(0)

}
