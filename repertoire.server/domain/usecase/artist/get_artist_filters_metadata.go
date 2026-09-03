package artist

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
)

type GetArtistFiltersMetadata struct {
	jwtService service.JwtService
	repository repository.ArtistRepository
}

func NewGetArtistFiltersMetadata(
	jwtService service.JwtService,
	repository repository.ArtistRepository,
) GetArtistFiltersMetadata {
	return GetArtistFiltersMetadata{
		jwtService: jwtService,
		repository: repository,
	}
}

func (g GetArtistFiltersMetadata) Handle(
	request requests.GetArtistFiltersMetadataRequest,
	token string,
) (metadata model.ArtistFiltersMetadata, e *httperror.ErrorCode) {
	userID, errCode := g.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return metadata, errCode
	}

	if err := g.repository.GetFiltersMetadata(&metadata, userID, request.SearchBy); err != nil {
		return metadata, httperror.DatabaseError(err)
	}
	return metadata, nil
}
