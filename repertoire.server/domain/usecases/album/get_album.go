package album

import (
	"errors"
	"repertoire/data/repository"
	"repertoire/models"
	"repertoire/utils"

	"github.com/google/uuid"
)

type GetAlbum struct {
	repository repository.AlbumRepository
}

func NewGetAlbum(repository repository.AlbumRepository) GetAlbum {
	return GetAlbum{
		repository: repository,
	}
}

func (g GetAlbum) Handle(id uuid.UUID) (album models.Album, e *utils.ErrorCode) {
	err := g.repository.Get(&album, id)
	if err != nil {
		return album, utils.InternalServerError(err)
	}
	if album.ID == uuid.Nil {
		return album, utils.NotFoundError(errors.New("album not found"))
	}
	return album, nil
}
