package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/LastSprint/GooodBack/api/auth"
	"github.com/LastSprint/GooodBack/api/auth/repos"
	"github.com/LastSprint/GooodBack/api/feedback"
	"github.com/LastSprint/GooodBack/api/slack"
	"github.com/LastSprint/GooodBack/common/middlewares"
	"github.com/LastSprint/GooodBack/oauth/providers"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"strings"
)

type config struct {
	BaseUri  string `env:"BASE_URI" envDefault:"http://localhost"`
	BasePath string `env:"BASE_PATH" envDefault:"/api/v1"`

	GoogleClientId     string `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET,unset"`

	JwtAccessTokenSeed        string `env:"JWT_ACCESS_TOKEN_SEED,unset"`
	JwtRefreshTokenPublicKey  string `env:"JWT_REFRESH_TOKEN_PUBLIC_KEY,unset"`
	JwtRefreshTokenPrivateKey string `env:"JWT_REFRESH_TOKEN_PRIVATE_KEY,unset"`

	AppDataMongoDbConnectionString string `env:"APP_DATA_MONGO_CONNECTION_STRING,unset" envDefault:"mongodb://root:root@localhost:6645"`

	SlackToken string `env:"SLACK_TOKEN,unset"`
}

func main() {

	cfg := config{}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalln("Can't parse env config with error", err.Error())
		return
	}

	cfg.GoogleClientId = strings.TrimSpace(cfg.GoogleClientId)
	cfg.GoogleClientSecret = strings.TrimSpace(cfg.GoogleClientSecret)
	cfg.SlackToken = strings.TrimSpace(cfg.SlackToken)
	cfg.AppDataMongoDbConnectionString = strings.TrimSpace(cfg.AppDataMongoDbConnectionString)

	fmt.Println(cfg)

	pubKey := []byte(strings.ReplaceAll(cfg.JwtRefreshTokenPublicKey, "\\n", "\n"))
	prKey := []byte(strings.ReplaceAll(cfg.JwtRefreshTokenPrivateKey, "\\n", "\n"))

	provs := map[string]auth.OAuth2Provider{
		auth.GoogleOAuthProvider: &providers.GoogleOAuth2Provider{
			ClientId:              cfg.GoogleClientId,
			ClientSecret:          cfg.GoogleClientSecret,
			ThisServerRedirectURL: cfg.BaseUri + cfg.BasePath + "/auth/thridparty/redirect",
		},
	}

	access, err := base64.StdEncoding.DecodeString(cfg.JwtAccessTokenSeed)

	if err != nil {
		log.Fatalf("Cant decode access token seed from string")
		return
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.AppDataMongoDbConnectionString))

	if err != nil {
		log.Fatalln("[ERR] Can't connect to mongodb ->", err.Error())
		return
	}

	userRepo := &repos.UserRepo{Client: client}

	authApi, err := auth.AssembleApi(userRepo, provs, access, pubKey, prKey)

	if err != nil {
		log.Fatalln("Error while building Auth API", err.Error())
		return
	}

	feedbackApi := feedback.AssembleFeedbackApi(client, userRepo)

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Set-Cookie"},
		AllowCredentials: true,
	}))
	r.Use(middleware.Logger)
	mdlw := &middlewares.AccessTokenValidatorMiddleware{Key: access}
	r.Use(mdlw.ExtractToken)

	slackApi := slack.AssembleSlackApi(client, userRepo, cfg.SlackToken)

	r.Route(cfg.BasePath, func(r chi.Router) {
		authApi.Start(r)
		feedbackApi.Start(r)
		slackApi.Start(r)
	})

	log.Fatalln(http.ListenAndServe(":80", r))
}
