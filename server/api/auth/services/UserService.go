package services

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/LastSprint/GooodBack/api/auth/entities"
	"github.com/LastSprint/GooodBack/api/auth/errors"
	"github.com/LastSprint/GooodBack/common"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"
	"time"
)

type UserRepo interface {
	GetUserById(userOauthId, authProvider string) (*entities.UserInfo, error)
	CreateUser(userOauthId, authProvider string) (*entities.UserInfo, error)

	SetRefreshTokenForUser(userId, token string) error

	CheckIfUserAllowedToAuth(provider, userId string) (bool, error)
}

type UserService struct {
	accessTokenKey         []byte
	refreshTokenPublicKey  *ecdsa.PublicKey
	refreshTokenPrivateKey *ecdsa.PrivateKey

	Repo UserRepo
}

func InitUserService(repo UserRepo, accessTokenKey, refreshTokenPublicKey, refreshTokenPrivateKey []byte) (*UserService, error) {
	privateKey, err := jwt.ParseECPrivateKeyFromPEM(refreshTokenPrivateKey)

	if err != nil {
		return nil, fmt.Errorf("couldn't parse private key %w", err)
	}

	publicKey, err := jwt.ParseECPublicKeyFromPEM(refreshTokenPublicKey)

	if err != nil {
		return nil, fmt.Errorf("couldn't parse public key %w", err)
	}

	return &UserService{
		accessTokenKey:         accessTokenKey,
		refreshTokenPublicKey:  publicKey,
		refreshTokenPrivateKey: privateKey,
		Repo:                   repo,
	}, nil
}

func (u *UserService) GetTokens(userId string, provider string) (*oauth2.Token, error) {

	allowed, err := u.Repo.CheckIfUserAllowedToAuth(provider, userId)

	if err != nil {
		return nil, fmt.Errorf("couldn't check if user is allowed to login -> %w", err)
	}

	if !allowed {
		return nil, errors.NotAllowedToLogin
	}

	user, err := u.Repo.GetUserById(userId, provider)

	if err != nil {
		return nil, fmt.Errorf("couldn't load user from repo with error %w", err)
	}

	if user == nil {
		user, err = u.Repo.CreateUser(userId, provider)
		if err != nil {
			return nil, fmt.Errorf("couldn't create user with error %w", err)
		}
	}

	if len(user.RefreshToken) == 0 {
		token, err := u.authTokensFromUser(user.ID.Hex())
		if err != nil {
			return nil, fmt.Errorf("couldn't create token pair with error %w", err)
		}

		if err = u.Repo.SetRefreshTokenForUser(user.ID.Hex(), token.RefreshToken); err != nil {
			return nil, fmt.Errorf("couldn't set refresh token for user %s with error %w", user.ID.Hex(), err)
		}

		return token, nil
	}

	accessToken, err := createAccessToken(user.ID.Hex()).SignedString(u.accessTokenKey)

	if err != nil {
		return nil, fmt.Errorf("couldn't create access token for user %s with error %w", user.ID.Hex(), err)
	}

	return &oauth2.Token{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		RefreshToken: user.RefreshToken,
	}, nil
}

func (u *UserService) RefreshAccessToken(refreshTokenStr string) (string, error) {
	claims := common.CustomJWTClaims{}
	_, err := jwt.ParseWithClaims(refreshTokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return u.refreshTokenPublicKey, nil
	})

	if err != nil {
		return "", fmt.Errorf("couldn't parse refresh token %s with error %w", refreshTokenStr, err)
	}

	result, err := createAccessToken(claims.ID).SignedString(u.accessTokenKey)

	if err != nil {
		return "", fmt.Errorf("couldn't sign new access token with error %w", err)
	}

	return result, nil
}

func (u *UserService) authTokensFromUser(userId string) (*oauth2.Token, error) {
	access, err := createAccessToken(userId).SignedString(u.accessTokenKey)

	if err != nil {
		return nil, fmt.Errorf("couldn't sign access token with error %w", err)
	}

	refresh, err := createRefreshToken(userId).SignedString(u.refreshTokenPrivateKey)

	if err != nil {
		return nil, fmt.Errorf("couldn't sign refresh token with error %w", err)
	}

	return &oauth2.Token{
		AccessToken:  access,
		RefreshToken: refresh,
		TokenType:    "Bearer",
	}, nil
}

func createAccessToken(userId string) *jwt.Token {

	now := time.Now()

	claims := &common.CustomJWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "lastsprint.goodback.debug",
			ExpiresAt: &jwt.NumericDate{
				Time: now.Add(time.Hour * 24 * 30),
			},
			IssuedAt: &jwt.NumericDate{
				Time: now,
			},
			ID: userId,
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
}

func createRefreshToken(userId string) *jwt.Token {
	now := time.Now()

	claims := &common.CustomJWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "lastsprint.goodback.debug",
			ExpiresAt: &jwt.NumericDate{
				Time: now.Add(time.Hour * 24 * 256),
			},
			IssuedAt: &jwt.NumericDate{
				Time: now,
			},
			ID: userId,
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodES256, claims)
}
