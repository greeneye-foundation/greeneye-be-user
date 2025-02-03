package services

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/greeneye-foundation/greeneye-be-user/internal/models"
	"github.com/greeneye-foundation/greeneye-be-user/internal/pkg/utils"
)

type UserService struct {
	collection *mongo.Collection
}

func NewUserService(client *mongo.Client, dbName string) *UserService {
	return &UserService{
		collection: client.Database(dbName).Collection("users"),
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	// Hash password
	hashedPassword, err := utils.HashPassword(user.PasswordHash)
	if err != nil {
		return err
	}
	user.PasswordHash = hashedPassword

	// Insert user
	_, err = s.collection.InsertOne(ctx, user)
	return err
}

func (s *UserService) AuthenticateUser(ctx context.Context, login *models.UserLogin) (*models.User, error) {
	var user models.User
	err := s.collection.FindOne(ctx, bson.M{"mobile_number": login.MobileNumber}).Decode(&user)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Compare passwords
	if err := utils.CheckPasswordHash(login.Password, user.PasswordHash); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}

func (s *UserService) GetUserByMobileNumber(ctx context.Context, mobileNumber string) (*models.User, error) {
	var user models.User
	err := s.collection.FindOne(ctx, bson.M{"mobile_number": mobileNumber}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, user *models.User) error {
	_, err := s.collection.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		bson.M{
			"$set": bson.M{
				"password_hash": user.PasswordHash,
				"updated_at":    user.UpdatedAt,
				// Add other fields as necessary
			},
		},
	)
	return err
}

func (s *UserService) GetUsers(ctx context.Context) ([]*models.User, error) {
	collection := s.collection.Database().Collection("users")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}
