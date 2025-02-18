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
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	initConfig()
	logs.InitLog()
	db := initMongo()
	// Initialize repositories, use cases, and handlers
	userRepo := repositories.NewUserRepositoryDB(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Create a new Fiber app
	app := fiber.New()

	// Define routes
	app.Get("/user", func(c *fiber.Ctx) error {
		return userHandler.GetUserHandler(c)
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
func initMongo() *mongo.Database {
	// Connect to MongoDB
	db, err := infrastructure.ConnectMongoDB(viper.GetString("mongo.uri"), viper.GetString("mongo.database"))
	if err != nil {
		logs.Error("Connect Mongo Error!")
	}
	logs.Info("Connect Mongo Successfully!")
	return db
}
