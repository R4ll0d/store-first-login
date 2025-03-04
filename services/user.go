package services

import (
	"store-first-login/models"
)

type UserService interface {
	InsertUser(models.UserRegister) error
	UpdateUser(string, models.UserUpdate) error
	DeleteUser(string) error
	GetUser(username string) (models.UserDetail, error)
	LoginUser(input models.UserLogin) (string, error)
	SendOTP(user models.UserRegister) error
}
