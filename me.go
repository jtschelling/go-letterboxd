package letterboxd

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

type MemberAccount struct {
	Member struct {
		ID          string `json:"id"`
		Username    string `json:"username"`
		DisplayName string `json:"displayName"`
		ShortName   string `json:"shortName"`
		Pronoun     struct {
			ID                  string `json:"id"`
			Label               string `json:"label"`
			SubjectPronoun      string `json:"subjectPronoun"`
			ObjectPronoun       string `json:"objectPronoun"`
			PossessiveAdjective string `json:"possessiveAdjective"`
			PossessivePronoun   string `json:"possessivePronoun"`
			Reflexive           string `json:"reflexive"`
		} `json:"pronoun"`
		Avatar struct {
			Sizes []struct {
				Width  int    `json:"width"`
				Height int    `json:"height"`
				URL    string `json:"url"`
			} `json:"sizes"`
		} `json:"avatar"`
		MemberStatus     string `json:"memberStatus"`
		HideAdsInContent bool   `json:"hideAdsInContent"`
		AccountStatus    string `json:"accountStatus"`
		HideAds          bool   `json:"hideAds"`
		BioLbml          string `json:"bioLbml"`
		FavoriteFilms    []any  `json:"favoriteFilms"`
		Links            []struct {
			Type string `json:"type"`
			ID   string `json:"id"`
			URL  string `json:"url"`
		} `json:"links"`
		PrivateWatchlist bool   `json:"privateWatchlist"`
		Bio              string `json:"bio"`
	} `json:"member"`
	HideAds              bool `json:"hideAds"`
	ShowCustomPostersAds bool `json:"showCustomPostersAds"`
	CanHaveCustomPosters bool `json:"canHaveCustomPosters"`
}

func (c *Client) GetMe(authZToken string) (MemberAccount, error) {
	// Create an HTTP GET request to c.BaseURL+"/me"
	req, err := http.NewRequest(http.MethodGet, c.BaseURL+"/me", nil)
	if err != nil {
		return MemberAccount{}, err
	}

	var bearer = "Bearer " + authZToken
	req.Header.Set("Authorization", bearer)

	// Perform the HTTP request
	client := &http.Client{
		Timeout: time.Second * 10, // You can adjust the timeout as needed
	}

	resp, err := client.Do(req)
	if err != nil {
		return MemberAccount{}, err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return MemberAccount{}, errors.New(string(b))
	}

	// Parse the response body into a MemberAccount struct
	var memberAccount MemberAccount
	err = json.NewDecoder(resp.Body).Decode(&memberAccount)
	if err != nil {
		return MemberAccount{}, err
	}

	return memberAccount, nil
}
