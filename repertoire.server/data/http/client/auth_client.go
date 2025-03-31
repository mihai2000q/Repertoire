package client

import (
	"github.com/go-resty/resty/v2"
	"repertoire/server/data/http"
	"repertoire/server/internal"
)

type AuthClient struct {
	env internal.Env
	*resty.Client
}

func NewAuthClient(client http.RestyClient, env internal.Env) AuthClient {
	return AuthClient{
		env:    env,
		Client: client.SetBaseURL(env.AuthUrl),
	}
}

func (client AuthClient) StorageToken(userID string, result any) (*resty.Response, error) {
	return client.R().
		SetFormData(map[string]string{
			"grant_type":    "client_credentials",
			"client_id":     client.env.AuthClientID,
			"client_secret": client.env.AuthClientSecret,
			"user_id":       userID,
		}).
		SetResult(result).
		Post("/storage/token")
}

func (client AuthClient) CentrifugoToken(userID string, result any) (*resty.Response, error) {
	return client.R().
		SetFormData(map[string]string{
			"grant_type":    "client_credentials",
			"client_id":     client.env.AuthClientID,
			"client_secret": client.env.AuthClientSecret,
			"user_id":       userID,
		}).
		SetResult(result).
		Post("/centrifugo/token")
}

func (client AuthClient) SignIn(email string, password string, result any) (*resty.Response, error) {
	return client.R().
		SetBody(struct{ Email, Password string }{email, password}).
		SetResult(result).
		Put("/sign-in")
}
