package services

import (
	"fmt"
	"store-first-login/errs"
	"store-first-login/logs"
	"store-first-login/models"
	"store-first-login/repositories"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return userService{userRepo: userRepo}
}

// InsertUser implements UserService.
func (s userService) InsertUser(user models.UserRegister) error {
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return errs.NewUnexpectedError()
	}
	user.Password = hashedPassword
	user.CreateDate = time.Now().UTC().Local().Format("2006-01-02T15:04:05.999-0700")
	err = s.userRepo.Insert(user)
	if err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}
	logs.Info(fmt.Sprintf("user %s Created", user.Username))
	return nil
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// UpdateUser implements UserService.
func (s userService) UpdateUser(username string, user models.UserUpdate) error {
	if user.Password != "" {
		hashedPassword, err := hashPassword(user.Password)
		if err != nil {
			return errs.NewUnexpectedError()
		}
		user.Password = hashedPassword
	}
	err := s.userRepo.Update(username, user)
	if err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}
	return nil
}

// DeleteUser implements UserService.
func (s userService) DeleteUser(username string) error {
	err := s.userRepo.Delete(username)
	if err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}
	return nil
}
