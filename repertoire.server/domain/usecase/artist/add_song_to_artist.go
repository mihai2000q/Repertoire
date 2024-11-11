package artist

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type AddSongToArtist struct {
	songRepository repository.SongRepository
}

func NewAddSongToArtist(songRepository repository.SongRepository) AddSongToArtist {
	return AddSongToArtist{songRepository: songRepository}
}

func (a AddSongToArtist) Handle(request requests.AddSongToArtistRequest) *wrapper.ErrorCode {
	var song model.Song
	err := a.songRepository.GetWithSongs(&song, request.SongID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return wrapper.NotFoundError(errors.New("song not found"))
	}
	if song.ArtistID != nil {
		return wrapper.BadRequestError(errors.New("song already has an artist"))
	}

	// update the whole album's artist, including the other songs
	if song.Album != nil {
		song.Album.ArtistID = &request.ID
		for i := range song.Album.Songs {
			song.Album.Songs[i].ArtistID = &request.ID
		}
	} else {
		song.ArtistID = &request.ID
	}

	err = a.songRepository.UpdateWithAssociations(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	return nil
}
