package service

import (
	"repertoire/auth/data/service"
	"repertoire/auth/internal"
	"repertoire/auth/internal/wrapper"
)

type CentrifugoService interface {
	Token(token string) (string, string, *wrapper.ErrorCode)
}

type centrifugoService struct {
	env        internal.Env
	jwtService service.JwtService
}

func NewCentrifugoService(
	env internal.Env,
	jwtService service.JwtService,
) CentrifugoService {
	return &centrifugoService{
		env:        env,
		jwtService: jwtService,
	}
}

func (c *centrifugoService) Token(token string) (string, string, *wrapper.ErrorCode) {
	userID, errCode := c.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return "", "", errCode
	}
	centrifugoToken, errCode := c.jwtService.CreateCentrifugoToken(userID.String())
	if errCode != nil {
		return "", "", errCode
	}
	return centrifugoToken, c.env.CentrifugoJwtExpirationTime, nil
}
