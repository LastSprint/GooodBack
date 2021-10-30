package main

import (
	"github.com/LastSprint/GooodBack/api/auth"
	"github.com/LastSprint/GooodBack/oauth/providers"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

type config struct {

	BaseUri string `env:"BASE_URI" envDefault:"http://localhost"`
	BasePath string `env:"BASE_PATH" envDefault:"/api/v1"`

	GoogleClientId string `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET,unset"`
}

func main() {

	cfg := config{}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalln("Can't parse env config with error", err.Error())
		return
	}

	provs := map[string]auth.OAuth2Provider{
		auth.GoogleOAuthProvider: &providers.GoogleOAuth2Provider{
			ClientId:              cfg.GoogleClientId,
			ClientSecret:          cfg.GoogleClientSecret,
			ThisServerRedirectURL: cfg.BaseUri + cfg.BasePath + "/auth/thridparty/redirect",
		},
	}

	authApi := auth.AssembleApi(provs)

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route(cfg.BasePath, func(r chi.Router) {
		authApi.Start(r)
	})

	log.Fatalln(http.ListenAndServe(":80", r))
}
