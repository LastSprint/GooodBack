package auth

import (
	"github.com/LastSprint/GooodBack/api/auth/repos"
	"github.com/LastSprint/GooodBack/api/auth/services"
)

func AssembleApi(providers map[string]OAuth2Provider, jwtAccessTokenSeed, jwtRefreshPubKey, jwtRefreshPrKey []byte) (*Api, error) {

	srv, err := services.InitUserService(&repos.UserRepo{}, jwtAccessTokenSeed, jwtRefreshPubKey, jwtRefreshPrKey)

	if err != nil {
		return nil, err
	}

	return &Api{
		Providers: providers,
		Srv:       srv,
	}, nil
}
