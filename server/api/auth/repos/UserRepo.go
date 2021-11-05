package repos

import (
	"context"
	"errors"
	"fmt"
	"github.com/LastSprint/GooodBack/api/auth/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	Client *mongo.Client
}

func (u *UserRepo) CheckIfUserAllowedToAuth(provider, userId string) (bool, error) {
	col := u.Client.Database("feedback_auth").Collection("allowed_users")

	search := col.FindOne(context.Background(), bson.M{
		"providerId": provider,
		"targets": bson.M{
			"$in": bson.A{userId},
		},
	})

	err := search.Err()

	if errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (u *UserRepo) GetUserIdByTarget(target string) (string, error) {
	col := u.Client.Database("feedback_auth").Collection("users")

	search := col.FindOne(context.Background(), bson.M{
		"third_party_id": target,
	})

	err := search.Err()

	if errors.Is(err, mongo.ErrNoDocuments) {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	result := &entities.UserInfo{}

	if err = search.Decode(result); err != nil {
		return "", err
	}

	return result.ID.Hex(), nil
}

func (u *UserRepo) GetUserById(userOauthId, authProvider string) (*entities.UserInfo, error) {
	col := u.Client.Database("feedback_auth").Collection("users")

	search := col.FindOne(context.Background(), bson.M{
		"third_party_id":       userOauthId,
		"third_party_provider": authProvider,
	})

	err := search.Err()

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	result := &entities.UserInfo{}

	if err := search.Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

func (u *UserRepo) CreateUser(userOauthId, authProvider string) (*entities.UserInfo, error) {
	col := u.Client.Database("feedback_auth").Collection("users")

	entity := entities.UserInfo{
		ThirdPartyProvider: authProvider,
		ThirdPartyId:       userOauthId,
	}

	insert, err := col.InsertOne(context.TODO(), entity)

	if err != nil {
		return nil, err
	}

	id, ok := insert.InsertedID.(primitive.ObjectID)

	if !ok {
		return nil, fmt.Errorf("couldn't parse id from %v", insert.InsertedID)
	}

	entity.ID = id

	return &entity, nil
}

func (u *UserRepo) SetRefreshTokenForUser(userId, token string) error {

	id, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return fmt.Errorf("couldnt convert userId to  ObjectID with error %w", err)
	}

	col := u.Client.Database("feedback_auth").Collection("users")

	update, err := col.UpdateOne(context.Background(), bson.M{
		"_id": id,
	}, bson.M{
		"$set": bson.M{
			"refresh_token": token,
		},
	})

	if err != nil {
		return err
	}

	if update.MatchedCount == 0 {
		return fmt.Errorf("user with id %s not found", userId)
	}

	return nil
}
