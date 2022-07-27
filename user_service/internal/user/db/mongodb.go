package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"user_service/internal/user"
	"user_service/pkg/logging"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func (d *db) Create(ctx context.Context, user user.User) (string, error) {
	d.logger.Debug("create user")
	result, err := d.collection.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to create user due to: %v", err)
	}
	d.logger.Debug("convert InsertedID to objectID")
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	d.logger.Trace(user)
	return "", fmt.Errorf("failed to convert objectId to hex. probable oid: %s", oid)
}
func (d *db) FindOne(ctx context.Context, id string) (u user.User, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return u, fmt.Errorf("failed to convert hex to objectID, hex: %s", id)
	}
	filter := bson.M{"_id": oid}
	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		// TODO 404
		return u, fmt.Errorf("failed to find one user by id: %s due to error: %v", id, result.Err())
	}
	if err = result.Decode(&u); err != nil {
		return u, fmt.Errorf("failed to decode user(id:%s) from DB due to error: %v", id, err)
	}
	return u, nil

}
func (d *db) Update(ctx context.Context, user user.User) error {

}
func (d *db) Delete(ctx context.Context, id string) error {

}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {

	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}
