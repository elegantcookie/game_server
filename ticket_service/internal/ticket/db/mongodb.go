package db

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"ticket_service/internal/ticket"
	"ticket_service/pkg/logging"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func (d *db) Create(ctx context.Context, ticket ticket.Ticket) (string, error) {
	result, err := d.collection.InsertOne(ctx, ticket)
	if err != nil {
		return "", fmt.Errorf("failed to create prize due to: %v", err)
	}
	d.logger.Debug("convert InsertedID to objectID")
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	d.logger.Trace(ticket)
	return "", fmt.Errorf("failed to convert objectId to hex. probable oid: %s", oid)
}

// FindById find prize by ticketID
func (d *db) FindById(ctx context.Context, id string) (ticket ticket.Ticket, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ticket, fmt.Errorf("failed to convert hex to objectID, hex: %s", id)
	}
	filter := bson.M{"_id": oid}
	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			// TODO ErrEntityNotFound
		}
		return ticket, fmt.Errorf("failed to find prize by id: %s due to error: %v", id, result.Err())
	}
	if err = result.Decode(&ticket); err != nil {
		return ticket, fmt.Errorf("failed to decode prize(id:%s) from DB due to error: %v", id, err)
	}
	return ticket, nil
}

func (d *db) FindAll(ctx context.Context) (users []ticket.Ticket, err error) {
	cursor, err := d.collection.Find(ctx, bson.M{})
	if cursor.Err() != nil {
		return users, fmt.Errorf("failed to find all tickets due to: %v", cursor.Err())
	}
	if err := cursor.All(ctx, &users); err != nil {
		return users, fmt.Errorf("failed to read all documents from cursor")
	}
	return users, nil
}

// Update by ticketID
func (d *db) Update(ctx context.Context, ticket ticket.Ticket) error {
	objectID, err := primitive.ObjectIDFromHex(ticket.ID)
	if err != nil {
		return fmt.Errorf("failed to convert prize ID to ObjectID. ID=%v", ticket.ID)
	}

	filter := bson.M{"_id": objectID}

	userBytes, err := bson.Marshal(ticket)
	if err != nil {
		return fmt.Errorf("failed to marshal prize due to: %v", err)
	}

	var updateUserObj bson.M
	err = bson.Unmarshal(userBytes, &updateUserObj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal prize bytes due to: %v", err)
	}
	delete(updateUserObj, "_id")
	update := bson.M{
		"$set": updateUserObj,
	}
	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to execute update prize query due to: %v", err)
	}

	if result.MatchedCount == 0 {
		// TODO ErrEntityNotFound
		return fmt.Errorf("not found")
	}
	return nil
}

// Delete prize by ticketID
func (d *db) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert prize ID to ObjectID. ID=%v", id)
	}

	filter := bson.M{"_id": objectID}
	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execute update prize query due to: %v", err)
	}

	if result.DeletedCount == 0 {
		// TODO ErrEntityNotFound
		return fmt.Errorf("not found")
	}
	d.logger.Tracef("Deleted %d documents", result.DeletedCount)
	return nil
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) ticket.Storage {

	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}
