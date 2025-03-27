package auth

import (
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
)

type GetCentrifugoToken struct {
	jwtService      service.JwtService
	realTimeService service.RealTimeService
}

func NewGetCentrifugoToken(jwtService service.JwtService, realTimeService service.RealTimeService) GetCentrifugoToken {
	return GetCentrifugoToken{
		jwtService:      jwtService,
		realTimeService: realTimeService,
	}
}

func (r *GetCentrifugoToken) Handle(token string) (string, *wrapper.ErrorCode) {
	userID, errCode := r.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return "", errCode
	}
	centrifugoToken := r.realTimeService.CreateToken(userID.String())
	return centrifugoToken, nil
}
