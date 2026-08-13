package data

import (
	"context"
	"errors"
	"strings"
	"time"

	"auth-Task-Manager/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	collection *mongo.Collection
}

func NewUserService(db *mongo.Database) *UserService {
	return &UserService{
		collection: db.Collection("users"),
	}
}

// CreateUser creates a new user.
func (service *UserService) CreateUser(
	username string,
	password string,
) (models.User, error) {

	username = strings.TrimSpace(username)

	if username == "" {
		return models.User{}, errors.New("username is required")
	}

	if len(password) < 8 {
		return models.User{}, errors.New(
			"password must be at least 8 characters",
		)
	}

	// Check whether username already exists.
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	var existingUser models.User

	err := service.collection.FindOne(
		ctx,
		bson.M{
			"username": username,
		},
	).Decode(&existingUser)

	if err == nil {
		return models.User{}, errors.New(
			"username already exists",
		)
	}

	if !errors.Is(err, mongo.ErrNoDocuments) {
		return models.User{}, err
	}

	// Hash password.
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return models.User{}, errors.New(
			"failed to hash password",
		)
	}

	// Determine role.
	userCount, err := service.collection.CountDocuments(
		ctx,
		bson.M{},
	)

	if err != nil {
		return models.User{}, err
	}

	role := "user"

	// First user becomes admin.
	if userCount == 0 {
		role = "admin"
	}

	user := models.User{
		ID:       bson.NewObjectID().Hex(),
		Username: username,
		Password: string(hashedPassword),
		Role:     role,
	}

	_, err = service.collection.InsertOne(
		ctx,
		user,
	)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// GetUserByUsername finds a user.
func (service *UserService) GetUserByUsername(
	username string,
) (models.User, error) {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	var user models.User

	err := service.collection.FindOne(
		ctx,
		bson.M{
			"username": username,
		},
	).Decode(&user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.User{}, errors.New(
				"user not found",
			)
		}

		return models.User{}, err
	}

	return user, nil
}

// GetUserByID finds a user by ID.
func (service *UserService) GetUserByID(
	id string,
) (models.User, error) {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	var user models.User

	err := service.collection.FindOne(
		ctx,
		bson.M{
			"_id": id,
		},
	).Decode(&user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.User{}, errors.New(
				"user not found",
			)
		}

		return models.User{}, err
	}

	return user, nil
}

// CheckPassword checks the plaintext password
// against the bcrypt hash stored in MongoDB.
func (service *UserService) CheckPassword(
	user models.User,
	password string,
) error {

	return bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)
}