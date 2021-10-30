package auth

import "github.com/LastSprint/GooodBack/api/auth/services"

func AssembleApi(providers map[string]OAuth2Provider) *Api {
	return &Api{
		Providers: providers,
		Srv:       &services.UserService{

		},
	}
}
