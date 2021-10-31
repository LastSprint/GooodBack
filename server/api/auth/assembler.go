package auth

import (
	"github.com/LastSprint/GooodBack/api/auth/services"
)

func AssembleApi(userRepo services.UserRepo, providers map[string]OAuth2Provider, jwtAccessTokenSeed, jwtRefreshPubKey, jwtRefreshPrKey []byte) (*Api, error) {

	srv, err := services.InitUserService(userRepo, jwtAccessTokenSeed, jwtRefreshPubKey, jwtRefreshPrKey)

	if err != nil {
		return nil, err
	}

	return &Api{
		Providers: providers,
		Srv:       srv,
	}, nil
}
