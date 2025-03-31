package service

import (
	"github.com/google/uuid"
	"repertoire/auth/data/service"
	"repertoire/auth/internal/wrapper"
	"repertoire/auth/model"
)

type CentrifugoService interface {
	Token(token string) (string, string, *wrapper.ErrorCode)
	PublicToken(clientCredentials model.ClientCredentials, userID uuid.UUID) (string, string, *wrapper.ErrorCode)
}

type centrifugoService struct {
	jwtService service.JwtService
}

func NewCentrifugoService(jwtService service.JwtService) CentrifugoService {
	return &centrifugoService{jwtService: jwtService}
}

func (c *centrifugoService) Token(token string) (string, string, *wrapper.ErrorCode) {
	userID, errCode := c.jwtService.GetUserIDFromJwt(token)
	if errCode != nil {
		return "", "", errCode
	}
	return c.jwtService.CreateCentrifugoToken(userID)
}

func (c *centrifugoService) PublicToken(
	clientCredentials model.ClientCredentials,
	userID uuid.UUID,
) (string, string, *wrapper.ErrorCode) {
	errCode := c.jwtService.ValidateCredentials(clientCredentials)
	if errCode != nil {
		return "", "", errCode
	}
	return c.jwtService.CreateCentrifugoToken(userID)
}
