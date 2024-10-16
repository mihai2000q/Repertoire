package album

import (
	"repertoire/data/repository"
	"repertoire/utils"

	"github.com/google/uuid"
)

type DeleteAlbum struct {
	repository repository.AlbumRepository
}

func NewDeleteAlbum(repository repository.AlbumRepository) DeleteAlbum {
	return DeleteAlbum{repository: repository}
}

func (d DeleteAlbum) Handle(id uuid.UUID) *utils.ErrorCode {
	err := d.repository.Delete(id)
	if err != nil {
		return utils.InternalServerError(err)
	}
	return nil
}
