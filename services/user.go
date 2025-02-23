package services

import (
	"store-first-login/models"
)

type UserService interface {
	InsertUser(models.UserRegister) error
	UpdateUser(string, models.UserUpdate) error
	DeleteUser(string) error
}
