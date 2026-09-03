package playlist

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
)

type GetPlaylistFiltersMetadata struct {
	jwtService         service.JwtService
	playlistRepository repository.PlaylistRepository
}

func NewGetPlaylistFiltersMetadata(
	jwtService service.JwtService,
	playlistRepository repository.PlaylistRepository,
) GetPlaylistFiltersMetadata {
	return GetPlaylistFiltersMetadata{
		jwtService:         jwtService,
		playlistRepository: playlistRepository,
	}
}

func (g GetPlaylistFiltersMetadata) Handle(
	request requests.GetPlaylistFiltersMetadataRequest,
	token string,
) (metadata model.PlaylistFiltersMetadata, e *httperror.ErrorCode) {
	userID, errCode := g.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return metadata, errCode
	}

	err := g.playlistRepository.GetFiltersMetadata(&metadata, userID, request.SearchBy)
	if err != nil {
		return metadata, httperror.DatabaseError(err)
	}
	return metadata, nil
}
