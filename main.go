package main

import (
	"fmt"
	"log"
	"strings"

	"store-first-login/handlers"
	"store-first-login/infrastructure"
	"store-first-login/logs"
	"store-first-login/repositories"
	"store-first-login/services"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func init() {
	initConfig()
	logs.InitLog()
}

func main() {
	db := infrastructure.InitMongo()
	// Initialize repositories, use cases, and handlers
	userRepo := repositories.NewUserRepositoryDB(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Create a new Fiber app
	app := fiber.New()

	// Define routes
	app.Post("store-first-login/user", func(c *fiber.Ctx) error {
		return userHandler.InsertUserHandler(c)
	})

	app.Put("store-first-login/user/:username", func(c *fiber.Ctx) error {
		return userHandler.UpdateUserHandler(c)
	})

	app.Delete("store-first-login/user/:username", func(c *fiber.Ctx) error {
		return userHandler.DeleteUserHandler(c)
	})

	// Start the server
	logs.Info(fmt.Sprintf("Server is running on port : %d", viper.GetInt("app.port")))
	if err := app.Listen(fmt.Sprintf(":%v", viper.GetInt("app.port"))); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
