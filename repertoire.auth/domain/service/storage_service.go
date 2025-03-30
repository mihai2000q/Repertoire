package service

import (
	"repertoire/auth/data/service"
	"repertoire/auth/internal"
	"repertoire/auth/internal/wrapper"
)

type StorageService interface {
	Token(grantType string, clientID string, clientSecret string, token string) (string, string, *wrapper.ErrorCode)
}

type storageService struct {
	env        internal.Env
	jwtService service.JwtService
}

func NewStorageService(
	env internal.Env,
	jwtService service.JwtService,
) StorageService {
	return &storageService{
		env:        env,
		jwtService: jwtService,
	}
}

func (c *storageService) Token(grantType string, clientID string, clientSecret string, token string) (string, string, *wrapper.ErrorCode) {
	errCode := c.jwtService.ValidateCredentials(grantType, clientID, clientSecret)
	if errCode != nil {
		return "", "", errCode
	}
	userID, errCode := c.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return "", "", errCode
	}
	storageToken, errCode := c.jwtService.CreateStorageToken(userID.String())
	if errCode != nil {
		return "", "", errCode
	}
	return storageToken, c.env.StorageJwtExpirationTime, nil
}
