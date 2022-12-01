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
		ID       primitive.ObjectID `bson:"id,omitempty" json:"id"`
		Username string             `bson:"username" json:"username"`
		Email    string             `bson:"email" json:"email"`
	}

	Room struct {
		ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
		Participants []Participant      `bson:"participants" json:"participants"`
		RoomType     RoomType           `bson:"roomType" json:"roomType"`
		CreatedAt    *time.Time         `bson:"createdAt,omitempty" json:"createdAt"`
		UpdatedAt    *time.Time         `bson:"updatedAt,omitempty" json:"updatedAt"`
		DeletedAt    *time.Time         `bson:"deletedAt,omitempty" json:"-"`
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
