package api

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	LoginUserInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	LoginUserOutput struct {
		User
		AuthToken string `json:"authToken"`
	}

	UpdateProfileInput struct {
		FirstName string  `json:"firstName"`
		LastName  *string `json:"lastName"`
		Avatar    *string `json:"avatar"`
		Email     string  `json:"email" validate:"email"`
		Password  string  `json:"password" validate:"min=8"`
	}

	UploadUserAvatarOutput struct {
		ImageURL string `json:"imageUrl"`
	}

	UserFunc struct {
		RegisterFunc           func(RegisterUserInput) error
		LoginFunc              func(LoginUserInput) (LoginUserOutput, error)
		ConfirmUserAccountFunc func(string) (*User, error)
		UpdateProfileFunc      func(string, UpdateProfileInput) error
		GetProfileFunc         func(string) (*User, error)
		UploadUserAvatarFunc   func(*multipart.FileHeader, string) (string, error)
	}

	UserStatus string

	User struct {
		ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
		FirstName string             `bson:"firstName,omitempty" json:"firstName"`
		LastName  *string            `bson:"lastName,omitempty" json:"lastName"`
		Username  string             `bson:"username,omitempty" json:"username"`
		Email     string             `bson:"email,omitempty" json:"email"`
		Avatar    *string            `bson:"avatar,omitempty" json:"avatar"`
		Password  string             `bson:"password,omitempty" json:"-"`
		Status    UserStatus         `bson:"status,omitempty" json:"-"`
		CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
		UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
		DeletedAt time.Time          `bson:"deletedAt,omitempty" json:"-"`
	}
)

const (
	Inactive UserStatus = "inactive"
	Active   UserStatus = "active"

	users string = "users"
)

var (
	AppName          = "talkbox"
	LoginExpDuration = time.Duration(730) * time.Hour
	JwtSigningMethod = jwt.SigningMethodHS256
	JwtSecretKey     = []byte(AppConfig.JwtSecret)
)

func (u *User) Save() error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	filter := bson.M{
		"_id": u.ID,
	}
	upsert := func() *bool {
		b := true
		return &b
	}()
	res, err := MongoDatabase.Collection(users).UpdateOne(context.Background(), filter, bson.M{"$set": u}, &options.UpdateOptions{
		Upsert: upsert,
	})
	if err != nil {
		return err
	}

	if u.ID == primitive.NilObjectID {
		u.ID = res.UpsertedID.(primitive.ObjectID)
	}

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

func FindUserByID(id string) (*User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id": objID,
	}

	var user User
	if err := MongoDatabase.Collection(users).FindOne(context.Background(), filter).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateUserToActive(id string) (*User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("[UpdateUserToActive] %v", err)
	}
	filter := bson.M{
		"_id": objID,
	}

	coll := MongoDatabase.Collection(users)

	updatedField := bson.M{
		"$set": bson.M{
			"status":    Active,
			"updatedAt": time.Now(),
		},
	}
	_, err = coll.UpdateByID(context.Background(), objID, updatedField)
	if err != nil {
		return nil, fmt.Errorf("[UpdateUserToActive] %v", err)
	}

	var user User
	if err := coll.FindOne(context.Background(), filter).Decode(&user); err != nil {
		return nil, fmt.Errorf("[UpdateUserToActive] %v", err)
	}

	return &user, nil
}

func FindUserByUsername(username string) (*User, error) {
	filter := bson.M{
		"username": username,
	}

	var user User
	if err := MongoDatabase.Collection(users).FindOne(context.Background(), filter).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateUserAvatarByID(userID, avatarURL string) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("[UpdateUserAvatarByID] %v", err)
	}

	filter := bson.M{
		"_id": objID,
	}

	_, err = MongoDatabase.Collection(users).UpdateOne(context.Background(), filter, bson.M{
		"$set": bson.M{
			"avatar": avatarURL,
		},
	})

	if err != nil {
		return fmt.Errorf("[UpdateUserAvatarByID] %v", err)
	}

	return nil
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
	fmtToken := fmt.Sprintf("%s$%s", user.ID.Hex(), token)
	encToken := base64.StdEncoding.EncodeToString([]byte(fmtToken))
	cacheKey := fmt.Sprintf("email_confirmation:%s", user.ID.Hex())
	exp := time.Duration(1) * time.Hour
	optStatus := RedisClient.Set(context.Background(), cacheKey, token, exp)
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

func GetUserProfile(userID string) (*User, error) {
	user, err := FindUserByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func UpdateUserProfile(userID string, input UpdateProfileInput) error {
	user, err := FindUserByID(userID)
	if err != nil {
		return err
	}

	user.FirstName = input.FirstName
	user.LastName = input.LastName
	user.Email = input.Email
	user.Avatar = input.Avatar

	if input.Password != "" {
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashPassword)
	}
	if err := user.Save(); err != nil {
		return err
	}
	return nil
}

func UploadUserAvatar(file *multipart.FileHeader, userID string) (string, error) {
	filename := file.Filename
	rawFile, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("[UploadUserAvatar] %v", err)
	}

	resp, err := cld.Upload.Upload(context.Background(), rawFile, uploader.UploadParams{PublicID: filename})
	if err != nil {
		return "", fmt.Errorf("[UploadUserAvatar] %v", err)
	}

	imageURL := resp.SecureURL
	if err := UpdateUserAvatarByID(userID, imageURL); err != nil {
		return "", fmt.Errorf("[UploadUserAvatar] %v", err)
	}

	return imageURL, nil
}

func UserDefaultHandler() *UserFunc {
	return &UserFunc{
		RegisterFunc:           RegisterUser,
		ConfirmUserAccountFunc: ConfirmUserAccount,
		LoginFunc:              Login,
		GetProfileFunc:         GetUserProfile,
		UpdateProfileFunc:      UpdateUserProfile,
		UploadUserAvatarFunc:   UploadUserAvatar,
	}
}

func (f *UserFunc) RegisterUserHandler(ctx *gin.Context) {
	input := RegisterUserInput{}
	if err := ctx.ShouldBind(&input); err != nil {
		log.Printf("[UserFunc.RegisterUserHander] %v", err)
		ctx.JSON(422, gin.H{
			"status":  "error",
			"message": "Failed to register the user",
		})
		return
	}

	err := f.RegisterFunc(input)
	if err != nil {
		log.Printf("[UserFunc.RegisterUserHander] %v", err)
		ctx.JSON(422, gin.H{
			"status":  "error",
			"message": "Failed to register the user",
		})
		return
	}

	ctx.JSON(201, gin.H{
		"status":  "success",
		"message": "User successfully registered, please check your email to confirm your account",
	})
}

func ConfirmUserAccount(token string) (*User, error) {
	decodedToken, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, fmt.Errorf("[ConfirmUserAccount] %v", err)
	}

	lst := strings.Split(string(decodedToken), "$")
	userID := lst[0]
	userToken := lst[1]

	cacheKey := fmt.Sprintf("email_confirmation:%s", userID)

	opts := RedisClient.Get(context.Background(), cacheKey)
	cacheToken, err := opts.Result()
	if err != nil {
		return nil, fmt.Errorf("[ConfirmUserAccount] %v", err)
	}

	if cacheToken != userToken {
		return nil, fmt.Errorf("[ConfirmUserAccount] %v", err)
	}

	user, err := UpdateUserToActive(userID)
	if err != nil {
		return nil, fmt.Errorf("[ConfirmUserAccount] %v", err)
	}

	if delOpts := RedisClient.Del(context.Background(), cacheKey); delOpts.Err() != nil {
		return nil, fmt.Errorf("[ConfirmUserAccount] %v", err)
	}

	return user, nil
}

func (f *UserFunc) ConfirmUserAccountHandler(ctx *gin.Context) {
	token := ctx.Query("token")
	if token == "" {
		ctx.JSON(422, gin.H{
			"status":  "error",
			"message": "Token is missing, please provide correct email confirmation token",
		})
		return
	}

	user, err := f.ConfirmUserAccountFunc(token)
	if err != nil {
		log.Printf("[UserFunc.ConfirmUserAccountHandler] %v", err)
		ctx.JSON(422, gin.H{
			"status":  "error",
			"message": "Failed to confirm user account, token has invalid",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"status":  "success",
		"message": "User confirmed successfully",
		"data":    user,
	})
}

func Login(input LoginUserInput) (LoginUserOutput, error) {
	user, err := FindUserByUsername(input.Username)
	if err != nil {
		return LoginUserOutput{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return LoginUserOutput{}, err
	}

	claims := struct {
		jwt.StandardClaims
		ID        string `json:"id"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Username  string `json:"username"`
		Email     string `json:"email"`
	}{
		StandardClaims: jwt.StandardClaims{
			Issuer:    AppName,
			ExpiresAt: time.Now().Add(LoginExpDuration).Unix(),
		},
		ID:        user.ID.Hex(),
		FirstName: user.FirstName,
		LastName:  *user.LastName,
		Username:  user.Username,
		Email:     user.Email,
	}

	token := jwt.NewWithClaims(JwtSigningMethod, claims)
	signedToken, err := token.SignedString(JwtSecretKey)
	if err != nil {
		return LoginUserOutput{}, err
	}

	return LoginUserOutput{
		User:      *user,
		AuthToken: signedToken,
	}, nil
}

func (f *UserFunc) LoginHandler(ctx *gin.Context) {
	input := LoginUserInput{}
	if err := ctx.ShouldBind(&input); err != nil {
		log.Printf("[UserFunc.LoginHandler] %v", err)
		ctx.JSON(400, gin.H{
			"status":  "error",
			"message": "User authentication failed, please check your request data",
		})
		return
	}

	loginOut, err := f.LoginFunc(input)
	if err != nil {
		log.Printf("[UserFunc.LoginHandler] %v", err)
		if err == mongo.ErrNoDocuments {
			ctx.JSON(422, gin.H{
				"status":  "error",
				"message": "Failed authenticate user, please check your username/password",
			})
			return
		}

		if err == bcrypt.ErrMismatchedHashAndPassword {
			ctx.JSON(422, gin.H{
				"status":  "error",
				"message": "Failed authenticate user, please check your username/password",
			})
			return
		}

		ctx.JSON(422, gin.H{
			"status":  "error",
			"message": "User authentication failed, not authorized",
		})
		return
	}

	splittedToken := strings.Split(loginOut.AuthToken, ".")
	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie("talkbox", strings.Join([]string{splittedToken[0], splittedToken[1]}, "."), 3600, "/", "", true, false)
	ctx.SetCookie("talkbox_sign", splittedToken[2], 3600, "/", "", true, true)

	ctx.JSON(200, gin.H{
		"status":  "success",
		"message": "User authenticated",
		"data":    loginOut,
	})
}

func (f *UserFunc) GetProfileHandler(ctx *gin.Context) {
	userCtx, ok := ctx.Get("user")
	if !ok {
		log.Println("[GetProfileHandler] Unable to get current user")
		ctx.JSON(422, gin.H{
			"status":   "error",
			"messages": "Failed to get user profile",
		})
		return
	}
	user := userCtx.(*User)

	userProfile, err := f.GetProfileFunc(user.ID.Hex())
	if err != nil {
		log.Printf("[GetProfileHandler] %v", err)
		ctx.JSON(422, gin.H{
			"status":   "error",
			"messages": "Failed to get user profile",
		})
	}

	ctx.JSON(200, gin.H{
		"status":  "success",
		"message": "Successfully get user profile",
		"data":    userProfile,
	})
}

func (f *UserFunc) UpdateProfileHandler(ctx *gin.Context) {
	userCtx, ok := ctx.Get("user")
	if !ok {
		log.Println("[GetProfileHandler] Unable to get current user")
		ctx.JSON(422, gin.H{
			"status":  "error",
			"message": "Failed to get user profile",
		})
		return
	}
	user := userCtx.(*User)

	input := UpdateProfileInput{}
	if err := ctx.ShouldBind(&input); err != nil {
		log.Printf("[UpdateProfileHandler] %v", err)
		ctx.JSON(400, gin.H{
			"status":  "error",
			"message": "Failed to get user profile, please check your request data",
		})
		return
	}

	err := f.UpdateProfileFunc(user.ID.Hex(), input)
	if err != nil {
		log.Printf("[UpdateProfileHandler] %v", err)
		ctx.JSON(422, gin.H{
			"status":  "error",
			"message": "Failed to update user profile",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"status":  "success",
		"message": "Successfully update user profile",
	})
}

func (f *UserFunc) UploadUserAvatarHandler(ctx *gin.Context) {
	userCtx, ok := ctx.Get("user")
	if !ok {
		log.Println("[UploadUserAvatarHandler] Unable to get current user")
		ctx.JSON(422, gin.H{
			"status":  "error",
			"message": "Failed to get user profile",
		})
		return
	}
	user := userCtx.(*User)

	file, err := ctx.FormFile("avatar")
	if err != nil {
		log.Printf("[UploadUserAvatarHandler] %v", err)
		ctx.JSON(422, gin.H{
			"status":  "error",
			"message": "Unable to get user avatar file",
		})
		return
	}

	userID := user.ID.Hex()
	url, err := f.UploadUserAvatarFunc(file, userID)
	if err != nil {
		log.Printf("[UploadUserAvatarHandler] %v", err)
		ctx.JSON(422, gin.H{
			"status":  "error",
			"message": "Failed to upload user avatar, please try again later",
		})
		return
	}

	res := UploadUserAvatarOutput{
		ImageURL: url,
	}

	ctx.JSON(200, gin.H{
		"status":  "success",
		"message": "Successfully upload user avatar",
		"data":    res,
	})
}
