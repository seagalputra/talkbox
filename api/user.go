package api

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type (
	RegisterUserInput struct {
		FirstName            string  `json:"firstName" validate:"required"`
		LastName             *string `json:"lastName"`
		Username             string  `json:"username" validate:"required"`
		Email                string  `json:"email" validate:"required,email"`
		Password             string  `json:"password" validate:"required,min=8"`
		PasswordConfirmation string  `json:"passwordConfirmation" validate:"required"`
	}

	RegisterUserOutput struct {
		User
	}

	UserFunc struct {
		RegisterFunc func(RegisterUserInput) *RegisterUserOutput
	}

	UserStatus string

	User struct {
		ID        primitive.ObjectID `bson:"_id" json:"id"`
		FirstName string             `bson:"firstName" json:"firstName"`
		LastName  *string            `bson:"lastName" json:"lastName"`
		Username  string             `bson:"username" json:"username"`
		Email     string             `bson:"email" json:"email"`
		Avatar    *string            `bson:"avatar" json:"avatar"`
		Password  string             `bson:"password" json:"-"`
		Status    UserStatus         `bson:"status" json:"-"`
		CreatedAt *time.Time         `bson:"createdAt" json:"createdAt"`
		UpdatedAt *time.Time         `bson:"updatedAt" json:"updatedAt"`
		DeletedAt *time.Time         `bson:"deletedAt" json:"deletedAt"`
	}
)

const (
	Inactive UserStatus = "inactive"
	Active   UserStatus = "active"

	users string = "users"
)

func (u *User) Save() error {
	res, err := MongoDatabase.Collection(users).InsertOne(context.TODO(), u)
	if err != nil {
		return err
	}
	u.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (u *User) IsAvailable() (bool, error) {
	filter := bson.M{
		"email":    u.Email,
		"username": u.Username,
	}

	var result bson.M
	err := MongoDatabase.Collection(users).FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return false, err
	}

	if result["email"] != nil || result["username"] != nil {
		return false, nil
	}

	return true, nil
}

func RegisterUser(input RegisterUserInput) error {
	user := &User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Username:  input.Username,
		Email:     input.Email,
	}

	isAvailable, err := user.IsAvailable()
	if err != nil {
		log.Printf("[RegisterUser] %v", err)
		return err
	}

	if !isAvailable {
		return errors.New("user already registered, please use other email/username")
	}

	if input.Password != input.PasswordConfirmation {
		return errors.New("your password and confirmation password doesn't match")
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[RegisterUser] %v", err)
		return err
	}

	user.Password = string(encryptedPassword)
	if err := user.Save(); err != nil {
		log.Printf("[RegisterUser] %v", err)
		return err
	}

	token := GenRandString(20)
	fmtToken := fmt.Sprintf("%s$%s", user.ID, token)
	encToken := base64.StdEncoding.EncodeToString([]byte(fmtToken))
	cacheKey := fmt.Sprintf("email_confirmation:%s", user.ID)
	exp := time.Duration(1) * time.Hour
	optStatus := RedisClient.Set(context.Background(), cacheKey, encToken, exp)
	if err := optStatus.Err(); err != nil {
		log.Printf("[RegisterUser] %v", err)
		return err
	}
	// TODO: send email as confirmation email
	return nil
}
