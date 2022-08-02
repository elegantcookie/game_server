package db

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"training_service/internal/config"
	"training_service/internal/table"
	"training_service/pkg/logging"
)

type db struct {
	database *mongo.Database
	logger   *logging.Logger
}

func (d *db) IsCollection(ctx context.Context, collectionName string) (bool, error) {
	names, err := d.database.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return false, err
	}
	for i := 0; i < len(names); i++ {
		if names[i] == collectionName {
			return true, nil
		}
	}
	return false, nil

}

func (d *db) CreateCollection(ctx context.Context, dto table.CollectionDTO) error {
	d.logger.Debug("CREATE COLLECTION")
	if dto.AccessKey != config.GetConfig().Keys.AccessKey {
		return fmt.Errorf("wrong access key")
	}
	collectionName := dto.Name
	if collectionName == "" {
		return fmt.Errorf("collection name can't be empty")
	}
	isCollection, err := d.IsCollection(ctx, collectionName)
	if err != nil {
		return err
	}
	if isCollection {
		return fmt.Errorf("collection with name: %s already exists", collectionName)
	}
	collection := d.database.Collection(collectionName)

	// COLLECTION CREATES ONLY IF SOMETHING IS INSERTED THERE
	_, err = collection.InsertOne(ctx, table.RecordDTO{
		TableName: dto.Name,
		ID:        "-1",
		UserID:    "-1",
		UserScore: "-1",
	})
	if err != nil {
		return err
	}
	_, err = collection.DeleteOne(ctx, bson.M{})
	if err != nil {
		return err
	}
	d.logger.Printf("new collection name: %s", collection.Name())
	return nil
}

func (d *db) DeleteCollectionByName(ctx context.Context, dto table.CollectionDTO) error {
	d.logger.Debug("DELETE COLLECTION")
	if dto.AccessKey != config.GetConfig().Keys.AccessKey {
		return fmt.Errorf("wrong access key")
	}
	collectionName := dto.Name
	if collectionName == "" {
		return fmt.Errorf("collection name can't be empty")
	}
	isCollection, err := d.IsCollection(ctx, collectionName)
	if err != nil {
		return err
	}
	if !isCollection {
		return fmt.Errorf("collection with name: %s doesn't exists", collectionName)
	}
	d.database.Collection(collectionName).Drop(ctx)
	return nil
}

func (d *db) Create(ctx context.Context, dto table.RecordDTO) (string, error) {
	d.logger.Debug("create record")
	isCollection, err := d.IsCollection(ctx, dto.TableName)
	if err != nil {
		return "", err
	}
	if !isCollection {
		return "", fmt.Errorf("no collection found with name: %s", dto.TableName)
	}
	if dto.UserScore == "" || dto.UserID == "" {
		return "", fmt.Errorf("wrong userscore or userid")
	}
	record := table.Record{
		ID:        dto.ID,
		UserID:    dto.UserID,
		UserScore: dto.UserScore,
	}
	collection := d.database.Collection(dto.TableName)

	result, err := collection.InsertOne(ctx, record)
	if err != nil {
		return "", fmt.Errorf("failed to create record due to: %v", err)
	}
	d.logger.Debug("convert InsertedID to objectID")
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	d.logger.Trace(record)
	return "", fmt.Errorf("failed to convert objectId to hex. probable oid: %s", oid)
}

// FindById find user by recordID
func (d *db) FindById(ctx context.Context, dto table.RecordDTO) (record table.Record, err error) {
	isCollection, err := d.IsCollection(ctx, dto.TableName)
	if err != nil {
		return record, err
	}
	if !isCollection {
		return record, fmt.Errorf("no collection found with name: %s", dto.TableName)
	}
	id := dto.ID
	tableName := dto.TableName
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return record, fmt.Errorf("failed to convert hex to objectID, hex: %s", id)
	}
	filter := bson.M{"_id": oid}
	collection := d.database.Collection(tableName)
	result := collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			// TODO ErrEntityNotFound
		}
		return record, fmt.Errorf("failed to find record by id: %s due to error: %v", id, result.Err())
	}
	if err = result.Decode(&record); err != nil {
		return record, fmt.Errorf("failed to decode record(id:%s) from DB due to error: %v", id, err)
	}
	return record, nil
}

func (d *db) FindByUserId(ctx context.Context, dto table.RecordDTO) (record table.Record, err error) {
	isCollection, err := d.IsCollection(ctx, dto.TableName)
	if err != nil {
		return record, err
	}
	if !isCollection {
		return record, fmt.Errorf("no collection found with name: %s", dto.TableName)
	}
	id := dto.UserID
	tableName := dto.TableName
	if err != nil {
		return record, fmt.Errorf("failed to convert hex to objectID, hex: %s", id)
	}
	filter := bson.M{"user_id": id}
	collection := d.database.Collection(tableName)
	result := collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			// TODO ErrEntityNotFound
		}
		return record, fmt.Errorf("failed to find record by userID: %s due to error: %v", id, result.Err())
	}
	if err = result.Decode(&record); err != nil {
		return record, fmt.Errorf("failed to decode record(userID:%s) from DB due to error: %v", id, err)
	}
	return record, nil
}

func (d *db) FindAll(ctx context.Context, dto table.RecordDTO) (users []table.Record, err error) {
	isCollection, err := d.IsCollection(ctx, dto.TableName)
	if err != nil {
		return users, err
	}
	if !isCollection {
		return users, fmt.Errorf("no collection found with name: %s", dto.TableName)
	}
	tableName := dto.TableName
	collection := d.database.Collection(tableName)
	cursor, err := collection.Find(ctx, bson.M{})
	if cursor.Err() != nil {
		return users, fmt.Errorf("failed to find all users due to: %v", cursor.Err())
	}
	if err := cursor.All(ctx, &users); err != nil {
		return users, fmt.Errorf("failed to read all documents from cursor")
	}
	return users, nil
}

func (d *db) FindCollectionNames(ctx context.Context) ([]string, error) {
	return d.database.ListCollectionNames(ctx, bson.M{})
}

// Update by userID
func (d *db) Update(ctx context.Context, dto table.RecordDTO) error {
	isCollection, err := d.IsCollection(ctx, dto.TableName)
	if err != nil {
		return err
	}
	if !isCollection {
		return fmt.Errorf("no collection found with name: %s", dto.TableName)
	}

	if dto.UserScore == "" {
		return fmt.Errorf("wrong userscore")
	}

	record := table.Record{
		ID:        dto.ID,
		UserID:    dto.UserID,
		UserScore: dto.UserScore,
	}
	tableName := dto.TableName

	filter := bson.M{"user_id": dto.UserID}

	userBytes, err := bson.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal record due to: %v", err)
	}

	var updateUserObj bson.M
	err = bson.Unmarshal(userBytes, &updateUserObj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal record bytes due to: %v", err)
	}
	delete(updateUserObj, "user_id")
	update := bson.M{
		"$set": updateUserObj,
	}
	collection := d.database.Collection(tableName)
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to execute update record query due to: %v", err)
	}

	if result.MatchedCount == 0 {
		// TODO ErrEntityNotFound
		return fmt.Errorf("not found")
	}

	return nil
}

// Delete record by recordID
func (d *db) Delete(ctx context.Context, dto table.RecordDTO) error {
	isCollection, err := d.IsCollection(ctx, dto.TableName)
	if err != nil {
		return err
	}
	if !isCollection {
		return fmt.Errorf("no collection found with name: %s", dto.TableName)
	}

	id := dto.ID
	tableName := dto.TableName
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert record ID to ObjectID. ID=%v", id)
	}

	filter := bson.M{"_id": objectID}
	collection := d.database.Collection(tableName)
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execute update record query due to: %v", err)
	}

	if result.DeletedCount == 0 {
		// TODO ErrEntityNotFound
		return fmt.Errorf("not found")
	}
	d.logger.Tracef("Deleted %d documents", result.DeletedCount)
	return nil
}

func NewStorage(database *mongo.Database, logger *logging.Logger) table.Storage {
	return &db{
		database: database,
		logger:   logger,
	}
}
