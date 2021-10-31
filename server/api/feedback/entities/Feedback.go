package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

type Feedback struct {
	ID           primitive.ObjectID `bson:"id,omitempty"`
	Message      string             `bson:"message"`
	CreationDate primitive.DateTime `bson:"creation_date"`
}
