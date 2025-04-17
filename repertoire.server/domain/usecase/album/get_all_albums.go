package album

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type GetAllAlbums struct {
	repository repository.AlbumRepository
	jwtService service.JwtService
}

func NewGetAllAlbums(repository repository.AlbumRepository, jwtService service.JwtService) GetAllAlbums {
	return GetAllAlbums{
		repository: repository,
		jwtService: jwtService,
	}
}

func (g GetAllAlbums) Handle(request requests.GetAlbumsRequest, token string) (result wrapper.WithTotalCount[model.EnhancedAlbum], e *wrapper.ErrorCode) {
	userID, errCode := g.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return result, errCode
	}

	err := g.repository.GetAllByUser(
		&result.Models,
		userID,
		request.CurrentPage,
		request.PageSize,
		request.OrderBy,
		request.SearchBy,
	)
	if err != nil {
		return result, wrapper.InternalServerError(err)
	}

	err = g.repository.GetAllByUserCount(&result.TotalCount, userID, request.SearchBy)
	if err != nil {
		return result, wrapper.InternalServerError(err)
	}

	return result, nil
}
