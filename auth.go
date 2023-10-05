package letterboxd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
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
	formData := url.Values{}
	formData.Set("grant_type", "authorization_code")
	formData.Set("code", code)
	formData.Set("redirect_uri", redirectURI)
	formData.Set("client_id", c.clientID)
	formData.Set("client_secret", c.clientSecret)

	fmt.Println(formData.Encode())

	// Create an HTTP POST request to c.BaseURL+"/auth/token"
	req, err := http.NewRequest(http.MethodPost, c.BaseURL+"/auth/token", strings.NewReader(formData.Encode()))
	if err != nil {
		return AccessToken{}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
		b, _ := io.ReadAll(resp.Body)
		return AccessToken{}, errors.New(string(b))
	}

	// Parse the response body into an AccessToken struct
	var accessToken AccessToken
	err = json.NewDecoder(resp.Body).Decode(&accessToken)
	if err != nil {
		return AccessToken{}, err
	}

	return accessToken, nil
}
