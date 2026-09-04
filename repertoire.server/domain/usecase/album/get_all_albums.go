package album

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/pagination"
	"repertoire/server/model"
)

type GetAllAlbums struct {
	albumRepository repository.AlbumRepository
	jwtService      service.JwtService
}

func NewGetAllAlbums(
	albumRepository repository.AlbumRepository,
	jwtService service.JwtService,
) GetAllAlbums {
	return GetAllAlbums{
		albumRepository: albumRepository,
		jwtService:      jwtService,
	}
}

func (g GetAllAlbums) Handle(request requests.GetAlbumsRequest, token string) (result pagination.WithTotalCount[model.EnhancedAlbum], e *httperror.ErrorCode) {
	userID, errCode := g.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return result, errCode
	}

	err := g.albumRepository.GetAllByUser(
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

	err = g.albumRepository.GetAllByUserCount(&result.TotalCount, userID, request.SearchBy)
	if err != nil {
		return result, httperror.DatabaseError(err)
	}

	return result, nil
}
