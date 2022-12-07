package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	SendMessageInput struct {
		Body       string  `json:"body"`
		Attachment *string `json:"attachment"`
	}

	SendMessageOutput struct {
		ID         string    `json:"id"`
		Body       string    `json:"body"`
		Attachment *string   `json:"attachment"`
		UserID     string    `json:"userId"`
		RoomID     string    `json:"roomId"`
		User       *User     `json:"user"`
		Room       *Room     `json:"room"`
		CreatedAt  time.Time `json:"createdAt"`
		UpdatedAt  time.Time `json:"updatedAt"`
	}

	GetMessagesInput struct {
		RoomID string
		Cursor string
		Limit  int64
	}

	GetMessageOutput struct {
		Cursor   string
		Limit    int64
		Messages []Message
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
		ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
		Body       string             `bson:"body" json:"body"`
		Attachment *string            `bson:"attachment" json:"attachment"`
		UserID     primitive.ObjectID `bson:"userId,omitempty" json:"userId"`
		RoomID     primitive.ObjectID `bson:"roomId,omitempty" json:"roomId"`
		CreatedAt  time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
		UpdatedAt  time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
		DeletedAt  time.Time          `bson:"deletedAt,omitempty" json:"-"`
		User       *User              `bson:"user,omitempty" json:"user"`
		Room       *Room              `bson:"room,omitempty" json:"room"`
	}

	MessageFunc struct {
		GetMessagesFunc func(GetMessagesInput) GetMessageOutput
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

func FindByRoomID(roomID string, cursorObj map[string]interface{}, limit int64) ([]Message, error) {
	objID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return []Message{}, err
	}

	pipeline := mongo.Pipeline{
		bson.D{{"$match", bson.D{{"roomId", bson.D{{"$eq", objID}}}}}},
		bson.D{{"$sort", bson.D{{"updatedAt", 1}, {"_id", 1}}}},
	}

	if len(cursorObj) != 0 {
		id := cursorObj["id"].(string)
		updatedAt, err := time.Parse(time.RFC3339, cursorObj["updatedAt"].(string))
		if err != nil {
			return []Message{}, err
		}

		msgObjID, err := primitive.ObjectIDFromHex(id)
		timeObj := primitive.NewDateTimeFromTime(updatedAt)
		if err != nil {
			return []Message{}, err
		}

		pipeline = append(pipeline, bson.D{
			{"$match",
				bson.D{
					{"updatedAt", bson.D{{"$gt", timeObj}}},
					{"_id", bson.D{{"$gt", msgObjID}}},
				},
			},
		})
	}

	pipeline = append(pipeline,
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "users"},
					{"localField", "userId"},
					{"foreignField", "_id"},
					{"as", "user"},
				},
			},
		},
		bson.D{{"$unwind", bson.D{{"path", "$user"}}}},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "rooms"},
					{"localField", "roomId"},
					{"foreignField", "_id"},
					{"as", "room"},
				},
			},
		},
		bson.D{{"$unwind", bson.D{{"path", "$room"}}}},
		bson.D{{"$limit", limit}})

	cursor, err := MongoDatabase.Collection(messages).Aggregate(context.Background(), pipeline)
	if err != nil {
		return []Message{}, err
	}

	var messages = make([]Message, 0)
	for cursor.Next(context.Background()) {
		var message Message
		if err := cursor.Decode(&message); err != nil {
			return []Message{}, err
		}
		messages = append(messages, message)
	}
	if err := cursor.Err(); err != nil {
		return []Message{}, err
	}

	return messages, nil
}

func GetMessages(input GetMessagesInput) GetMessageOutput {
	cursorInput := make(map[string]interface{})
	if input.Cursor != "" {
		d, err := base64.StdEncoding.DecodeString(input.Cursor)
		if err != nil {
			log.Printf("[GetMessages] %v", err)
			return GetMessageOutput{
				Limit:    input.Limit,
				Messages: []Message{},
			}
		}
		err = json.Unmarshal(d, &cursorInput)
		if err != nil {
			log.Printf("[GetMessages] %v", err)
			return GetMessageOutput{
				Limit:    input.Limit,
				Messages: []Message{},
			}
		}
	}

	messages, err := FindByRoomID(input.RoomID, cursorInput, input.Limit)
	if err != nil {
		log.Printf("[GetMessages] %v", err)
		return GetMessageOutput{
			Limit:    input.Limit,
			Messages: []Message{},
		}
	}

	var cursor string
	if len(messages) != 0 {
		lastMessage := messages[len(messages)-1]
		c := map[string]interface{}{
			"id":        lastMessage.ID.Hex(),
			"updatedAt": lastMessage.UpdatedAt,
		}
		jsonCursor, err := json.Marshal(c)
		if err != nil {
			log.Printf("[GetMessages] %v", err)
			return GetMessageOutput{
				Limit:    input.Limit,
				Messages: []Message{},
			}
		}

		cursor = base64.StdEncoding.EncodeToString(jsonCursor)
	}

	output := GetMessageOutput{
		Cursor:   cursor,
		Limit:    input.Limit,
		Messages: messages,
	}

	return output
}

func MessageDefaultHandler() *MessageFunc {
	return &MessageFunc{
		GetMessagesFunc: GetMessages,
	}
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

		roomObjID, err := primitive.ObjectIDFromHex(roomID)
		if err != nil {
			log.Printf("[WSHandler] %v", err)
			continue
		}

		userObjID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			log.Printf("[WSHandler] %v", err)
			continue
		}
		message := Message{
			Body:       input.Body,
			Attachment: input.Attachment,
			RoomID:     roomObjID,
			UserID:     userObjID,
		}
		if err := message.Save(); err != nil {
			log.Printf("[WSHandler] %v", err)
		}

		recipientConn, present := userConnection[recipientID]
		if present {
			recipientConn.WriteJSON(SendMessageOutput{
				ID:         message.ID.Hex(),
				Body:       input.Body,
				Attachment: input.Attachment,
				UserID:     user.ID.Hex(),
				RoomID:     room.ID.Hex(),
				User:       user,
				Room:       room,
				CreatedAt:  message.CreatedAt,
				UpdatedAt:  message.UpdatedAt,
			})
		}
	}
}

func (f *MessageFunc) GetMessagesHandler(ctx *gin.Context) {
	roomID := ctx.Param("room_id")
	cursor := ctx.Query("cursor")
	limitQuery := ctx.Query("limit")

	limit := 10
	if limitQuery != "" {
		var err error
		limit, err = strconv.Atoi(limitQuery)
		if err != nil {
			log.Printf("[GetMessagesHandler] %v", err)
			ctx.JSON(422, gin.H{
				"status":  "error",
				"message": "Failed to get messages",
			})
			return
		}
	}

	input := GetMessagesInput{
		RoomID: roomID,
		Cursor: cursor,
		Limit:  int64(limit),
	}
	output := f.GetMessagesFunc(input)

	ctx.JSON(200, gin.H{
		"status": "success",
		"meta": gin.H{
			"limit":  input.Limit,
			"cursor": output.Cursor,
		},
		"data": output.Messages,
	})
}
