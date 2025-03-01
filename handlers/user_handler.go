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
