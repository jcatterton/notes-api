package service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"mime/multipart"
	"notes-api/pkg/dao"
	"notes-api/pkg/external"
	"notes-api/pkg/models"
	"path/filepath"
	"strings"
	"time"
)

type NotesService struct {
	Dao dao.NotesDao
	Ext external.ExtAPI
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

func (svc *NotesService) SendToContentService(ctx context.Context, id string) error {
	logger := logrus.WithContext(ctx)

	notes, err := svc.GetNotes(ctx, id)
	if err != nil {
		return err
	}
	note := notes[0]

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	fileName := strings.Replace(note.Name, " ", "", -1)
	extension := filepath.Ext(fileName)
	if extension != "" {
		fileName = fileName[0 : len(fileName)-len(extension)]
	}
	fileName = fmt.Sprintf("%v.txt", fileName)

	formFile, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return err
	}

	if _, err := io.Copy(formFile, bytes.NewBuffer([]byte(note.Text))); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		logger.WithError(err).Error("Error closing multipart writer")
	}

	if err := svc.Ext.SendToContentService(ctx, body, writer.FormDataContentType()); err != nil {
		return err
	}

	return nil
}

func (svc *NotesService) ValidateToken(ctx context.Context, token string) error {
	if err := svc.Ext.ValidateToken(ctx, token); err != nil {
		return err
	}

	return nil
}

func (svc *NotesService) SetToken(token string) {
	svc.Ext.Token = token
}
