package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"repertoire/server/data/http/client"
	"repertoire/server/internal/wrapper"
)

type AuthService interface {
	SignIn(email string, password string) (string, *wrapper.ErrorCode)
}

type authService struct {
	authClient client.AuthClient
}

func NewAuthService(authClient client.AuthClient) AuthService {
	return &authService{authClient: authClient}
}

func (a authService) SignIn(email string, password string) (string, *wrapper.ErrorCode) {
	response, err := a.authClient.SignIn(email, password)
	if err != nil {
		return "", wrapper.InternalServerError(err)
	}
	if response.StatusCode() != http.StatusOK {
		return "", wrapper.InternalServerError(errors.New("failed to sign in" + response.String()))
	}
	var token string
	if err = json.Unmarshal(response.Body(), &token); err != nil {
		return "", wrapper.InternalServerError(err)
	}
	return token, nil
}
