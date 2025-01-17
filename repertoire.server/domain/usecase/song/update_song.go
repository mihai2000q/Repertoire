package song

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type UpdateSong struct {
	repository repository.SongRepository
}

func NewUpdateSong(repository repository.SongRepository) UpdateSong {
	return UpdateSong{repository: repository}
}

func (u UpdateSong) Handle(request requests.UpdateSongRequest) *wrapper.ErrorCode {
	var song model.Song
	err := u.repository.Get(&song, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return wrapper.NotFoundError(errors.New("song not found"))
	}

	song.Title = request.Title
	song.Description = request.Description
	song.IsRecorded = request.IsRecorded
	song.Bpm = request.Bpm
	song.SongsterrLink = request.SongsterrLink
	song.YoutubeLink = request.YoutubeLink
	song.ReleaseDate = request.ReleaseDate
	song.Difficulty = request.Difficulty
	song.GuitarTuningID = request.GuitarTuningID

	err = u.repository.Update(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
