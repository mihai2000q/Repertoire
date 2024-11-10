package album

import (
	"errors"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"slices"

	"github.com/google/uuid"
)

type RemoveSongFromAlbum struct {
	repository repository.AlbumRepository
}

func NewRemoveSongFromAlbum(repository repository.AlbumRepository) RemoveSongFromAlbum {
	return RemoveSongFromAlbum{repository: repository}
}

func (r RemoveSongFromAlbum) Handle(id uuid.UUID, songID uuid.UUID) *wrapper.ErrorCode {
	var album model.Album
	err := r.repository.GetWithSongs(&album, id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(album).IsZero() {
		return wrapper.NotFoundError(errors.New("album not found"))
	}

	index := slices.IndexFunc(album.Songs, func(s model.Song) bool {
		return s.ID == songID
	})
	if index == -1 {
		return wrapper.NotFoundError(errors.New("song not found"))
	}

	for i := index + 1; i < len(album.Songs); i++ {
		*album.Songs[i].AlbumTrackNo = *album.Songs[i].AlbumTrackNo - 1
	}

	err = r.repository.UpdateWithAssociations(&album)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = r.repository.RemoveSong(&album, &album.Songs[index])
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
