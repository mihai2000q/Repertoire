package song

import (
	"errors"
	"github.com/google/uuid"
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/models"
	"repertoire/utils/wrapper"
)

type UpdateSong struct {
	repository repository.SongRepository
}

func NewUpdateSong(repository repository.SongRepository) UpdateSong {
	return UpdateSong{repository: repository}
}

func (u UpdateSong) Handle(request requests.UpdateSongRequest) *wrapper.ErrorCode {
	var song models.Song
	err := u.repository.Get(&song, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if song.ID == uuid.Nil {
		return wrapper.NotFoundError(errors.New("song not found"))
	}

	song.Title = request.Title
	song.IsRecorded = request.IsRecorded

	err = u.repository.Update(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
