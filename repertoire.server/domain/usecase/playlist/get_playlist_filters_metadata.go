package playlist

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type GetPlaylistFiltersMetadata struct {
	jwtService service.JwtService
	repository repository.PlaylistRepository
}

func NewGetPlaylistFiltersMetadata(
	jwtService service.JwtService,
	repository repository.PlaylistRepository,
) GetPlaylistFiltersMetadata {
	return GetPlaylistFiltersMetadata{
		jwtService: jwtService,
		repository: repository,
	}
}

func (g GetPlaylistFiltersMetadata) Handle(
	request requests.GetPlaylistFiltersMetadataRequest,
	token string,
) (metadata model.PlaylistFiltersMetadata, e *wrapper.ErrorCode) {
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
