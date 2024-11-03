package album

import (
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

	index := slices.IndexFunc(album.Songs, func(s model.Song) bool {
		return s.ID == songID
	})

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
