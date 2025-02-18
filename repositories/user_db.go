package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepositoryDB struct {
	collection *mongo.Collection
}

func NewUserRepositoryDB(db *mongo.Database) UserRepository {
	return userRepositoryDB{collection: db.Collection("profile")}
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
