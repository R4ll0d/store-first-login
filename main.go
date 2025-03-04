package main

import (
	"fmt"
	"log"
	"os"

	"store-first-login/handlers"
	"store-first-login/infrastructure"
	"store-first-login/logs"
	"store-first-login/repositories"
	"store-first-login/services"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func init() {
	initConfig()
	logs.InitLog()
}

func main() {
	port := os.Getenv("PORT") // Cloud Run ใช้ตัวแปร PORT
	if port == "" {
		port = "8080" // ตั้งค่าเริ่มต้นเป็น 8080
	}
	db := infrastructure.InitMongo()
	// Initialize repositories, use cases, and handlers
	userRepo := repositories.NewUserRepositoryDB(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Create a new Fiber app
	app := fiber.New()

	// Define routes

	app.Put("store-first-login/updateUser/:username", func(c *fiber.Ctx) error {
		return userHandler.UpdateUserHandler(c)
	})

	app.Delete("store-first-login/deleteUser/:username", func(c *fiber.Ctx) error {
		return userHandler.DeleteUserHandler(c)
	})

	app.Get("store-first-login/getUser/:username", func(c *fiber.Ctx) error {
		return userHandler.GetUserHandler(c)
	})

	app.Post("store-first-login/register", func(c *fiber.Ctx) error {
		return userHandler.InsertUserHandler(c)
	})

	app.Post("store-first-login/login", func(c *fiber.Ctx) error {
		return userHandler.LoginUserHandler(c)
	})

	app.Post("store-first-login/sendOTP", func(c *fiber.Ctx) error {
		return userHandler.SendOTPHandler(c)
	})

	// app.Post("store-first-login/validateOTP", func(c *fiber.Ctx) error {
	// 	return userHandler.ValidateOTPHandler(c)
	// })

	// Start the server
	logs.Info(fmt.Sprintf("Server is running on port: %s", port))
	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func initConfig() {
	_ = godotenv.Load()
}
