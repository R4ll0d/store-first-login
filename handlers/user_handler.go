package handlers

import (
	"store-first-login/models"
	"store-first-login/services"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	userSrv services.UserService
}

func NewUserHandler(userSrv services.UserService) userHandler {
	return userHandler{userSrv: userSrv}
}
func (h *userHandler) InsertUserHandler(c *fiber.Ctx) error {
	var user models.UserRegister
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ResponseJson{
			StatusCode: strconv.Itoa(fiber.StatusBadRequest),
			Status:     "error",
			Message:    "Invalid request body",
		})
	}

	err := h.userSrv.InsertUser(user)
	if err != nil {
		if strings.Contains(err.Error(), "Already Exits") {
			return c.Status(fiber.StatusNotFound).JSON(models.ResponseJson{
				StatusCode: strconv.Itoa(fiber.StatusNotFound),
				Status:     "error",
				Message:    "user already exits",
			})
		}
		if strings.Contains(err.Error(), "Invalid ConfirmPassword") {
			return c.Status(fiber.StatusNotFound).JSON(models.ResponseJson{
				StatusCode: strconv.Itoa(fiber.StatusNotFound),
				Status:     "error",
				Message:    "password and confirm password do not match",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ResponseJson{
			StatusCode: strconv.Itoa(fiber.StatusInternalServerError),
			Status:     "error",
			Message:    "Failed to insert user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.ResponseJson{
		StatusCode: strconv.Itoa(fiber.StatusCreated),
		Status:     "success",
		Message:    "User created successfully",
	})
}
func (h *userHandler) UpdateUserHandler(c *fiber.Ctx) error {
	var user models.UserUpdate
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ResponseJson{
			StatusCode: strconv.Itoa(fiber.StatusBadRequest),
			Status:     "error",
			Message:    "Invalid request body",
		})
	}
	username := c.Params("username")
	err := h.userSrv.UpdateUser(username, user)
	if err != nil {
		if strings.Contains(err.Error(), "unexpected error") {
			return c.Status(fiber.StatusNotFound).JSON(models.ResponseJson{
				StatusCode: strconv.Itoa(fiber.StatusNotFound),
				Status:     "error",
				Message:    "unexpected error",
			})
		}
		if strings.Contains(err.Error(), "User Not Found") {
			return c.Status(fiber.StatusNotFound).JSON(models.ResponseJson{
				StatusCode: strconv.Itoa(fiber.StatusNotFound),
				Status:     "error",
				Message:    "user not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ResponseJson{
			StatusCode: strconv.Itoa(fiber.StatusInternalServerError),
			Status:     "error",
			Message:    "Failed to insert user",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(models.ResponseJson{
		StatusCode: strconv.Itoa(fiber.StatusCreated),
		Status:     "success",
		Message:    "User updated successfully",
	})
}
func (h *userHandler) DeleteUserHandler(c *fiber.Ctx) error {
	username := c.Params("username")
	err := h.userSrv.DeleteUser(username)
	if err != nil {
		if strings.Contains(err.Error(), "unexpected error") {
			return c.Status(fiber.StatusNotFound).JSON(models.ResponseJson{
				StatusCode: strconv.Itoa(fiber.StatusNotFound),
				Status:     "error",
				Message:    "user not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ResponseJson{
			StatusCode: strconv.Itoa(fiber.StatusInternalServerError),
			Status:     "error",
			Message:    "Failed to insert user",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(models.ResponseJson{
		StatusCode: strconv.Itoa(fiber.StatusCreated),
		Status:     "success",
		Message:    "User deleted successfully",
	})
}

func (h *userHandler) GetUserHandler(c *fiber.Ctx) error {
	username := c.Params("username")
	result, err := h.userSrv.GetUser(username)
	if err != nil {
		if strings.Contains(err.Error(), "unexpected error") {
			return c.Status(fiber.StatusNotFound).JSON(models.ResponseJson{
				StatusCode: strconv.Itoa(fiber.StatusNotFound),
				Status:     "error",
				Message:    "user not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ResponseJson{
			StatusCode: strconv.Itoa(fiber.StatusInternalServerError),
			Status:     "error",
			Message:    "Failed to find user",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(models.ResponseUser{
		StatusCode: strconv.Itoa(fiber.StatusCreated),
		Status:     "success",
		Message:    "User geted successfully",
		Data:       result,
	})
}
func (h *userHandler) LoginUserHandler(c *fiber.Ctx) error {
	var user models.UserLogin
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ResponseJson{
			StatusCode: strconv.Itoa(fiber.StatusBadRequest),
			Status:     "error",
			Message:    "Invalid request body",
		})
	}

	jwtToken, err := h.userSrv.LoginUser(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ResponseJson{
			StatusCode: strconv.Itoa(fiber.StatusInternalServerError),
			Status:     "error",
			Message:    "Failed to login",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.ResponseUserLogin{
		StatusCode: strconv.Itoa(fiber.StatusCreated),
		Status:     "success",
		Message:    "Login successfully",
		Data: struct {
			JwtToken string `bson:"jwt-token"`
		}{JwtToken: jwtToken},
	})
}

// Send OTP
func (h *userHandler) SendOTPHandler(c *fiber.Ctx) error {
	var request models.UserRegister
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ResponseJson{
			StatusCode: strconv.Itoa(fiber.StatusBadRequest),
			Status:     "error",
			Message:    "Invalid request body",
		})
	}

	err := h.userSrv.SendOTP(request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ResponseJson{
			StatusCode: strconv.Itoa(fiber.StatusInternalServerError),
			Status:     "error",
			Message:    "Failed to send OTP",
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.ResponseJson{
		StatusCode: strconv.Itoa(fiber.StatusOK),
		Status:     "success",
		Message:    "OTP sent successfully",
	})
}

// // Validate OTP
// func (h *userHandler) ValidateOTPHandler(c *fiber.Ctx) error {
// 	var request models.ValidateOTPRequest
// 	if err := c.BodyParser(&request); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(models.ResponseJson{
// 			StatusCode: strconv.Itoa(fiber.StatusBadRequest),
// 			Status:     "error",
// 			Message:    "Invalid request body",
// 		})
// 	}

// 	valid, err := h.userSrv.ValidateOTP(request)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(models.ResponseJson{
// 			StatusCode: strconv.Itoa(fiber.StatusInternalServerError),
// 			Status:     "error",
// 			Message:    "Failed to validate OTP",
// 		})
// 	}

// 	if !valid {
// 		return c.Status(fiber.StatusUnauthorized).JSON(models.ResponseJson{
// 			StatusCode: strconv.Itoa(fiber.StatusUnauthorized),
// 			Status:     "error",
// 			Message:    "Invalid OTP",
// 		})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(models.ResponseJson{
// 		StatusCode: strconv.Itoa(fiber.StatusOK),
// 		Status:     "success",
// 		Message:    "OTP validated successfully",
// 	})
// }
