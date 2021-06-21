package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"

	"notes-api/pkg/models"
)

type NoteDaoHandler interface {
	Ping(ctx context.Context) error
	GetNotes(ctx context.Context, filter map[string]interface{}) ([]models.Note, error)
	UpdateNote(ctx context.Context, filter map[string]interface{}, updates bson.M) error
	DeleteNote(ctx context.Context, filter map[string]interface{}) error
	CreateNote(ctx context.Context, note models.Note) error
}
