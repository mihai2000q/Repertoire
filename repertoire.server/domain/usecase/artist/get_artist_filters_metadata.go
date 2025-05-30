package artist

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
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
) (metadata model.ArtistFiltersMetadata, e *wrapper.ErrorCode) {
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
