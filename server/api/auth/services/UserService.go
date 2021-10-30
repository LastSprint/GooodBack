package services

import (
	"fmt"
	"golang.org/x/oauth2"
)

type UserService struct {
	
}

func (u *UserService) GetTokens(userId string, provider string) (*oauth2.Token, error) {
	return nil, fmt.Errorf("not implemented yet")
}

