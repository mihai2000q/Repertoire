package song

import (
	"errors"
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/model"
	"repertoire/utils/wrapper"

	"github.com/google/uuid"
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
	if song.ID == uuid.Nil {
		return wrapper.NotFoundError(errors.New("song not found"))
	}

	song.Title = request.Title
	song.Description = request.Description
	song.IsRecorded = request.IsRecorded
	song.Bpm = request.Bpm
	song.SongsterrLink = request.SongsterrLink
	song.ReleaseDate = request.ReleaseDate
	song.Difficulty = request.Difficulty
	song.GuitarTuningID = request.GuitarTuningID
	song.AlbumID = request.AlbumID
	song.ArtistID = request.ArtistID

	err = u.repository.Update(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
