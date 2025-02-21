package repositories

import (
	"context"
	"fmt"
	"store-first-login/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepositoryDB struct {
	collection *mongo.Collection
}

func NewUserRepositoryDB(db *mongo.Database) UserRepository {
	return userRepositoryDB{collection: db.Collection("user")}
}

// GetOne implements UserRepository.
func (u userRepositoryDB) GetOne(username string) (map[string]interface{}, error) {
	var user map[string]interface{}
	err := u.collection.FindOne(context.Background(), bson.M{"Username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no user found with username: %s", username)
		}
		return nil, err
	}
	return user, nil
}

// GetAllUsers implements UserRepository.
func (u userRepositoryDB) GetAll() ([]map[string]interface{}, error) {
	cursor, err := u.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	var users []map[string]interface{}
	if err := cursor.All(context.Background(), &users); err != nil {
		return nil, err
	}
	return users, nil
}

// Insert implements UserRepository.
func (u userRepositoryDB) Insert(user models.UserRegister) error {
	_, err := u.collection.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}
	return nil
}

// Update implements UserRepository.
func (u userRepositoryDB) Update(username string, updatedUser models.UserUpdate) error {
	// Create a filter to find the user by username
	filter := bson.M{
		"username":   username,
		"deleteDate": "",
	}
	// Create an update document with the fields you want to update
	update := bson.M{
		"$set": updatedUser, // Assuming you want to replace the user details
	}
	// Perform the update operation
	result, err := u.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	// Check if any documents were matched and modified
	if result.MatchedCount == 0 {
		return fmt.Errorf("no user found with username: %s", username)
	}
	return nil
}

// Delete implements UserRepository for soft delete.
func (u userRepositoryDB) Delete(username string) error {
	deleteDate := time.Now().UTC().Local().Format("2006-01-02T15:04:05.999-0700")
	// Create a filter to find the user by username
	filter := bson.M{
		"username":   username,
		"deleteDate": "",
	}

	// Create an update document to set the deletedAt timestamp
	update := bson.M{
		"$set": bson.M{"deleteDate": deleteDate}, // Set the delete timestamp field
	}

	// Perform the update operation (soft delete)
	result, err := u.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	// Check if any documents were matched and modified
	if result.MatchedCount == 0 {
		return fmt.Errorf("no user found with username: %s", username)
	}
	return nil
}
