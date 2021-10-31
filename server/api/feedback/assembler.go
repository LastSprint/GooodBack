package feedback

import (
	"github.com/LastSprint/GooodBack/api/feedback/repos"
	"github.com/LastSprint/GooodBack/api/feedback/services"
	"go.mongodb.org/mongo-driver/mongo"
)

func AssembleFeedbackApi(client *mongo.Client, userRepo services.UserRepository) *Api {
	return &Api{
		Srv: &services.FeedbackService{
			Repo: &repos.FeedbackRepository{
				Client: client,
			},
			UserRepo: userRepo,
		},
	}
}
