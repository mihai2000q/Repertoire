package artist

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
)

type GetArtistFiltersMetadata struct {
	jwtService       service.JwtService
	artistRepository repository.ArtistRepository
}

func NewGetArtistFiltersMetadata(
	jwtService service.JwtService,
	artistRepository repository.ArtistRepository,
) GetArtistFiltersMetadata {
	return GetArtistFiltersMetadata{
		jwtService:       jwtService,
		artistRepository: artistRepository,
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

	if err := g.artistRepository.GetFiltersMetadata(&metadata, userID, request.SearchBy); err != nil {
		return metadata, httperror.DatabaseError(err)
	}
	return metadata, nil
}
