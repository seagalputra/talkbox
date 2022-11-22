package api

import "time"

type (
	Message struct {
		ID         string     `bson:"_id"`
		Body       string     `bson:"body"`
		Attachment *string    `bson:"attachment"`
		UserID     string     `bson:"userId"`
		RoomID     string     `bson:"roomId"`
		CreatedAt  *time.Time `bson:"createdAt"`
		UpdatedAt  *time.Time `bson:"updatedAt"`
		DeletedAt  *time.Time `bson:"deletedAt"`
	}
)
