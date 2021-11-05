package providers

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"net/url"
)

// GoogleOAuth2Provider implement Google OpenID
//
// For details see https://developers.google.com/identity/protocols/oauth2/web-server#httprest_1
type GoogleOAuth2Provider struct {
	ClientId string
	ClientSecret string
	ThisServerRedirectURL string
}

func (g *GoogleOAuth2Provider) GetRedirectUrl() (*url.URL, error) {

	result, err := url.Parse("https://accounts.google.com/o/oauth2/v2/auth")

	if err != nil {
		return nil, err
	}

	query := &url.Values{}

	query.Add("client_id", g.ClientId)
	query.Add("response_type", "code")
	query.Add("scope", "https://www.googleapis.com/auth/userinfo.email")
	query.Add("redirect_uri", g.ThisServerRedirectURL)

	result.RawQuery = query.Encode()

	return result, nil
}

func (g *GoogleOAuth2Provider) ExchangeAuthCode(code, redirectUrl string) (*oauth2.Token, error) {
	requestBody := url.Values{}

	if len(redirectUrl) == 0 {
		redirectUrl = g.ThisServerRedirectURL
	}

	requestBody.Add("code", code)
	requestBody.Add("client_id", g.ClientId)
	requestBody.Add("client_secret", g.ClientSecret)
	requestBody.Add("redirect_uri", redirectUrl)
	requestBody.Add("grant_type", "authorization_code")

	response, err := http.PostForm("https://oauth2.googleapis.com/token", requestBody)

	if err != nil {
		return nil, fmt.Errorf("token request filed with error %w", err)
	}

	data, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, fmt.Errorf("couldn't read response body with error %w", err)
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code is %v; response is: %s", response.StatusCode, string(data))
	}

	tokens := &oauth2.Token{}

	if err = json.Unmarshal(data, tokens); err != nil {
		return nil, fmt.Errorf("couldn't parse token pair from response with error %w", err)
	}

	return tokens, nil
}

func (g *GoogleOAuth2Provider) GetUserID(token *oauth2.Token) (string, error) {
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)

	if err != nil {
		return "", fmt.Errorf("can't get user info from api with error %w", err)
	}

	data, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", fmt.Errorf("couldn't read response body with error %w", err)
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("response status code is %v; response is: %s", response.StatusCode, string(data))
	}

	result := struct {
		Email string `json:"email"`
	}{}

	if err = json.Unmarshal(data, &result); err != nil {
		return "", fmt.Errorf("couldn't parse response %s with error %w", string(data), err.Error())
	}

	return result.Email, nil
}

