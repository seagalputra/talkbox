package api

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	RoomType string

	Participant struct {
		ID       primitive.ObjectID `bson:"id,omitempty"`
		Username string             `bson:"username"`
		Email    string             `bson:"email"`
	}

	Room struct {
		ID           primitive.ObjectID `bson:"_id,omitempty"`
		Participants []Participant      `bson:"participants"`
		RoomType     RoomType           `bson:"roomType"`
		CreatedAt    *time.Time         `bson:"created_at,omitempty"`
		UpdatedAt    *time.Time         `bson:"updated_at,omitempty"`
		DeletedAt    *time.Time         `bson:"deleted_at,omitempty"`
	}
)

const (
	Private RoomType = "private"
	Group   RoomType = "group"

	rooms string = "rooms"
)

func FindRoomByID(id string) (*Room, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id": objID,
	}

	var room Room
	if err := MongoDatabase.Collection(rooms).FindOne(context.Background(), filter).Decode(&room); err != nil {
		return nil, err
	}

	return &room, nil
}
