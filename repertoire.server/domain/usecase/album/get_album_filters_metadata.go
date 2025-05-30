package album

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type GetAlbumFiltersMetadata struct {
	jwtService service.JwtService
	repository repository.AlbumRepository
}

func NewGetAlbumFiltersMetadata(
	jwtService service.JwtService,
	repository repository.AlbumRepository,
) GetAlbumFiltersMetadata {
	return GetAlbumFiltersMetadata{
		jwtService: jwtService,
		repository: repository,
	}
}

func (g GetAlbumFiltersMetadata) Handle(
	request requests.GetAlbumFiltersMetadataRequest,
	token string,
) (metadata model.AlbumFiltersMetadata, e *wrapper.ErrorCode) {
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
