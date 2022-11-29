package api

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	SendMessageInput struct {
		Body       string  `json:"body"`
		Attachment *string `json:"attachment"`
	}

	SendMessageOutput struct {
		Body       string    `json:"body"`
		Attachment *string   `json:"attachment"`
		From       *User     `json:"from"`
		SentAt     time.Time `json:"sentAt"`
	}

	WebSocketConnection struct {
		*websocket.Conn
		UserID   string
		Username string
		Email    string
		RoomID   string
	}

	// key using user id
	UserConnection map[string]WebSocketConnection

	Message struct {
		ID         primitive.ObjectID `bson:"_id,omitempty"`
		Body       string             `bson:"body"`
		Attachment *string            `bson:"attachment"`
		UserID     string             `bson:"userId"`
		RoomID     string             `bson:"roomId"`
		CreatedAt  time.Time          `bson:"createdAt,omitempty"`
		UpdatedAt  time.Time          `bson:"updatedAt,omitempty"`
		DeletedAt  time.Time          `bson:"deletedAt,omitempty"`
	}

	MessageFunc struct {
	}
)

const (
	messages string = "messages"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: should check the origin request from client
		return true
	},
}

var userConnection = make(UserConnection)

func (m *Message) Save() error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	res, err := MongoDatabase.Collection(messages).InsertOne(context.TODO(), m)
	if err != nil {
		return err
	}
	m.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func MessageDefaultHandler() *MessageFunc {
	return &MessageFunc{}
}

func (f *MessageFunc) WSHandler(ctx *gin.Context) {
	conn, err := wsUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("wsHandler: %v", err)
		return
	}

	userCtx, ok := ctx.Get("user")
	if !ok {
		return
	}
	user := userCtx.(*User)

	userID := user.ID.Hex()
	roomID := ctx.Param("room_id")
	wsConn := WebSocketConnection{
		Conn:     conn,
		UserID:   userID,
		Username: user.Username,
		Email:    user.Email,
		RoomID:   roomID,
	}
	userConnection[userID] = wsConn

	for {
		var input SendMessageInput
		err := conn.ReadJSON(&input)
		if err != nil {
			if strings.Contains(err.Error(), "websocket: close") {
				delete(userConnection, userID)
				return
			}

			log.Printf("[WSHandler] %v", err)
			continue
		}
		log.Printf("%v", input)

		// check if recipient has active connection, if not save the message and send notification
		room, err := FindRoomByID(roomID)
		if err != nil {
			log.Printf("[WSHandler] %v", err)
			if err == mongo.ErrNoDocuments {
				return
			}
			continue
		}

		var recipientID string
		for _, participant := range room.Participants {
			if participant.ID.Hex() != userID {
				recipientID = participant.ID.Hex()
			}
		}
		log.Printf("recipient: %s", recipientID)

		recipientConn, present := userConnection[recipientID]
		if present {
			recipientConn.WriteJSON(SendMessageOutput{
				Body:       input.Body,
				Attachment: input.Attachment,
				From:       user,
				SentAt:     time.Now(),
			})
		}

		message := Message{
			Body:       input.Body,
			Attachment: input.Attachment,
			RoomID:     roomID,
			UserID:     userID,
		}
		if err := message.Save(); err != nil {
			log.Printf("[WSHandler] %v", err)
		}
	}
}
