package infrastructure

import (
	"context"
	"fmt"
	"log"
	"store-first-login/logs"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func ConnectMongoDB(uri, dbName string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// ตรวจสอบการเชื่อมต่อ
	if err := client.Ping(ctx, nil); err != nil {
		client.Disconnect(ctx) // ปิดการเชื่อมต่อถ้าล้มเหลว
		return nil, err
	}

	// เก็บ client ไว้เพื่อใช้ใน HealthCheck
	mongoClient = client
	return client.Database(dbName), nil
}

func InitMongo() *mongo.Database {
	// Connect to MongoDB
	db, err := ConnectMongoDB(viper.GetString("mongo.uri"), viper.GetString("mongo.database"))
	if err != nil {
		log.Fatal("Connect Mongo Error:", err)
	}
	logs.Info("Connect Mongo Successfully!")

	// ตรวจสอบสุขภาพของ MongoDB
	go HealthCheckMongo()

	return db
}

func HealthCheckMongo() {
	ticker := time.NewTicker(10 * time.Second) // เช็คทุก 10 วินาที
	defer ticker.Stop()

	for range ticker.C {
		if mongoClient == nil {
			logs.Info("MongoDB client is nil, skipping health check")
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := mongoClient.Ping(ctx, nil)
		cancel()

		if err != nil {
			logs.Info(fmt.Sprintf("MongoDB connection lost: %s", err))
		}
		// } else {
		// 	logs.Info("MongoDB is healthy")
		// }
	}
}
