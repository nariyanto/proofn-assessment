package dao

import (
	"fmt"

	"proofn/task5/models"

	"github.com/go-pg/pg"
)

type User struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

func (d *User) Connect() error {
	var n int

	//conn string
	db = pg.Connect(&pg.Options{
		User:     d.User,
		Password: d.Password,
		Addr:     fmt.Sprintf("%s:%s", d.Host, d.Port),
		Database: d.Database,
	})

	//Check our connection
	_, err := db.QueryOne(pg.Scan(&n), "SELECT 1")
	return err
}

func (d *User) Close() error {
	err := db.Close()
	return err
}

func (d *User) FindByEmail(user models.User) (models.User, error) {
	//Go get the users
	err := db.Model(&user).Where("email = ?", user.Email).Select()
	if err != nil {
		return user, err
	}

	return user, nil
}

func (d *User) FindAll() ([]models.User, error) {
	var users []models.User

	//Go get the users
	err := db.Model(&users).Select()
	if err != nil {
		return []models.User{}, err
	}

	return users, nil
}

func (d *User) DeleteAll() error {
	var ids []int

	//Find the user ids
	err := db.Model(&User{}).Column("id").Select(&ids)
	if err != nil {
		return err
	}

	//Delete the user ids if we have results
	if len(ids) > 0 {
		pgids := pg.In(ids)
		_, err := db.Model(&User{}).Where("id IN (?)", pgids).Delete()
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *User) Insert(user models.User) (models.User, error) {
	err := db.Insert(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (d *User) Update(user models.User) (models.User, error) {
	err := db.Update(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}
