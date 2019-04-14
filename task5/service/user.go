package service

import (
	"time"

	"proofn/task5/client"
	"proofn/task5/dao"
	"proofn/task5/models"
)

type User struct {
	Vault      *client.Vault
	Dao        *dao.User
	Encyrption Transit
}

func (o *User) GetUsersByEmail(user models.User) (models.User, error) {

	user, err := o.Dao.FindByEmail(user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (o *User) GetUsers() (models.UsersResp, error) {
	var eUsers []models.User

	eUsers, err := o.Dao.FindAll()
	if err != nil {
		return models.UsersResp{}, err
	}

	//Create our response payload
	usersResp := models.UsersResp{}
	usersResp.Users = eUsers

	return usersResp, nil
}

func (o *User) CreateUser(user models.User) (models.User, error) {
	user.CreatedDate = time.Now()

	user, err := o.Dao.Insert(user)

	return user, err
}

func (o *User) DeleteUsers() error {
	err := o.Dao.DeleteAll()
	return err
}

func (o *User) UpdateUser(user models.User) (models.User, error) {
	user, err := o.Dao.Update(user)

	return user, err
}
