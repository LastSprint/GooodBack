package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserInfo struct {
	RefreshToken  string             `bson:"refresh_token,omitempty"`
	ID            primitive.ObjectID `bson:"id,omitempty"`
	OauthProvider string             `bson:"oauth_provider"`
}
