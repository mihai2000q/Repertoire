package album

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type AddSongsToAlbum struct {
	repository     repository.AlbumRepository
	songRepository repository.SongRepository
}

func NewAddSongsToAlbum(
	albumRepository repository.AlbumRepository,
	repository repository.SongRepository,
) AddSongsToAlbum {
	return AddSongsToAlbum{
		repository:     albumRepository,
		songRepository: repository,
	}
}

func (a AddSongsToAlbum) Handle(request requests.AddSongsToAlbumRequest) *wrapper.ErrorCode {
	var songs []model.Song
	err := a.songRepository.GetAllByIDs(&songs, request.SongIDs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	var album model.Album
	err = a.repository.GetWithSongs(&album, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(album).IsZero() {
		return wrapper.NotFoundError(errors.New("album not found"))
	}

	songsLength := len(album.Songs) + 1
	for i, song := range songs {
		if song.ArtistID != nil && album.ArtistID != nil && song.ArtistID != album.ArtistID {
			return wrapper.BadRequestError(errors.New("song " + song.ID.String() + " and album do not share the same artist"))
		}
		if song.AlbumID != nil {
			return wrapper.BadRequestError(errors.New("song " + song.ID.String() + " already has an album"))
		}

		song.AlbumID = &request.ID
		trackNo := uint(songsLength + i)
		song.AlbumTrackNo = &trackNo

		// synchronize artist
		song.ArtistID, album.ArtistID = album.ArtistID, song.ArtistID

		err = a.songRepository.Update(&song)
		if err != nil {
			return wrapper.InternalServerError(err)
		}
	}

	return nil
}
