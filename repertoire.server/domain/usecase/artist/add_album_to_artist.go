package artist

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type AddAlbumToArtist struct {
	albumRepository repository.AlbumRepository
}

func NewAddAlbumToArtist(albumRepository repository.AlbumRepository) AddAlbumToArtist {
	return AddAlbumToArtist{albumRepository: albumRepository}
}

func (a AddAlbumToArtist) Handle(request requests.AddAlbumToArtistRequest) *wrapper.ErrorCode {
	var album model.Album
	err := a.albumRepository.GetWithSongs(&album, request.AlbumID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(album).IsZero() {
		return wrapper.NotFoundError(errors.New("album not found"))
	}
	if album.ArtistID != nil {
		return wrapper.BadRequestError(errors.New("album already has an artist"))
	}

	// update the whole album's artist
	album.ArtistID = &request.ID
	for i := range album.Songs {
		album.Songs[i].ArtistID = &request.ID
	}

	err = a.albumRepository.UpdateWithAssociations(&album)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	return nil
}
