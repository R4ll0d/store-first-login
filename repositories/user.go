package repositories

import (
	"store-first-login/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	GetAll() ([]map[string]interface{}, error)
	GetOne(string) (map[string]interface{}, error)
	Insert(models.UserRegister) (interface{}, error)
	Update(string, models.UserUpdate) (*mongo.UpdateResult, error)
	Delete(string) error
}
