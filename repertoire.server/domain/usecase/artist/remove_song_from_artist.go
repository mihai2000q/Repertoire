package artist

import (
	"errors"
	"github.com/google/uuid"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type RemoveSongFromArtist struct {
	songRepository repository.SongRepository
}

func NewRemoveSongFromArtist(songRepository repository.SongRepository) RemoveSongFromArtist {
	return RemoveSongFromArtist{songRepository: songRepository}
}

func (r RemoveSongFromArtist) Handle(id uuid.UUID, songID uuid.UUID) *wrapper.ErrorCode {
	var song model.Song
	err := r.songRepository.Get(&song, songID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return wrapper.NotFoundError(errors.New("song not found"))
	}
	if song.ArtistID == nil || *song.ArtistID != id {
		return wrapper.BadRequestError(errors.New("song is not owned by this artist"))
	}

	song.ArtistID = nil

	err = r.songRepository.Update(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	return nil
}
