package album

import (
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/model"
	"repertoire/utils/wrapper"
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

func (g GetAllAlbums) Handle(request requests.GetAlbumsRequest, token string) (result wrapper.WithTotalCount[model.Album], e *wrapper.ErrorCode) {
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
	)
	if err != nil {
		return result, wrapper.InternalServerError(err)
	}

	err = g.repository.GetAllByUserCount(&result.TotalCount, userID)
	if err != nil {
		return result, wrapper.InternalServerError(err)
	}

	return result, nil
}
