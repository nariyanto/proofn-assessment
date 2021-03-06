package app

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"proofn/task5/client"
	"proofn/task5/config"
	"proofn/task5/dao"
	"proofn/task5/helper"
	"proofn/task5/models"
	"proofn/task5/service"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type Data struct {
	AccessToken string      `json:"access_token"`
	User        models.User `json:"user"`
}

var authService = service.User{}

func initAuth(vault client.Vault) {
	var config = config.Config{}
	config.Read()

	//DAO config
	var authDao = dao.User{
		Host:     config.Database.Host,
		Port:     config.Database.Port,
		Database: config.Database.Name,
		User:     config.Database.Username,
		Password: config.Database.Password,
	}

	//Check our DAO connection
	err := authDao.Connect()
	if err != nil {
		log.Fatal(err)
	}

	//Create service
	authService.Vault = &vault
	authService.Dao = &authDao
	authService.Encyrption.Key = config.Vault.Transit.Key
	authService.Encyrption.Mount = config.Vault.Transit.Mount
}

func Signup(c echo.Context, vault client.Vault) error {
	initAuth(vault)

	var user models.User
	var result helper.BaseResponse

	var bodyBytes []byte
	if c.Request().Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request().Body)
	}
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
		result = helper.CreateErrorResponse(http.StatusBadRequest, "Invalid request payload.", nil)
		return c.JSONPretty(result.Status, result, "  ")
	}

	if len(user.TokenVerification) <= 0 || len(user.Email) <= 0 {
		result = helper.CreateErrorResponse(http.StatusBadRequest, "Please fill all required field.", nil)
		return c.JSONPretty(result.Status, result, "  ")
	}

	_, err := authService.GetUsersByEmail(user)
	if err == nil {
		result = helper.CreateErrorResponse(http.StatusForbidden, "Email is already exist.", nil)
		return c.JSONPretty(result.Status, result, "  ")
	}

	if len(user.Password) < 8 {
		result = helper.CreateErrorResponse(http.StatusForbidden, "Password is less than 8 characters.", nil)
		return c.JSONPretty(result.Status, result, "  ")
	}

	json.Unmarshal(bodyBytes, &user)
	hashPassword, _ := helper.HashPassword(user.Password)

	newUser := models.User{
		Name:              user.Name,
		Email:             user.Email,
		Password:          hashPassword,
		Token:             "",
		TokenVerification: helper.RandStringBytes(10),
		Status:            0,
	}

	user, err = authService.CreateUser(newUser)
	if err != nil {
		result = helper.CreateErrorResponse(http.StatusInternalServerError, "Something Wrong", err)
		return c.JSONPretty(result.Status, result, "  ")
	}

	result = helper.CreateSuccessResponse(user, "success")
	return c.JSONPretty(result.Status, result, "  ")
}

func Login(c echo.Context, vault client.Vault) error {
	initAuth(vault)

	var user models.User
	var result helper.BaseResponse

	var bodyBytes []byte
	if c.Request().Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request().Body)
	}
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
		result = helper.CreateErrorResponse(http.StatusBadRequest, "Invalid request payload.", nil)
		return c.JSONPretty(result.Status, result, "  ")
	}

	if len(user.Email) <= 0 || len(user.Password) <= 0 {
		result = helper.CreateErrorResponse(http.StatusBadRequest, "Please fill all required field.", nil)
		return c.JSONPretty(result.Status, result, "  ")
	}

	duser, err := authService.GetUsersByEmail(user)
	if err != nil {
		result = helper.CreateErrorResponse(http.StatusBadRequest, "Email not found.", nil)
		return c.JSONPretty(result.Status, result, "  ")
	}

	if len(user.Password) < 8 {
		result = helper.CreateErrorResponse(http.StatusBadRequest, "Password is less than 8 characters.", nil)
		return c.JSONPretty(result.Status, result, "  ")
	}

	if !helper.CheckPasswordHash(user.Password, duser.Password) {
		result = helper.CreateErrorResponse(http.StatusUnauthorized, "Invalid password", nil)

		return c.JSONPretty(result.Status, result, "  ")
	}

	if duser.Status != 1 {
		result = helper.CreateErrorResponse(http.StatusUnauthorized, "User not yet verified, please verify the account first", nil)

		return c.JSONPretty(result.Status, result, "  ")
	}

	claims := helper.JwtCustomClaims{
		ID:       duser.ID,
		Name:     duser.Name,
		Email:    duser.Email,
		Password: duser.Password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 8760).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := os.Getenv("JWT_SECRET")
	t, err := token.SignedString([]byte(secretKey))

	if err != nil {
		result = helper.CreateErrorResponse(http.StatusUnauthorized, "Fail", err)

		return c.JSONPretty(result.Status, result, "  ")
	}

	duser.Token = t
	_, err = authService.UpdateUser(duser)
	if err != nil {
		result = helper.CreateErrorResponse(http.StatusInternalServerError, "Something Wrong", err)
		return c.JSONPretty(result.Status, result, "  ")
	}

	result = helper.CreateSuccessResponse(Data{
		AccessToken: "Bearer " + t,
		User:        duser,
	}, "success")

	return c.JSONPretty(result.Status, result, "  ")
}

func Verify(c echo.Context, vault client.Vault) error {
	initAuth(vault)

	var user models.User
	var result helper.BaseResponse

	var bodyBytes []byte
	if c.Request().Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request().Body)
	}
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
		result = helper.CreateErrorResponse(http.StatusBadRequest, "Invalid request payload.", nil)
		return c.JSONPretty(result.Status, result, "  ")
	}

	if len(user.Email) <= 0 || len(user.TokenVerification) <= 0 {
		result = helper.CreateErrorResponse(http.StatusBadRequest, "Please fill all required field.", nil)
		return c.JSONPretty(result.Status, result, "  ")
	}

	duser, err := authService.GetUsersByEmail(user)
	if err != nil {
		result = helper.CreateErrorResponse(http.StatusBadRequest, "Email not found.", nil)
		return c.JSONPretty(result.Status, result, "  ")
	}

	if user.TokenVerification != duser.TokenVerification {
		result = helper.CreateErrorResponse(http.StatusBadRequest, "Invalid verification code", nil)
		return c.JSONPretty(result.Status, result, "  ")
	}

	duser.Status = 1
	duser.TokenVerification = ""

	user, err = authService.UpdateUser(duser)
	if err != nil {
		result = helper.CreateErrorResponse(http.StatusInternalServerError, "Something Wrong", err)
		return c.JSONPretty(result.Status, result, "  ")
	}

	result = helper.CreateSuccessResponse(user, "success")
	return c.JSONPretty(result.Status, result, "  ")
}
