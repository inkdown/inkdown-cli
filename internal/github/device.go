package github

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

func RequestDeviceCode(clientID string) (*DeviceCodeResponse, error) {
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("scope", "public_repo")

	req, _ := http.NewRequest(
		"POST",
		"https://github.com/login/device/code",
		strings.NewReader(data.Encode()),
	)

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res DeviceCodeResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	return &res, err
}
