package song

import (
	"errors"
	"github.com/google/uuid"
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/models"
	"repertoire/utils"
)

type UpdateSong struct {
	repository repository.SongRepository
}

func NewUpdateSong(repository repository.SongRepository) UpdateSong {
	return UpdateSong{repository: repository}
}

func (u UpdateSong) Handle(request requests.UpdateSongRequest) *utils.ErrorCode {
	var song models.Song
	err := u.repository.Get(&song, request.ID)
	if err != nil {
		return utils.InternalServerError(err)
	}
	if song.ID == uuid.Nil {
		return utils.NotFoundError(errors.New("song not found"))
	}

	song.Title = request.Title
	song.IsRecorded = request.IsRecorded

	err = u.repository.Update(&song)
	if err != nil {
		return utils.InternalServerError(err)
	}

	return nil
}
