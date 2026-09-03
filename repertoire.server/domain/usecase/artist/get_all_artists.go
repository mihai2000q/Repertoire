package artist

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/pagination"
	"repertoire/server/model"
)

type GetAllArtists struct {
	artistRepository repository.ArtistRepository
	jwtService       service.JwtService
}

func NewGetAllArtists(
	artistRepository repository.ArtistRepository,
	jwtService service.JwtService,
) GetAllArtists {
	return GetAllArtists{
		artistRepository: artistRepository,
		jwtService:       jwtService,
	}
}

func (g GetAllArtists) Handle(request requests.GetArtistsRequest, token string) (result pagination.WithTotalCount[model.EnhancedArtist], e *httperror.ErrorCode) {
	userID, errCode := g.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return result, errCode
	}

	err := g.artistRepository.GetAllByUser(
		&result.Models,
		userID,
		request.CurrentPage,
		request.PageSize,
		request.OrderBy,
		request.SearchBy,
	)
	if err != nil {
		return result, httperror.DatabaseError(err)
	}

	err = g.artistRepository.GetAllByUserCount(&result.TotalCount, userID, request.SearchBy)
	if err != nil {
		return result, httperror.DatabaseError(err)
	}

	return result, nil
}
