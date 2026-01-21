package service

import (
	"repertoire/auth/data/service"
	"repertoire/auth/internal/wrapper"
	"repertoire/auth/model"

	"github.com/google/uuid"
)

type StorageService interface {
	Token(clientCredentials model.ClientCredentials, userID uuid.UUID) (string, string, *wrapper.ErrorCode)
}

type storageService struct {
	jwtService service.JwtService
}

func NewStorageService(jwtService service.JwtService) StorageService {
	return &storageService{jwtService: jwtService}
}

func (c *storageService) Token(
	clientCredentials model.ClientCredentials,
	userID uuid.UUID,
) (string, string, *wrapper.ErrorCode) {
	errCode := c.jwtService.ValidateCredentials(clientCredentials)
	if errCode != nil {
		return "", "", errCode
	}
	return c.jwtService.CreateStorageToken(userID)
}
