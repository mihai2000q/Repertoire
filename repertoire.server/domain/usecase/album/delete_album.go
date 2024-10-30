package album

import (
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"

	"github.com/google/uuid"
)

type DeleteAlbum struct {
	repository repository.AlbumRepository
}

func NewDeleteAlbum(repository repository.AlbumRepository) DeleteAlbum {
	return DeleteAlbum{repository: repository}
}

func (d DeleteAlbum) Handle(id uuid.UUID) *wrapper.ErrorCode {
	err := d.repository.Delete(id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	return nil
}
