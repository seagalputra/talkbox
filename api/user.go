package api

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
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
		RegisterFunc func(RegisterUserInput) error
	}

	UserStatus string

	User struct {
		ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
		FirstName string             `bson:"firstName" json:"firstName"`
		LastName  *string            `bson:"lastName" json:"lastName"`
		Username  string             `bson:"username" json:"username"`
		Email     string             `bson:"email" json:"email"`
		Avatar    *string            `bson:"avatar" json:"avatar"`
		Password  string             `bson:"password" json:"-"`
		Status    UserStatus         `bson:"status" json:"-"`
		CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
		UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
		DeletedAt time.Time          `bson:"deletedAt,omitempty" json:"deletedAt"`
	}
)

const (
	Inactive UserStatus = "inactive"
	Active   UserStatus = "active"

	users string = "users"
)

func (u *User) Save() error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
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

	count, err := MongoDatabase.Collection(users).CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}

	if count > 0 {
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
		Status:    Inactive,
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
	fmtToken := fmt.Sprintf("%s$%s", user.ID.String(), token)
	encToken := base64.StdEncoding.EncodeToString([]byte(fmtToken))
	cacheKey := fmt.Sprintf("email_confirmation:%s", user.ID.String())
	exp := time.Duration(1) * time.Hour
	optStatus := RedisClient.Set(context.Background(), cacheKey, encToken, exp)
	if err := optStatus.Err(); err != nil {
		log.Printf("[RegisterUser] %v", err)
		return err
	}
	go sendConfirmationEmail(user.Email, encToken)
	return nil
}

func sendConfirmationEmail(to, token string) {
	urlVal := url.Values{}
	urlVal.Set("token", token)
	// TODO: move body email to html file
	body := fmt.Sprintf(`
	<div>
		<p>Click link below to verify your account</p>
		<p>%s</p>
	</div>
	`, AppConfig.EmailConfirmationURL+"?"+urlVal.Encode())

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", AppConfig.EmailSenderName)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", "Verify your account - Talkbox")
	mailer.SetBody("text/html", body)

	smtpPort, err := strconv.Atoi(AppConfig.SMTPPort)
	if err != nil {
		log.Printf("[sendConfirmationEmail] %v", err)
		return
	}

	dialer := gomail.NewDialer(
		AppConfig.SMTPHost,
		smtpPort,
		AppConfig.SMTPUsername,
		AppConfig.SMTPPassword,
	)

	if err = dialer.DialAndSend(mailer); err != nil {
		log.Printf("[sendConfirmationEmail] %v", err)
		return
	}
	log.Printf("[sendConfirmationEmail] Confirmation email successfully sent to %s", to)
}

func (f *UserFunc) RegisterUserHandler(ctx *gin.Context) {
	input := RegisterUserInput{}
	if err := ctx.ShouldBind(&input); err != nil {
		log.Printf("[UserFunc.RegisterUserHander] %v", err)
		ctx.JSON(422, gin.H{
			"status":  "error",
			"message": "failed to register the user",
		})
		return
	}

	err := f.RegisterFunc(input)
	if err != nil {
		log.Printf("[UserFunc.RegisterUserHander] %v", err)
		ctx.JSON(422, gin.H{
			"status":  "error",
			"message": "failed to register the user",
		})
		return
	}

	ctx.JSON(201, gin.H{
		"status":  "success",
		"message": "User successfully registered, please check your email to confirm your account",
	})
}
