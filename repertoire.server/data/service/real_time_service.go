package service

import (
	"context"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"repertoire/server/data/cache"
	"repertoire/server/data/realtime"
	"repertoire/server/internal"
	"time"
)

type RealTimeService interface {
	Publish(channel string, userID string, payload any) error
	CreateToken(userID string) string
}

type realTimeService struct {
	env    internal.Env
	cache  cache.CentrifugoCache
	client realtime.CentrifugoClient
}

func NewRealTimeService(
	env internal.Env,
	cache cache.CentrifugoCache,
	client realtime.CentrifugoClient,
) RealTimeService {
	return realTimeService{
		env:    env,
		cache:  cache,
		client: client,
	}
}

func (r realTimeService) Publish(channel string, userID string, payload any) error {
	parsedPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	r.client.SetToken(r.getToken(userID))
	err = r.client.Connect()
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

func (r realTimeService) getToken(userID string) string {
	// get from cache
	tokenKey := "token#" + userID
	token, found := r.cache.Get(tokenKey)
	if found {
		return token.(string)
	}

	// create token and set in cache
	createdToken := r.CreateToken(userID)
	r.cache.Set(tokenKey, createdToken, time.Hour)
	return createdToken
}

func (r realTimeService) CreateToken(userID string) string {
	env := internal.NewEnv()

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti": uuid.New().String(),
		"sub": userID,
		"iss": env.CentrifugoJWTIssuer,
		"aud": env.CentrifugoJWTAudience,
		"iat": time.Now().UTC().Unix(),
		"exp": time.Now().UTC().Add(time.Hour).Unix(),
	})
	token, _ := claims.SignedString([]byte(env.CentrifugoJWTSecretKey))
	return token
}
