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
	repository repository.AlbumRepository
}

func NewGetAlbum(repository repository.AlbumRepository) GetAlbum {
	return GetAlbum{
		repository: repository,
	}
}

func (g GetAlbum) Handle(request requests.GetAlbumRequest) (album model.Album, e *httperror.ErrorCode) {
	if len(request.SongsOrderBy) == 0 {
		request.SongsOrderBy = []string{"album_track_no"}
	}

	if err := g.repository.GetWithAssociations(&album, request.ID, request.SongsOrderBy); err != nil {
		return album, httperror.DatabaseError(err)
	}
	if reflect.ValueOf(album).IsZero() {
		return album, httperror.NotFoundError(errors.New("album not found"))
	}
	return album, nil
}
