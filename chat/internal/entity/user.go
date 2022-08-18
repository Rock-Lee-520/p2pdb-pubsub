package entity

import (
	"fmt"
)

type User struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
	Avatar      string `json:"avatar"`
}

func NewUser(id, name string) *User {
	return &User{
		Id:          id,
		DisplayName: name,
		Avatar:      fmt.Sprintf("https://api.multiavatar.com/v1/%s", id),
	}
}
