package album

import (
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/models"
	"repertoire/utils/wrapper"
)

type GetAllAlbums struct {
	repository repository.AlbumRepository
}

func NewGetAllAlbums(repository repository.AlbumRepository) GetAllAlbums {
	return GetAllAlbums{
		repository: repository,
	}
}

func (g GetAllAlbums) Handle(request requests.GetAlbumsRequest) (result wrapper.WithTotalCount[models.Album], e *wrapper.ErrorCode) {
	err := g.repository.GetAllByUser(&result.Data, request.UserID, request.CurrentPage, request.PageSize)
	if err != nil {
		return result, wrapper.InternalServerError(err)
	}
	err = g.repository.GetAllByUserCount(&result.TotalCount, request.UserID)
	if err != nil {
		return result, wrapper.InternalServerError(err)
	}
	return result, nil
}
