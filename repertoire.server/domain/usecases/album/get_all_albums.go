package album

import (
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/models"
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

func (g GetAllAlbums) Handle(request requests.GetAlbumsRequest, token string) (result wrapper.WithTotalCount[models.Album], e *wrapper.ErrorCode) {
	userId, errCode := g.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return result, errCode
	}

	err := g.repository.GetAllByUser(&result.Data, userId, request.CurrentPage, request.PageSize)
	if err != nil {
		return result, wrapper.InternalServerError(err)
	}

	err = g.repository.GetAllByUserCount(&result.TotalCount, userId)
	if err != nil {
		return result, wrapper.InternalServerError(err)
	}

	return result, nil
}
