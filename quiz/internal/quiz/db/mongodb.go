package db

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"quiz_service/internal/quiz"
	"quiz_service/pkg/logging"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func (d *db) Create(ctx context.Context, gs quiz.Quiz) (string, error) {
	result, err := d.collection.InsertOne(ctx, gs)
	if err != nil {
		return "", fmt.Errorf("failed to create lobby due to: %v", err)
	}
	d.logger.Debug("convert InsertedID to objectID")
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	d.logger.Trace(gs)
	return "", fmt.Errorf("failed to convert objectId to hex. probable oid: %s", oid)
}

// FindById find lobby by gsID
func (d *db) FindById(ctx context.Context, id string) (gs quiz.Quiz, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return gs, fmt.Errorf("failed to convert hex to objectID, hex: %s", id)
	}
	filter := bson.M{"_id": oid}
	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			// TODO ErrEntityNotFound
		}
		return gs, fmt.Errorf("failed to find lobby by id: %s due to error: %v", id, result.Err())
	}
	if err = result.Decode(&gs); err != nil {
		return gs, fmt.Errorf("failed to decode lobby(id:%s) from DB due to error: %v", id, err)
	}
	return gs, nil
}

func (d *db) FindAll(ctx context.Context) (users []quiz.Quiz, err error) {
	cursor, err := d.collection.Find(ctx, bson.M{})
	if cursor.Err() != nil {
		return users, fmt.Errorf("failed to find all gss due to: %v", cursor.Err())
	}
	if err := cursor.All(ctx, &users); err != nil {
		return users, fmt.Errorf("failed to read all documents from cursor")
	}
	return users, nil
}

// Update by gsID
func (d *db) Update(ctx context.Context, gs quiz.Quiz) error {
	objectID, err := primitive.ObjectIDFromHex(gs.ID)
	if err != nil {
		return fmt.Errorf("failed to convert lobby ID to ObjectID. ID=%v", gs.ID)
	}

	filter := bson.M{"_id": objectID}

	userBytes, err := bson.Marshal(gs)
	if err != nil {
		return fmt.Errorf("failed to marshal lobby due to: %v", err)
	}

	var updateUserObj bson.M
	err = bson.Unmarshal(userBytes, &updateUserObj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal lobby bytes due to: %v", err)
	}
	delete(updateUserObj, "_id")
	update := bson.M{
		"$set": updateUserObj,
	}
	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to execute update lobby query due to: %v", err)
	}

	if result.MatchedCount == 0 {
		// TODO ErrEntityNotFound
		return fmt.Errorf("not found")
	}
	return nil
}

// Delete lobby by gsID
func (d *db) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert lobby ID to ObjectID. ID=%v", id)
	}

	filter := bson.M{"_id": objectID}
	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execute update lobby query due to: %v", err)
	}

	if result.DeletedCount == 0 {
		// TODO ErrEntityNotFound
		return fmt.Errorf("not found")
	}
	d.logger.Tracef("Deleted %d documents", result.DeletedCount)
	return nil
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) quiz.Storage {

	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}
