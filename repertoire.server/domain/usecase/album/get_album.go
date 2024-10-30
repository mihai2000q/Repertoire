package album

import (
	"errors"
	"repertoire/server/data/repository"
	"repertoire/server/model"
	"repertoire/server/utils/wrapper"

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
	err := g.repository.GetWithAssociations(&album, id)
	if err != nil {
		return album, wrapper.InternalServerError(err)
	}
	if album.ID == uuid.Nil {
		return album, wrapper.NotFoundError(errors.New("album not found"))
	}
	return album, nil
}
