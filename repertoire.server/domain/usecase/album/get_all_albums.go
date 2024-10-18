package album

import (
	"repertoire/api/request"
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

func (g GetAllAlbums) Handle(request request.GetAlbumsRequest, token string) (result wrapper.WithTotalCount[model.Album], e *wrapper.ErrorCode) {
	userId, errCode := g.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return result, errCode
	}

	err := g.repository.GetAllByUser(&result.Models, userId, request.CurrentPage, request.PageSize)
	if err != nil {
		return result, wrapper.InternalServerError(err)
	}

	err = g.repository.GetAllByUserCount(&result.TotalCount, userId)
	if err != nil {
		return result, wrapper.InternalServerError(err)
	}

	return result, nil
}
