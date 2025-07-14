package services

import (
	"encoding/json"
	"net/http"
	"os"
)

type BasicUser struct {
	ID           string `json:"_id"`
	IsVendor     bool   `json:"isVendor"`
	IsRegistered bool   `json:"isRegistered"`
}

func CheckUserToken(token string) (BasicUser, error) {
	authURL := os.Getenv("CLIENT_SERVER")
	userReq, err := http.NewRequest("GET", authURL+"/api/v2/auth/check", nil)
	if err != nil {
		return BasicUser{}, err
	}
	userReq.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(userReq)
	if err != nil || resp.StatusCode != http.StatusOK {
		return BasicUser{}, err
	}
	defer resp.Body.Close()

	var user BasicUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return BasicUser{}, err
	}

	return user, nil
}
