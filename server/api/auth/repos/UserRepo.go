package repos

import (
	"github.com/LastSprint/GooodBack/api/auth/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepo struct {
}

func (u UserRepo) GetUserById(userOauthId, authProvider string) (*entities.UserInfo, error) {
	return nil, nil
}

func (u UserRepo) CreateUser(userOauthId, authProvider string) (*entities.UserInfo, error) {
	return &entities.UserInfo{
		RefreshToken:  "",
		ID:            primitive.NewObjectID(),
		OauthProvider: "google",
	}, nil
}

func (u UserRepo) SetRefreshTokenForUser(userId string) error {
	return nil
}
