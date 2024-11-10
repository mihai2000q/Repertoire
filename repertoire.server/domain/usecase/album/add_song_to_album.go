package album

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type AddSongToAlbum struct {
	repository     repository.AlbumRepository
	songRepository repository.SongRepository
}

func NewAddSongToAlbum(
	albumRepository repository.AlbumRepository,
	repository repository.SongRepository,
) AddSongToAlbum {
	return AddSongToAlbum{
		repository:     albumRepository,
		songRepository: repository,
	}
}

func (a AddSongToAlbum) Handle(request requests.AddSongToAlbumRequest) *wrapper.ErrorCode {
	var song model.Song
	err := a.songRepository.Get(&song, request.SongID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return wrapper.NotFoundError(errors.New("song not found"))
	}
	if song.AlbumID != nil {
		return wrapper.BadRequestError(errors.New("song already has an album"))
	}

	var album model.Album
	err = a.repository.GetWithSongs(&album, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(album).IsZero() {
		return wrapper.NotFoundError(errors.New("album not found"))
	}
	if song.ArtistID != album.ArtistID {
		return wrapper.BadRequestError(errors.New("song and album do not share the same artist"))
	}

	song.AlbumID = &request.ID
	trackNo := uint(len(album.Songs)) + 1
	song.AlbumTrackNo = &trackNo

	err = a.songRepository.Update(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
