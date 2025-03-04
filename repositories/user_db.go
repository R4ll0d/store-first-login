package repositories

import (
	"context"
	"fmt"
	"store-first-login/models"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	err := u.collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
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
func (u userRepositoryDB) Insert(user models.UserRegister) (interface{}, error) {
	result, err := u.collection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

// Update implements UserRepository.
func (u userRepositoryDB) Update(username string, updatedUser models.UserUpdate) (*mongo.UpdateResult, error) {
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
		return nil, err
	}
	return result, nil
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

func (u userRepositoryDB) GetByFilter(dbName string, collection string, filter string, field string) ([]map[string]interface{}, error) {
	// Parse the filter and field into BSON format
	filterBson, err := parseFilterStringToBson(filter)
	if err != nil {
		filterBson = bson.M{}
	}

	fieldBson, err := parseFieldStringToBson(field)
	if err != nil {
		fieldBson = bson.M{}
	}

	// Define options for Find based on projection
	findOptions := options.Find().SetProjection(fieldBson)

	// Perform the query with appropriate filter and projection
	cursor, err := u.collection.Find(context.Background(), filterBson, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background()) // Ensure the cursor is closed after use

	// Collect results
	var results []map[string]interface{}
	for cursor.Next(context.Background()) {
		var doc map[string]interface{}
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("error decoding document: %v", err)
		}
		results = append(results, doc)
	}

	// Check for errors encountered during cursor iteration
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	// Return results
	if len(results) == 0 {
		return nil, fmt.Errorf("mongo: no documents in result")
	}
	return results, nil
}

func parseFilterStringToBson(filterString string) (bson.M, error) {
	filterBson := bson.M{}
	// Removing parentheses and splitting by comma
	parts := strings.Split(strings.Trim(filterString, "()"), ",")

	for _, part := range parts {
		keyValue := strings.Split(part, "=")
		if len(keyValue) != 2 {
			return nil, fmt.Errorf("invalid key-value pair: %s", part)
		}

		key := strings.TrimSpace(keyValue[0])
		value := strings.TrimSpace(keyValue[1])

		// Creating BSON document
		filterBson[key] = value
	}
	return filterBson, nil
}

func parseFieldStringToBson(fieldString string) (bson.M, error) {
	fieldBson := bson.M{}

	// Removing parentheses and splitting by comma
	parts := strings.Split(strings.Trim(fieldString, "()"), ",")

	for _, part := range parts {
		key := strings.TrimSpace(part)
		if key == "" {
			return nil, fmt.Errorf("invalid key-value pair: %s", part)
		}
		fieldBson[key] = 1
	}
	return fieldBson, nil
}
