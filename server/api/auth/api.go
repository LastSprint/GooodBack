package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/LastSprint/GooodBack/api/auth/entries"
	auth_errors "github.com/LastSprint/GooodBack/api/auth/errors"
	"github.com/LastSprint/GooodBack/common"
	"github.com/go-chi/chi/v5"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"net/url"
)

/*
Contains OAuth 2.0 for Google implementation
*/

type OAuth2Provider interface {
	// GetRedirectUrl have to return URL for first redirection from this service to third-party service
	GetRedirectUrl() (*url.URL, error)
	ExchangeAuthCode(code, redirectUrl string) (*oauth2.Token, error)
	// GetUserID have to return user ID in terms of third party service.
	// for example email or username (we don't care)
	GetUserID(token *oauth2.Token) (string, error)
}

type Service interface {
	// GetTokens check if user already exist and if it is will return a token pair
	// if not will create new user and return their tokens
	GetTokens(userId string, provider string) (*oauth2.Token, error)
	RefreshAccessToken(refreshTokenStr string) (string, error)
}

const (
	GoogleOAuthProvider string = "google"
)

type Api struct {
	Providers map[string]OAuth2Provider
	Srv       Service
}

func (api *Api) Start(router chi.Router) {
	router.Route("/auth", func(r chi.Router) {
		r.Get("/thridparty/{provider}", api.handleThirdPartyAuthIntent)
		r.Get("/thridparty/redirect", api.handleThirdPartyServiceCodeRedirect)
		r.Post("/thirparty/code", api.handleThirdPartyCodeExchange)
		r.Post("/thirparty/code", api.handleThirdPartyCodeExchange)

		r.Post("/refresh", api.handleAccessTokenRefresh)
	})
}

// handleThirdPartyAuthIntent will be called when user want to authorize with some third party service
func (api *Api) handleThirdPartyAuthIntent(w http.ResponseWriter, r *http.Request) {
	providerName := chi.URLParam(r, "provider")

	provider, ok := api.Providers[providerName]

	if !ok {
		http.Error(w, fmt.Sprintf("provider with name \"%s\" not found", provider), http.StatusBadRequest)
		return
	}

	redirectUrl, err := provider.GetRedirectUrl()

	if err != nil {
		log.Printf("[ERR] provider \"%s\" couldn't build refirect url with error \"%s\"", providerName, err.Error())
		http.Error(w, "Couldn't build url for redirection", http.StatusInternalServerError)
		return
	}

	query := redirectUrl.Query()

	// FIXME: Add random seed to state

	query.Set("state", providerName)

	redirectUrl.RawQuery = query.Encode()

	http.Redirect(w, r, redirectUrl.String(), http.StatusPermanentRedirect)
}

// handleThirdPartyServiceCodeRedirect catches redirect from third party service
// and tries to parse `code` from URL query
// then exchanges that code on access tokens
func (api *Api) handleThirdPartyServiceCodeRedirect(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	if len(code) == 0 {
		log.Printf("[ERR] third party sent request without code: %s", r.URL)
		http.Error(w, "service didnt return code", http.StatusInternalServerError)
		return
	}

	// FIXME: Add seed validation

	providerName := r.URL.Query().Get("state")

	tokens, err := api.getUserTokens(code, providerName, "")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// FIXME: Write token pair in cookies

	if err = json.NewEncoder(w).Encode(tokens); err != nil {
		log.Printf("Can't serialize aouth tokens to JSON with error \"%s\"", err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}

// handleThirdPartyCodeExchange just like handleThirdPartyServiceCodeRedirect BUT this method considered to be called by client not by third-party service
func (api *Api) handleThirdPartyCodeExchange(w http.ResponseWriter, r *http.Request) {
	entry := entries.OAuth2CodeExchange{}

	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, "couldn't parse response", http.StatusBadRequest)
		return
	}

	tokens, err := api.getUserTokens(entry.Code, entry.Provider, entry.RedirectUrl)

	if errors.Is(err, auth_errors.NotAllowedToLogin) {
		http.Error(w, "You are not allowed to login 🧙‍️", http.StatusForbidden)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = common.SetCookie(w, r, "Authorization", tokens.AccessToken, 60 * 60 * 24 * 30); err != nil {
		log.Println("[ERR] Couldn't set cookie ->", err.Error())
		http.Error(w, "something went wrong", http.StatusBadRequest)
		return
	}

	if err = common.SetCookie(w, r, "Refreshing", tokens.RefreshToken, 60 * 60 * 24 * 30 * 12); err != nil {
		log.Println("[ERR] Couldn't set cookie ->", err.Error())
		http.Error(w, "something went wrong", http.StatusBadRequest)
		return
	}

	fmt.Println(w.Header())
}

func (api *Api) handleAccessTokenRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken := struct {
		RefreshToken string `json:"refresh_token"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&refreshToken); err != nil {
		http.Error(w, "can't parse request", http.StatusBadRequest)
		return
	}

	newToken, err := api.Srv.RefreshAccessToken(refreshToken.RefreshToken)

	if err != nil {
		log.Printf("[ERR] while refreshing token %s", err.Error())
		http.Error(w, "can't parse request", http.StatusBadRequest)
		return
	}

	result := oauth2.Token{
		AccessToken:  newToken,
		RefreshToken: refreshToken.RefreshToken,
	}

	if err = json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("[ERR] couldn't encode token to json with error %s", err.Error())
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
}

func (api *Api) getUserTokens(oauth2Code, oauth2ProviderName string, redirectUrl string) (*oauth2.Token, error) {
	provider, ok := api.Providers[oauth2ProviderName]

	if !ok {
		return nil, fmt.Errorf("provider \"%s\" not found by name \"%s\"", provider, oauth2ProviderName)
	}

	token, err := provider.ExchangeAuthCode(oauth2Code, redirectUrl)

	if err != nil {
		log.Printf("[ERR] provider \"%s\" couldn't exchange code \"%s\" with err \"%s\"", oauth2ProviderName, oauth2Code, err.Error())
		return nil, fmt.Errorf("couldn't exchange code")
	}

	thirdPartyUserId, err := provider.GetUserID(token)

	if err != nil {
		log.Printf("[ERR] couldn't get user ID from provider \"%s\" with error \"%s\"", oauth2ProviderName, err.Error())
		return nil, fmt.Errorf("can't get user ID from third party")
	}

	tokens, err := api.Srv.GetTokens(thirdPartyUserId, oauth2ProviderName)

	if err != nil {
		log.Printf("[ERR] couldn't get tokens for userId \"%s\" with error \"%s\"", thirdPartyUserId, err.Error())
		return nil, fmt.Errorf("can't get tokens for user %w", err)
	}

	return tokens, err
}
