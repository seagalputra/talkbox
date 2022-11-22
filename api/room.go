package api

import "time"

type (
	RoomType string

	Participant struct {
		ID       string `bson:"_id"`
		Username string `bson:"username"`
		Email    string `bson:"email"`
	}

	Room struct {
		ID           string        `bson:"_id"`
		Participants []Participant `bson:"participants"`
		RoomType     RoomType      `bson:"roomType"`
		CreatedAt    *time.Time    `bson:"created_at"`
		UpdatedAt    *time.Time    `bson:"updated_at"`
		DeletedAt    *time.Time    `bson:"deleted_at"`
	}
)

const (
	Private RoomType = "private"
	Group   RoomType = "group"
)
