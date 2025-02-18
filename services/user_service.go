package services

import (
	"store-first-login/errs"
	"store-first-login/logs"
	"store-first-login/repositories"
)

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return userService{userRepo: userRepo}
}
func (s userService) GetUsers() ([]map[string]interface{}, error) {
	accounts, err := s.userRepo.GetAll()
	if err != nil {
		logs.Error(err)
		return nil, errs.NewUnexpectedError()
	}
	return accounts, nil
}
