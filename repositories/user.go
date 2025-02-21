package repositories

import (
	"store-first-login/models"
)

type UserRepository interface {
	GetAll() ([]map[string]interface{}, error)
	GetOne(string) (map[string]interface{}, error)
	Insert(models.UserRegister) error
	Update(string, models.UserUpdate) error
	Delete(string) error
}
