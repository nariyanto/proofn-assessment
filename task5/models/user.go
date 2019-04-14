package models

import "time"

type User struct {
	ID                int64     `json:"id"`
	Name              string    `json:"name"`
	Email             string    `json:"email"`
	Password          string    `json:"password"`
	Token             string    `json:"token"`
	TokenVerification string    `json:"tokenVerification"`
	Status            int8      `json:"status"`
	CreatedDate       time.Time `json:"createdDate"`
}

type UsersResp struct {
	Users []User `json:"users"`
}
