package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/arseniisemenow/s21auto-client-go/util"

	"github.com/go-resty/resty/v2"
)

type tokenResponse struct {
	Error            string `json:"error"`
	AccessToken      string `json:"access_token"`
	ExpiresIn        int64  `json:"expires_in"`
	RefreshExpiresIn int64  `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	IDToken          string `json:"id_token"`
	NotBeforePolicy  int64  `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

type Token struct {
	AccessToken  string
	RefreshToken string
	Username     string
	Password     string

	IssueTime  int64
	ExpiryTime int64
}

const (
	tokenUrl = "https://auth.21-school.ru/auth/realms/EduPowerKeycloak/protocol/openid-connect/token"
	clientID = "school21"
)

func (token *Token) Refresh(ctx context.Context) error {
	// If token is still valid (with 60 second buffer), no need to refresh
	if token.AccessToken != "" && (time.Now().Unix() < token.ExpiryTime-60) {
		return nil
	}

	client := resty.New()

	var formData map[string]string
	if token.RefreshToken != "" {
		// Use refresh token if available
		formData = map[string]string{
			"client_id":     clientID,
			"grant_type":    "refresh_token",
			"refresh_token": token.RefreshToken,
		}
	} else {
		// Use username/password for initial auth
		formData = map[string]string{
			"client_id":  clientID,
			"grant_type": "password",
			"username":   token.Username,
			"password":   token.Password,
		}
	}

	res, err := client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(formData).
		Post(tokenUrl)

	if err != nil {
		return fmt.Errorf("token request failed: %w", err)
	}

	if res.StatusCode() != 200 {
		return fmt.Errorf("token request failed with status %d: %s", res.StatusCode(), res.String())
	}

	tokenResp, err := util.UnmarshalJson[tokenResponse](res.Body())
	if err != nil {
		return fmt.Errorf("failed to parse token response: %w", err)
	}

	if tokenResp.Error != "" {
		return fmt.Errorf("unable to get access token: %s", tokenResp.Error)
	}

	token.AccessToken = tokenResp.AccessToken
	token.RefreshToken = tokenResp.RefreshToken
	token.IssueTime = time.Now().Unix()
	token.ExpiryTime = token.IssueTime + tokenResp.ExpiresIn

	return nil
}

func RequestToken(username, password string, ctx context.Context) (Token, error) {
	token := Token{
		Username: username,
		Password: password,
	}

	err := token.Refresh(ctx)
	if err != nil {
		return Token{}, err
	}

	return token, nil
}
