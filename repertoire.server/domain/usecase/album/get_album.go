package album

import (
	"errors"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

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
	if reflect.ValueOf(album).IsZero() {
		return album, wrapper.NotFoundError(errors.New("album not found"))
	}
	return album, nil
}
