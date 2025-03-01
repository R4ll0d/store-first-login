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
	app.Post("store-first-login/register/user", func(c *fiber.Ctx) error {
		return userHandler.InsertUserHandler(c)
	})

	app.Put("store-first-login/update/user/:username", func(c *fiber.Ctx) error {
		return userHandler.UpdateUserHandler(c)
	})

	app.Delete("store-first-login/delete/user/:username", func(c *fiber.Ctx) error {
		return userHandler.DeleteUserHandler(c)
	})

	app.Get("store-first-login/get/user/:username", func(c *fiber.Ctx) error {
		return userHandler.GetUserHandler(c)
	})

	app.Post("store-first-login/login/user", func(c *fiber.Ctx) error {
		return userHandler.LoginUserHandler(c)
	})
	// Start the server
	logs.Info(fmt.Sprintf("Server is running on port: %s", port))
	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func initConfig() {
	_ = godotenv.Load()
}
