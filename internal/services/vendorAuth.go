package services

import (
	"encoding/json"
	"net/http"
	"os"
	"shop-dashboard/internal/models"
)

func CheckVendorToken(token string) (models.UserInfo, error) {
	url := os.Getenv("CLIENT_SERVER")
	req, err := http.NewRequest("GET", url+"/api/v2/auth/vendor", nil)
	if err != nil {
		return models.UserInfo{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return models.UserInfo{}, err
	}
	defer resp.Body.Close()

	var user models.UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return models.UserInfo{}, err
	}

	return user, nil
}
