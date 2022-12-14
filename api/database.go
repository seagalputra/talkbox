package api

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDatabase *mongo.Database

func ConnectDatabase() error {
	url := AppConfig.DatabaseURL
	if url == "" {
		return errors.New("unable to connect the database, you should fill the DATABASE_URL in the config")
	}
	MongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(url))
	if err != nil {
		return err
	}

	dbName := AppConfig.DatabaseName
	MongoDatabase = MongoClient.Database(dbName)

	return nil
}

var RedisClient *redis.Client

func ConnectToRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: AppConfig.RedisHost,
	})
}
