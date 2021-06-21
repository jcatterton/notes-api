package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NoteRequest struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

type Note struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	LastEditedTs time.Time          `json:"lastEditedTs" bson:"lastEditedTs"`
	Text         string             `json:"text" bson:"text"`
}
