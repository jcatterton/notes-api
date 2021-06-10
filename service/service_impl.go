package service

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"notes-api/dao"
	"notes-api/models"
	"time"
)

type NotesService struct {
	Dao dao.NotesDao
}

func (svc *NotesService) Ping(ctx context.Context) error {
	return svc.Dao.Ping(ctx)
}

func (svc *NotesService) GetNotes(ctx context.Context, id string) ([]models.Note, error) {
	filter := make(map[string]interface{})

	if id != "" {
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		filter["_id"] = objectId
	}

	return svc.Dao.GetNotes(ctx, filter)
}

func (svc *NotesService) UpdateNote(ctx context.Context, id string, noteRequest models.NoteRequest) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := map[string]interface{}{
		"_id": objectId,
	}

	updates := bson.M{
		"$set": bson.M{
			"name": noteRequest.Name,
			"text": noteRequest.Text,
		},
	}

	return svc.Dao.UpdateNote(ctx, filter, updates)
}

func (svc *NotesService) DeleteNote(ctx context.Context, id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := map[string]interface{}{
		"_id": objectId,
	}

	if err := svc.Dao.DeleteNote(ctx, filter); err != nil {
		return err
	}

	return nil
}

func (svc *NotesService) CreateNote(ctx context.Context, noteRequest models.NoteRequest) (string, error) {
	id := primitive.NewObjectID()

	note := models.Note{
		ID:           id,
		Name:         noteRequest.Name,
		LastEditedTs: time.Now(),
		Text:         noteRequest.Text,
	}

	if err := svc.Dao.CreateNote(ctx, note); err != nil {
		return "", err
	}

	return id.Hex(), nil
}
