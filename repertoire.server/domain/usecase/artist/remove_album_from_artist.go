package artist

import (
	"errors"
	"github.com/google/uuid"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type RemoveAlbumFromArtist struct {
	albumRepository repository.AlbumRepository
}

func NewRemoveAlbumFromArtist(albumRepository repository.AlbumRepository) RemoveAlbumFromArtist {
	return RemoveAlbumFromArtist{albumRepository: albumRepository}
}

func (r RemoveAlbumFromArtist) Handle(id uuid.UUID, albumID uuid.UUID) *wrapper.ErrorCode {
	var album model.Album
	err := r.albumRepository.GetWithSongs(&album, albumID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(album).IsZero() {
		return wrapper.NotFoundError(errors.New("album not found"))
	}
	if album.ArtistID == nil || *album.ArtistID != id {
		return wrapper.BadRequestError(errors.New("album is not owned by this artist"))
	}

	album.ArtistID = nil
	for i := range album.Songs {
		album.Songs[i].ArtistID = nil
	}

	err = r.albumRepository.UpdateWithAssociations(&album)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	return nil
}
