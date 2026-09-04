package album

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
)

type GetAlbum struct {
	albumRepository repository.AlbumRepository
}

func NewGetAlbum(albumRepository repository.AlbumRepository) GetAlbum {
	return GetAlbum{
		albumRepository: albumRepository,
	}
}

func (g GetAlbum) Handle(request requests.GetAlbumRequest) (album model.Album, e *httperror.ErrorCode) {
	if len(request.SongsOrderBy) == 0 {
		request.SongsOrderBy = []string{"album_track_no"}
	}

	if err := g.albumRepository.GetWithAssociations(&album, request.ID, request.SongsOrderBy); err != nil {
		return album, httperror.DatabaseError(err)
	}
	if reflect.ValueOf(album).IsZero() {
		return album, httperror.NotFoundError(errors.New("album not found"))
	}
	return album, nil
}
