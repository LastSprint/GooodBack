package main

import (
	"context"
	"encoding/base64"
	"github.com/LastSprint/GooodBack/api/auth"
	"github.com/LastSprint/GooodBack/api/auth/repos"
	"github.com/LastSprint/GooodBack/api/feedback"
	"github.com/LastSprint/GooodBack/common/middlewares"
	"github.com/LastSprint/GooodBack/oauth/providers"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"net/http"
)

type config struct {
	BaseUri  string `env:"BASE_URI" envDefault:"http://localhost"`
	BasePath string `env:"BASE_PATH" envDefault:"/api/v1"`

	GoogleClientId     string `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET,unset"`

	JwtAccessTokenSeed        string `env:"JWT_ACCESS_TOKEN_SEED,unset"`
	JwtRefreshTokenPublicKey  string `env:"JWT_REFRESH_TOKEN_PUBLIC_KEY_PATH,unset"`
	JwtRefreshTokenPrivateKey string `env:"JWT_REFRESH_TOKEN_PRIVATE_KEY_PATH,unset"`

	AppDataMongoDbConnectionString string `env:"APP_DATA_MONGO_CONNECTION_STRING,unset" envDefault:"mongodb://root:root@localhost:6645"`
}

func main() {

	cfg := config{}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalln("Can't parse env config with error", err.Error())
		return
	}

	pubKey, err := ioutil.ReadFile(cfg.JwtRefreshTokenPublicKey)

	if err != nil {
		log.Fatalf("Couldn't read refresh token public key")
	}

	prKey, err := ioutil.ReadFile(cfg.JwtRefreshTokenPrivateKey)

	if err != nil {
		log.Fatalf("Couldn't read refresh token private key")
	}

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

	userRepo := &repos.UserRepo{Client: client}

	authApi, err := auth.AssembleApi(userRepo, provs, access, pubKey, prKey)

	if err != nil {
		log.Fatalln("Error while building Auth API", err.Error())
		return
	}

	feedbackApi := feedback.AssembleFeedbackApi(client, userRepo)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	mdlw := &middlewares.AccessTokenValidatorMiddleware{Key: access}
	r.Use(mdlw.ExtractToken)

	r.Route(cfg.BasePath, func(r chi.Router) {
		authApi.Start(r)
		feedbackApi.Start(r)
	})

	log.Fatalln(http.ListenAndServe(":80", r))
}
