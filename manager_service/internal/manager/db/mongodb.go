package db

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"manager_service/internal/manager"
	"manager_service/pkg/logging"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func (d *db) Create(ctx context.Context, dto manager.LobbyRecordDTO) (string, error) {
	lr := manager.LobbyRecord{
		Type:       dto.Type,
		LobbyID:    dto.LobbyID,
		GameType:   dto.GameType,
		Expiration: dto.Expiration,
	}
	result, err := d.collection.InsertOne(ctx, lr)
	if err != nil {
		return "", fmt.Errorf("failed to create dto due to: %v", err)
	}
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	return "", fmt.Errorf("failed to convert objectId to hex. probable oid: %s", oid)
}

func (d *db) FindById(ctx context.Context, id string) (lr manager.LobbyRecord, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return lr, fmt.Errorf("failed to convert hex to objectID, hex: %s", id)
	}
	filter := bson.M{"_id": oid}
	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			// TODO ErrEntityNotFound
		}
		return lr, fmt.Errorf("failed to find lr by id: %s due to error: %v", id, result.Err())
	}
	if err = result.Decode(&lr); err != nil {
		return lr, fmt.Errorf("failed to decode lr(id:%s) from DB due to error: %v", id, err)
	}
	return lr, nil
}

func (d *db) FindAll(ctx context.Context) (lrs []manager.LobbyRecord, err error) {
	cursor, err := d.collection.Find(ctx, bson.D{})
	if cursor.Err() != nil {
		return lrs, fmt.Errorf("failed to find all lrs due to: %v", cursor.Err())
	}
	if err := cursor.All(ctx, &lrs); err != nil {
		return lrs, fmt.Errorf("failed to read all documents from cursor")
	}
	return lrs, nil
}

func (d *db) Update(ctx context.Context, lr manager.LobbyRecord) error {
	objectID, err := primitive.ObjectIDFromHex(lr.ID)
	if err != nil {
		return fmt.Errorf("failed to convert lr ID to ObjectID. ID=%v", lr.ID)
	}

	filter := bson.M{"_id": objectID}

	userBytes, err := bson.Marshal(lr)
	if err != nil {
		return fmt.Errorf("failed to marshal lr due to: %v", err)
	}

	var updateUserObj bson.M
	err = bson.Unmarshal(userBytes, &updateUserObj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal lr bytes due to: %v", err)
	}
	delete(updateUserObj, "_id")
	update := bson.M{
		"$set": updateUserObj,
	}
	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to execute update lr query due to: %v", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("not found")
	}

	return nil
}

func (d *db) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert lr ID to ObjectID. ID=%v", id)
	}

	filter := bson.M{"_id": objectID}

	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execute update lr query due to: %v", err)
	}

	if result.DeletedCount == 0 {
		// TODO ErrEntityNotFound
		return fmt.Errorf("not found")
	}
	return nil
}

func (d *db) DeleteAll(ctx context.Context) error {

	filter := bson.M{}

	result, err := d.collection.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execute update lr query due to: %v", err)
	}

	if result.DeletedCount == 0 {
		// TODO ErrEntityNotFound
		return fmt.Errorf("not found")
	}
	return nil
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) manager.Storage {

	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}
