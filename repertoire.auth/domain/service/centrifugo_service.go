package service

import (
	"github.com/google/uuid"
	"repertoire/auth/data/service"
	"repertoire/auth/internal/wrapper"
)

type CentrifugoService interface {
	Token(token string) (string, string, *wrapper.ErrorCode)
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
	if errCode != nil {
		return "", "", errCode
	}
	return centrifugoToken, c.env.CentrifugoJwtExpirationTime, nil
}
