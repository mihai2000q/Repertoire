package album

import (
	"errors"
	"repertoire/data/repository"
	"repertoire/model"
	"repertoire/utils/wrapper"

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

func (g GetAlbum) Handle(id uuid.UUID) (album model.Album, e *wrapper.ErrorCode) {
	err := g.repository.Get(&album, id)
	if err != nil {
		return album, wrapper.InternalServerError(err)
	}
	if album.ID == uuid.Nil {
		return album, wrapper.NotFoundError(errors.New("album not found"))
	}
	return album, nil
}
