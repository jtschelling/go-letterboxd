package letterboxd

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// AccessToken holds information to sign every request with a valid access token
// data provided by the getAccessToken() function
type AccessToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	NotBefore    int    `json:"notBefore"`
	Issuer       string `json:"issuer"`
	EncodedToken string `json:"encodedToken"`
}

type AuthorizationCodeBody struct {
	GrantType    string `json:"grant_type"`
	Code         string `json:"code"`
	RedirectURI  string `json:"redirect_uri"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type RefreshTokenBody struct {
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

// OAuthError is included with 400 responses
type OAuthError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"errorDescription"`
}

// Use the authorization_code grant type to obtain a token from /auth/token
func (c *Client) PostAuthTokenAuthorizationCode(code string, redirectURI string) (AccessToken, error) {
	authorizationCodeBody := AuthorizationCodeBody{
		GrantType:    "authorization_code",
		Code:         code,
		RedirectURI:  redirectURI,
		ClientID:     c.clientID,
		ClientSecret: c.clientSecret,
	}

	// Marshal the authorizationCodeBody into JSON
	body, err := json.Marshal(authorizationCodeBody)
	if err != nil {
		return AccessToken{}, err
	}

	// Create an HTTP POST request to c.BaseURL/auth/token
	req, err := http.NewRequest("POST", c.BaseURL+"/auth/token", bytes.NewReader(body))
	if err != nil {
		return AccessToken{}, err
	}

	// Set the Content-Type header to application/json
	req.Header.Set("Content-Type", "application/json")

	// Perform the HTTP request
	client := &http.Client{
		Timeout: time.Second * 10, // You can adjust the timeout as needed
	}

	resp, err := client.Do(req)
	if err != nil {
		return AccessToken{}, err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return AccessToken{}, errors.New("failed to obtain access token")
	}

	// Parse the response body into an AccessToken struct
	var accessToken AccessToken
	err = json.NewDecoder(resp.Body).Decode(&accessToken)
	if err != nil {
		return AccessToken{}, err
	}

	return accessToken, nil
}
