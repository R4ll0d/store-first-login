package handlers

import (
	"store-first-login/models"
	"store-first-login/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	userSrv services.UserService
}

func NewUserHandler(userSrv services.UserService) userHandler {
	return userHandler{userSrv: userSrv}
}
func (h *userHandler) GetUserHandler(c *fiber.Ctx) error {
	users, err := h.userSrv.GetUsers()
	if err != nil {
		return c.Status(fiber.StatusCreated).JSON(models.ResponseJson{
			StatusCode: strconv.Itoa(fiber.StatusCreated),
			Status:     "error",
			Message:    "User get unsuccessfully",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(models.ResponseData{
		StatusCode: strconv.Itoa(fiber.StatusCreated),
		Status:     "success",
		Message:    "User get successfully",
		Data:       users,
	})
}
