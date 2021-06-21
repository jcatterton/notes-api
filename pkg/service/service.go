package service

import (
	"context"

	"notes-api/pkg/models"
)

type NoteServiceHandler interface {
	Ping(ctx context.Context) error
	GetNotes(ctx context.Context, id string) ([]models.Note, error)
	UpdateNote(ctx context.Context, id string, noteRequest models.NoteRequest) error
	DeleteNote(ctx context.Context, id string) error
	CreateNote(ctx context.Context, noteRequest models.NoteRequest) (string, error)
	SendToContentService(ctx context.Context, id string) error
	ValidateToken(ctx context.Context, token string) error
	SetToken(token string)
}
