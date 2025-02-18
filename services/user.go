package services

type UserService interface {
	GetUsers() ([]map[string]interface{}, error)
}
