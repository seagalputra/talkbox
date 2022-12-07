package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	GetRoomsInput struct {
		Cursor string
		Limit  int64
		User   *User
	}

	GetRoomsOutput struct {
		Cursor string
		Limit  int64
		Rooms  []Room
	}

	RoomType string

	Participant struct {
		ID       primitive.ObjectID `bson:"id,omitempty" json:"id"`
		Username string             `bson:"username" json:"username"`
		Email    string             `bson:"email" json:"email"`
		Avatar   string             `bson:"avatar" json:"avatar"`
	}

	Room struct {
		ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
		Participants []Participant      `bson:"participants" json:"participants"`
		RoomType     RoomType           `bson:"roomType" json:"roomType"`
		LastMessage  string             `bson:"lastMessage" json:"lastMessage"`
		CreatedAt    *time.Time         `bson:"createdAt,omitempty" json:"createdAt"`
		UpdatedAt    *time.Time         `bson:"updatedAt,omitempty" json:"updatedAt"`
		DeletedAt    *time.Time         `bson:"deletedAt,omitempty" json:"-"`
	}

	RoomFunc struct {
		GetRoomsFunc func(GetRoomsInput) GetRoomsOutput
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

func FindRoomsByUserID(userID *primitive.ObjectID, cursorObj map[string]interface{}, limit int64) ([]Room, error) {
	pipeline := mongo.Pipeline{
		bson.D{{"$sort", bson.D{{"updatedAt", -1}, {"_id", -1}}}},
	}

	if userID != nil {
		// it prepend to the pipline array
		pipeline = append(mongo.Pipeline{
			bson.D{{"$match", bson.D{{"participants", bson.D{{"$elemMatch", bson.D{{"id", userID}}}}}}}},
		}, pipeline...)
	}

	if len(cursorObj) != 0 {
		id := cursorObj["id"].(string)
		updatedAt, err := time.Parse(time.RFC3339, cursorObj["updatedAt"].(string))
		if err != nil {
			return []Room{}, err
		}

		roomObjID, err := primitive.ObjectIDFromHex(id)
		timeObj := primitive.NewDateTimeFromTime(updatedAt)
		if err != nil {
			return []Room{}, err
		}

		pipeline = append(pipeline, bson.D{
			{
				"$match",
				bson.D{
					{"updatedAt", bson.D{{"$gt", timeObj}}},
					{"_id", bson.D{{"$gt", roomObjID}}},
				},
			},
		})
	}

	pipeline = append(pipeline,
		bson.D{
			{"$limit", limit},
		})

	cursor, err := MongoDatabase.Collection(rooms).Aggregate(context.Background(), pipeline)
	if err != nil {
		return []Room{}, err
	}

	var rooms = make([]Room, 0)
	for cursor.Next(context.Background()) {
		var room Room
		if err := cursor.Decode(&room); err != nil {
			return []Room{}, err
		}
		rooms = append(rooms, room)
	}

	if err := cursor.Err(); err != nil {
		return []Room{}, err
	}

	return rooms, nil
}

func RoomDefaultHandler() *RoomFunc {
	return &RoomFunc{
		GetRoomsFunc: GetRooms,
	}
}

func GetRooms(input GetRoomsInput) GetRoomsOutput {
	cursorInput := make(map[string]interface{})
	if input.Cursor != "" {
		decoded, err := base64.StdEncoding.DecodeString(input.Cursor)
		if err != nil {
			log.Printf("[GetRooms] %v", err)
			return GetRoomsOutput{
				Limit: input.Limit,
				Rooms: []Room{},
			}
		}

		err = json.Unmarshal(decoded, &cursorInput)
		if err != nil {
			log.Printf("[GetRooms] %v", err)
			return GetRoomsOutput{
				Limit: input.Limit,
				Rooms: []Room{},
			}
		}
	}

	var userID *primitive.ObjectID
	if input.User != nil {
		userID = &input.User.ID
	}

	rooms, err := FindRoomsByUserID(userID, cursorInput, input.Limit)
	if err != nil {
		log.Printf("[GetRooms] %v", err)
		return GetRoomsOutput{
			Limit: input.Limit,
			Rooms: []Room{},
		}
	}

	var cursor string
	if len(rooms) != 0 {
		lastRoom := rooms[len(rooms)-1]
		c := map[string]interface{}{
			"id":        lastRoom.ID.Hex(),
			"updatedAt": lastRoom.UpdatedAt,
		}
		jsonCursor, err := json.Marshal(c)
		if err != nil {
			log.Printf("[GetRooms] %v", err)
			return GetRoomsOutput{
				Limit: input.Limit,
				Rooms: []Room{},
			}
		}

		cursor = base64.StdEncoding.EncodeToString(jsonCursor)
	}

	output := GetRoomsOutput{
		Cursor: cursor,
		Limit:  input.Limit,
		Rooms:  rooms,
	}

	return output
}

func (f *RoomFunc) GetRoomsHandler(ctx *gin.Context) {
	cursor := ctx.Query("cursor")
	limitQuery := ctx.Query("limit")

	userCtx, ok := ctx.Get("user")
	if !ok {
		log.Println("[GetRoomsHandler] Unable to get current user")
		ctx.JSON(422, gin.H{
			"status":   "error",
			"messages": "Failed to get rooms",
		})
	}
	user := userCtx.(*User)

	limit := 10
	if limitQuery != "" {
		var err error
		limit, err = strconv.Atoi(limitQuery)
		if err != nil {
			log.Printf("[GetRoomsHandler] %v", err)
			ctx.JSON(422, gin.H{
				"status":   "error",
				"messages": "Failed to get rooms",
			})
		}
	}

	input := GetRoomsInput{
		Cursor: cursor,
		Limit:  int64(limit),
		User:   user,
	}
	output := f.GetRoomsFunc(input)

	ctx.JSON(200, gin.H{
		"status": "success",
		"meta": gin.H{
			"limit":  input.Limit,
			"cursor": output.Cursor,
		},
		"data": output.Rooms,
	})
}
