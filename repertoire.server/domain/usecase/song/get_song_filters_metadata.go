package song

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
)

type GetSongFiltersMetadata struct {
	jwtService     service.JwtService
	songRepository repository.SongRepository
}

func NewGetSongFiltersMetadata(
	jwtService service.JwtService,
	songRepository repository.SongRepository,
) GetSongFiltersMetadata {
	return GetSongFiltersMetadata{
		jwtService:     jwtService,
		songRepository: songRepository,
	}
}

func (g GetSongFiltersMetadata) Handle(
	request requests.GetSongFiltersMetadataRequest,
	token string,
) (metadata model.SongFiltersMetadata, e *httperror.ErrorCode) {
	userID, errCode := g.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return metadata, errCode
	}

	err := g.songRepository.GetFiltersMetadata(&metadata, userID, request.SearchBy)
	if err != nil {
		return metadata, httperror.DatabaseError(err)
	}
	return metadata, nil
}
