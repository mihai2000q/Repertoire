package section

import (
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
)

type GetSongSectionTypes struct {
	songSectionRepository repository.SongSectionRepository
	jwtService            service.JwtService
}

func NewGetSongSectionTypes(
	songSectionRepository repository.SongSectionRepository,
	jwtService service.JwtService,
) GetSongSectionTypes {
	return GetSongSectionTypes{
		songSectionRepository: songSectionRepository,
		jwtService:            jwtService,
	}
}

func (g GetSongSectionTypes) Handle(token string) (result []model.SongSectionType, e *httperror.ErrorCode) {
	userID, errCode := g.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return result, errCode
	}

	if err := g.songSectionRepository.GetTypes(&result, userID); err != nil {
		return result, httperror.DatabaseError(err)
	}

	return result, nil
}
