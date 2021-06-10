package dao

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"notes-api/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type NotesDao struct {
	Client     *mongo.Client
	Database   string
	Collection string
}

func (dao *NotesDao) Ping(ctx context.Context) error {
	return dao.Client.Ping(ctx, readpref.Primary())
}

func (dao *NotesDao) GetNotes(ctx context.Context, filter map[string]interface{}) ([]models.Note, error) {
	cursor, err := dao.getCollection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var notes []models.Note
	if err := cursor.All(ctx, &notes); err != nil {
		return nil, err
	}

	return notes, nil
}

func (dao *NotesDao) UpdateNote(ctx context.Context, filter map[string]interface{}, updates bson.M) error {
	result := dao.getCollection().FindOneAndUpdate(ctx, filter, updates)
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (dao *NotesDao) DeleteNote(ctx context.Context, filter map[string]interface{}) error {
	result, err := dao.getCollection().DeleteOne(ctx, filter)
	if err != nil {
		return err
	} else if result.DeletedCount == 0 {
		return errors.New("no notes were deleted")
	}
	return nil
}

func (dao *NotesDao) CreateNote(ctx context.Context, note models.Note) error {
	_, err := dao.getCollection().InsertOne(ctx, note)
	if err != nil {
		return err
	}

	return nil
}

func (dao *NotesDao) getCollection() *mongo.Collection {
	return dao.Client.Database(dao.Database).Collection(dao.Collection)
}
