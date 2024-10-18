package song

import (
	"github.com/google/uuid"
	"repertoire/data/repository"
	"repertoire/utils/wrapper"
)

type DeleteSong struct {
	repository repository.SongRepository
}

func NewDeleteSong(repository repository.SongRepository) DeleteSong {
	return DeleteSong{repository: repository}
}

func (d DeleteSong) Handle(id uuid.UUID) *wrapper.ErrorCode {
	err := d.repository.Delete(id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	return nil
}
