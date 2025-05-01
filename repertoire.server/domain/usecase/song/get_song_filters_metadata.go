package song

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type GetSongFiltersMetadata struct {
	jwtService service.JwtService
	repository repository.SongRepository
}

func NewGetSongFiltersMetadata(
	jwtService service.JwtService,
	repository repository.SongRepository,
) GetSongFiltersMetadata {
	return GetSongFiltersMetadata{
		jwtService: jwtService,
		repository: repository,
	}
}

func (g GetSongFiltersMetadata) Handle(
	request requests.GetSongFiltersMetadataRequest,
	token string,
) (metadata model.SongFiltersMetadata, e *wrapper.ErrorCode) {
	userID, errCode := g.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return metadata, errCode
	}

	err := g.repository.GetFiltersMetadata(&metadata, userID, request.SearchBy)
	if err != nil {
		return metadata, wrapper.InternalServerError(err)
	}
	return metadata, nil
}
