package slack

import (
	repos2 "github.com/LastSprint/GooodBack/api/feedback/repos"
	"github.com/LastSprint/GooodBack/api/feedback/services"
	"github.com/LastSprint/GooodBack/api/slack/repos"
	services2 "github.com/LastSprint/GooodBack/api/slack/services"
	"go.mongodb.org/mongo-driver/mongo"
)

func AssembleSlackApi(client *mongo.Client, userRepo services.UserRepository, slackToken string) *Api {
	return &Api{
		CmdHandler: &services2.FeedbackCommandHandler{
			Repo: &repos.SlackRepo{
				SlackToken: slackToken,
			},
		},
		InteractionHandler: &services2.SlackInteractivityWebHandler{
			FeedbackSrv: &services.FeedbackService{
				Repo: &repos2.FeedbackRepository{
					Client: client,
				},
				UserRepo: userRepo,
			}},
	}
}
