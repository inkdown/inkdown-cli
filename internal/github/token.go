package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	Error       string `json:"error"`
}

func PollForToken(clientID, deviceCode string, interval int) (string, error) {
	for {
		time.Sleep(time.Duration(interval) * time.Second)

		data := url.Values{}
		data.Set("client_id", clientID)
		data.Set("device_code", deviceCode)
		data.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")

		req, _ := http.NewRequest(
			"POST",
			"https://github.com/login/oauth/access_token",
			strings.NewReader(data.Encode()),
		)

		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", err
		}

		var res TokenResponse
		json.NewDecoder(resp.Body).Decode(&res)
		resp.Body.Close()

		if res.AccessToken != "" {
			return res.AccessToken, nil
		}

		if res.Error != "" && res.Error != "authorization_pending" {
			return "", fmt.Errorf("oauth error: %s", res.Error)
		}
	}
}

func ValidateToken(token string) error {
	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "community-cli")

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return errors.New("Invalid token")
	}
	return nil
}

func SaveToken(token string) error {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".community-cli")

	os.MkdirAll(dir, 0700)

	path := filepath.Join(dir, "token")
	return os.WriteFile(path, []byte(token), 0600)
}

func LoadToken() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(home, ".community-cli", "token")
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}
