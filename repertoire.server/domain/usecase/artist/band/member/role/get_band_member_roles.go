package role

import (
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type GetBandMemberRoles struct {
	repository repository.ArtistRepository
	jwtService service.JwtService
}

func NewGetBandMemberRoles(repository repository.ArtistRepository, jwtService service.JwtService) GetBandMemberRoles {
	return GetBandMemberRoles{
		repository: repository,
		jwtService: jwtService,
	}
}

func (g GetBandMemberRoles) Handle(token string) (result []model.BandMemberRole, e *wrapper.ErrorCode) {
	userID, errCode := g.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return result, errCode
	}

	err := g.repository.GetBandMemberRoles(&result, userID)
	if err != nil {
		return result, wrapper.InternalServerError(err)
	}

	return result, nil
}
