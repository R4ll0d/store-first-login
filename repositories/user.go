package repositories

type UserRepository interface {
	GetAll() ([]map[string]interface{}, error)
}
