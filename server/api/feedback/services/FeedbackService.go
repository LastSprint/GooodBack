package services

import (
	"github.com/LastSprint/GooodBack/api/feedback/entries"
	"github.com/LastSprint/GooodBack/common"
	"log"
)

type UserRepository interface {
	GetUserIdByTarget(target string) (string, error)
}

type Repository interface {
	CreateNew(message, userID string) error
	Read(userID string) ([]entries.Feedback, error)
}

type FeedbackService struct {
	Repo     Repository
	UserRepo UserRepository
}

func (f *FeedbackService) Write(feedback entries.NewFeedback) error {

	id, err := f.UserRepo.GetUserIdByTarget(feedback.Target)

	if err != nil {
		log.Println("[ERR] couldn't check user existing due to", err.Error())
		return common.NotFound
	}

	if err = f.Repo.CreateNew(feedback.Message, id); err != nil {
		log.Printf("[ERR] couldn't create feedback with error %s", err.Error())
	}

	return nil
}

func (f *FeedbackService) Read(userId string) ([]entries.Feedback, error) {

	res, err := f.Repo.Read(userId)

	if err != nil {
		log.Printf("[ERR] couldn't read feedback with error %s", err.Error())
	}

	return res, nil
}
