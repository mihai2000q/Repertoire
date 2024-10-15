package song

import (
	"github.com/google/uuid"
	"repertoire/data/repository"
	"repertoire/utils"
)

type DeleteSong struct {
	repository repository.SongRepository
}

func NewDeleteSong(repository repository.SongRepository) *DeleteSong {
	return &DeleteSong{repository: repository}
}

func (d DeleteSong) Handle(id uuid.UUID) *utils.ErrorCode {
	err := d.repository.Delete(id)
	if err != nil {
		return utils.InternalServerError(err)
	}
	return nil
}
