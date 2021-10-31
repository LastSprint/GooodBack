package repos

import (
	"context"
	"fmt"
	"github.com/LastSprint/GooodBack/api/feedback/entities"
	"github.com/LastSprint/GooodBack/api/feedback/entries"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sort"
	"time"
)

type FeedbackRepository struct {
	Client *mongo.Client
}

func (f *FeedbackRepository) CreateNew(message, userId string) error {
	db := f.Client.Database("feedback")
	col := db.Collection("feedbacks")

	id, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return fmt.Errorf("couldnt convert userId to  ObjectID with error %w", err)
	}

	entry := entities.Feedback{
		ID:           primitive.NewObjectID(),
		Message:      message,
		CreationDate: primitive.NewDateTimeFromTime(time.Now()),
	}

	opts := options.Update().SetUpsert(true)

	_, err = col.UpdateOne(context.TODO(), bson.M{"userId": id}, bson.M{"$push": bson.M{"feedbacks": entry}}, opts)

	return err
}

func (f *FeedbackRepository) Read(userId string) ([]entries.Feedback, error) {

	id, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return nil, fmt.Errorf("couldnt convert userId to  ObjectID with error %w", err)
	}

	col := f.Client.Database("feedback").Collection("feedbacks")

	result := col.FindOne(context.TODO(), bson.M{"userId": id})

	if err = result.Err(); err != nil {
		return nil, err
	}

	found := struct {
		Arr []entities.Feedback `bson:"feedbacks"`
	}{}

	if err = result.Decode(&found); err != nil {
		return nil, err
	}

	resultArr := make([]entries.Feedback, len(found.Arr))

	for index, it := range found.Arr {
		resultArr[index] = entries.Feedback{
			NewFeedback: entries.NewFeedback{
				Message: it.Message,
				Target:  it.ID.Hex(),
			},
			CreationDate: it.CreationDate.Time(),
		}
	}

	sort.Sort(entries.FeedbackSortable(resultArr))

	return resultArr, nil
}
