package bandmember

import (
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
)

type GetBandMemberRoles struct {
	artistRepository repository.ArtistRepository
	jwtService       service.JwtService
}

func NewGetBandMemberRoles(
	artistRepository repository.ArtistRepository,
	jwtService service.JwtService,
) GetBandMemberRoles {
	return GetBandMemberRoles{
		artistRepository: artistRepository,
		jwtService:       jwtService,
	}
}

func (g GetBandMemberRoles) Handle(token string) (result []model.BandMemberRole, e *httperror.ErrorCode) {
	userID, errCode := g.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return result, errCode
	}

	if err := g.artistRepository.GetBandMemberRoles(&result, userID); err != nil {
		return result, httperror.DatabaseError(err)
	}

	return result, nil
}
