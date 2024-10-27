package song

import (
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/model"
	"repertoire/utils/wrapper"
)

type GetSongSectionTypes struct {
	repository repository.SongRepository
	jwtService service.JwtService
}

func NewGetSongSectionTypes(repository repository.SongRepository, jwtService service.JwtService) GetSongSectionTypes {
	return GetSongSectionTypes{
		repository: repository,
		jwtService: jwtService,
	}
}

func (g GetSongSectionTypes) Handle(token string) (result []model.SongSectionType, e *wrapper.ErrorCode) {
	userID, errCode := g.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return result, errCode
	}

	err := g.repository.GetSectionTypes(&result, userID)
	if err != nil {
		return result, wrapper.InternalServerError(err)
	}

	return result, nil
}
