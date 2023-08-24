package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	cfg = Get()
)

func validateSecretKey(secretKey string) (record SecretKey, err error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/auth/validate-secret-key", cfg.AUTH_URL), nil)
	if err != nil {
		return
	}

	req.Header.Set("Secret-Key", secretKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	bodyString := string(body)

	httpResp := ValidateSecretKeyResponse{}
	err = json.Unmarshal([]byte(bodyString), &httpResp)
	if err != nil {
		return
	}

	if httpResp.Code != 200 {
		return record, fmt.Errorf(httpResp.Status)
	}

	return httpResp.Data, nil
}

func ValidateCurrentUser(c *gin.Context) (user GetCurrentUserResponse, err error) {

	secretKey := c.GetHeader("Secret-Key")

	if secretKey != "" {
		record, err := validateSecretKey(secretKey)
		if err != nil {
			return user, err
		}

		if record.ID == "" {
			return user, fmt.Errorf("Secret key not found")
		}

		simulatedUser := User{
			ID:       record.ID,
			Username: record.Serial,
			Email:    record.Name,
			Role:     record.Role,
		}

		user = GetCurrentUserResponse{
			Status: "OK",
			Code:   200,
			Data: CurrentUserData{
				User: simulatedUser,
			},
		}

		return user, nil
	}

	token := c.GetHeader("Authorization")

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/auth/current-user", cfg.AUTH_URL), nil)
	if err != nil {
		return
	}

	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	bodyString := string(body)

	err = json.Unmarshal([]byte(bodyString), &user)
	if err != nil {
		return
	}

	if user.Code != 200 {
		return user, fmt.Errorf(user.Status)
	}

	return user, nil
}
