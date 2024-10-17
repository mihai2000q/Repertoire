package album

import (
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/models"
	"repertoire/utils"
)

type GetAllAlbums struct {
	repository repository.AlbumRepository
}

func NewGetAllAlbums(repository repository.AlbumRepository) GetAllAlbums {
	return GetAllAlbums{
		repository: repository,
	}
}

func (g GetAllAlbums) Handle(request requests.GetAlbumsRequest) (albums []models.Album, e *utils.ErrorCode) {
	err := g.repository.GetAllByUser(&albums, request.UserID, request.CurrentPage, request.PageSize)
	if err != nil {
		return albums, utils.InternalServerError(err)
	}
	return albums, nil
}
