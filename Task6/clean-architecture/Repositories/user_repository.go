package repositories

import (
	"context"
	"time"

	domain "clean-architecture/Domain"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository interface {
	CreateUser(user domain.User) error
	GetUsers() ([]domain.User, error)
	GetUserByUsername(username string) (domain.User, error)
	GetUserByID(id string) (domain.User, error)
}

type UserDocument struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	Firstname string        `bson:"firstname"`
	Lastname  string        `bson:"lastname"`
	Username  string        `bson:"username"`
	Password  string        `bson:"password"`
	Role      string        `bson:"role"`
}

type UserRepositoryImpl struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		collection: db.Collection("user"),
	}
}



func (service *UserRepositoryImpl) CreateUser(
	user domain.User,
) error {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	newUser := UserDocument{
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Username:  user.Username,
		Password:  user.Password,
		Role:      user.Role,
	}

	_, err := service.collection.InsertOne(
		ctx,
		newUser,
	)

	return err
}

func (service *UserRepositoryImpl) GetUsers() ([]domain.User, error) {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	cursor, err := service.collection.Find(
		ctx,
		bson.M{},
	)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var documents []UserDocument

	if err := cursor.All(ctx, &documents); err != nil {
		return nil, err
	}

	users := make([]domain.User, 0, len(documents))

	for _, document := range documents {

		user := domain.User{
			ID:        document.ID.Hex(),
			Firstname: document.Firstname,
			Lastname:  document.Lastname,
			Username:  document.Username,
			Password:  document.Password,
			Role:      document.Role,
		}

		users = append(users, user)
	}

	return users, nil
}

func (service *UserRepositoryImpl) GetUserByUsername(
	username string,
) (domain.User, error) {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	var document UserDocument

	err := service.collection.FindOne(
		ctx,
		bson.M{
			"username": username,
		},
	).Decode(&document)

	if err != nil {
		return domain.User{}, err
	}

	user := domain.User{
		ID:        document.ID.Hex(),
		Firstname: document.Firstname,
		Lastname:  document.Lastname,
		Username:  document.Username,
		Password:  document.Password,
		Role:      document.Role,
	}

	return user, nil
}

func (service *UserRepositoryImpl) GetUserByID(
	id string,
) (domain.User, error) {

	objectID, err := bson.ObjectIDFromHex(id)

	if err != nil {
		return domain.User{}, err
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	var document UserDocument

	err = service.collection.FindOne(
		ctx,
		bson.M{
			"_id": objectID,
		},
	).Decode(&document)

	if err != nil {
		return domain.User{}, err
	}

	user := domain.User{
		ID:        document.ID.Hex(),
		Firstname: document.Firstname,
		Lastname:  document.Lastname,
		Username:  document.Username,
		Password:  document.Password,
		Role:      document.Role,
	}

	return user, nil
}
