package services

import (
	"encoding/json"
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
	if user.Password != user.ConfirmPassword {
		logs.Error("Invalid ConfirmPassword")
		return errs.NewValidationError("Invalid ConfirmPassword")
	}
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		logs.Error("Invalid hashPassword")
		return errs.NewValidationError("Invalid hashPassword")
	}
	user.Password = hashedPassword
	user.CreateDate = time.Now().UTC().Local().Format("2006-01-02T15:04:05.999-0700")
	insertedID, err := s.userRepo.Insert(user)
	if err != nil {
		if insertedID == nil {
			logs.Error("Username Already Exits")
			return errs.NewAlreadyExits("Already Exits")
		}
		logs.Error("Unexpected error")
		return errs.NewUnexpectedError()
	}
	logs.Info(fmt.Sprintf("User %s Created", user.Username))
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
			logs.Error("Invalid hashPassword")
			return errs.NewValidationError("Invalid hashPassword")
		}
		user.Password = hashedPassword
	}
	result, err := s.userRepo.Update(username, user)
	if err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}
	if result.MatchedCount == 0 {
		logs.Error("User Not Found")
		return errs.NewNotFoundError("User Not Found")
	}
	if result.ModifiedCount == 0 {
		logs.Error(errs.NewUnexpectedError())
		return errs.NewUnexpectedError()
	}
	logs.Info(fmt.Sprintf("User %s Updated with details: %+v", username, user))
	return nil
}

// DeleteUser implements UserService.
func (s userService) DeleteUser(username string) error {
	err := s.userRepo.Delete(username)
	if err != nil {
		logs.Error(err.Error())
		return errs.NewUnexpectedError()
	}
	logs.Info(fmt.Sprintf("User %s Deleted", username))
	return nil
}

// GetUser implements UserService.
func (s userService) GetUser(username string) (models.UserDetail, error) {
	// result เป็น map[string]interface{}
	result, err := s.userRepo.GetOne(username)
	if err != nil {
		logs.Error(err.Error())
		return models.UserDetail{}, err // คืนค่า struct เปล่าแทน nil
	}

	// แปลง map[string]interface{} เป็น JSON
	jsonData, err := json.Marshal(result)
	if err != nil {
		logs.Error(err.Error())
		return models.UserDetail{}, fmt.Errorf("failed to marshal result: %v", err)
	}

	// แปลง JSON เป็น models.UserDetail
	var userDetails models.UserDetail
	if err := json.Unmarshal(jsonData, &userDetails); err != nil {
		logs.Error(err.Error())
		return models.UserDetail{}, fmt.Errorf("failed to unmarshal to models.UserDetail: %v", err)
	}
	logs.Info(fmt.Sprintf("User %s Geted with details: %+v", username, userDetails))
	return userDetails, nil
}
