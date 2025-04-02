package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"repertoire/server/data/cache"
	"repertoire/server/data/http/auth"
	"repertoire/server/data/http/client"
	"repertoire/server/data/realtime"
	"time"
)

type RealTimeService interface {
	Publish(channel string, userID string, payload any) error
}

type realTimeService struct {
	cache      cache.CentrifugoCache
	client     realtime.CentrifugoClient
	authClient client.AuthClient
}

func NewRealTimeService(
	cache cache.CentrifugoCache,
	client realtime.CentrifugoClient,
	authClient client.AuthClient,
) RealTimeService {
	return realTimeService{
		cache:      cache,
		client:     client,
		authClient: authClient,
	}
}

func (r realTimeService) Publish(channel string, userID string, payload any) error {
	accessToken, err := r.getToken(userID)
	if err != nil {
		return err
	}
	r.client.SetToken(accessToken)

	err = r.client.Connect()
	if err != nil {
		return err
	}

	parsedPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = r.client.Publish(context.Background(), channel+":"+userID, parsedPayload)
	if err != nil {
		return err
	}

	err = r.client.Disconnect()
	return err
}

func (r realTimeService) getToken(userID string) (string, error) {
	// get from cache
	tokenKey := "token#" + userID
	token, found := r.cache.Get(tokenKey)
	if found {
		return token.(string), nil
	}

	// fetch token and set in cache
	tokenResult, err := r.fetchToken(userID)
	if err != nil {
		return "", err
	}
	expiresIn, _ := time.ParseDuration(tokenResult.ExpiresIn)
	r.cache.Set(tokenKey, tokenResult.Token, expiresIn)
	return tokenResult.Token, nil
}

func (r realTimeService) fetchToken(userID string) (auth.TokenResponse, error) {
	var result auth.TokenResponse
	response, err := r.authClient.CentrifugoToken(userID, &result)

	if err != nil {
		return auth.TokenResponse{}, err
	}
	if response.StatusCode() != http.StatusOK {
		return auth.TokenResponse{}, errors.New("failed to fetch token: " + response.String())
	}

	return result, nil
}
