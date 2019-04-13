package models

import "time"

type User struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	CreatedDate time.Time `json:"createdDate"`
}

type UsersResp struct {
	Users []User `json:"users"`
}
