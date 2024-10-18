package playlist

import (
	"repertoire/data/repository"
	"repertoire/utils/wrapper"

	"github.com/google/uuid"
)

type DeletePlaylist struct {
	repository repository.PlaylistRepository
}

func NewDeletePlaylist(repository repository.PlaylistRepository) DeletePlaylist {
	return DeletePlaylist{repository: repository}
}

func (d DeletePlaylist) Handle(id uuid.UUID) *wrapper.ErrorCode {
	err := d.repository.Delete(id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	return nil
}
