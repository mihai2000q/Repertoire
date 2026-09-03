package album

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
)

type GetAlbumFiltersMetadata struct {
	jwtService      service.JwtService
	albumRepository repository.AlbumRepository
}

func NewGetAlbumFiltersMetadata(
	jwtService service.JwtService,
	albumRepository repository.AlbumRepository,
) GetAlbumFiltersMetadata {
	return GetAlbumFiltersMetadata{
		jwtService:      jwtService,
		albumRepository: albumRepository,
	}
}

func (g GetAlbumFiltersMetadata) Handle(
	request requests.GetAlbumFiltersMetadataRequest,
	token string,
) (metadata model.AlbumFiltersMetadata, e *httperror.ErrorCode) {
	userID, errCode := g.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return metadata, errCode
	}

	if err := g.albumRepository.GetFiltersMetadata(&metadata, userID, request.SearchBy); err != nil {
		return metadata, httperror.DatabaseError(err)
	}
	return metadata, nil
}
