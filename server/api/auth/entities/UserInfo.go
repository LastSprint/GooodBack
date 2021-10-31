package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserInfo struct {
	RefreshToken       string             `bson:"refresh_token,omitempty"`
	ID                 primitive.ObjectID `bson:"_id,omitempty"`
	ThirdPartyProvider string             `bson:"third_party_provider"`
	ThirdPartyId       string             `bson:"third_party_id"`
}
