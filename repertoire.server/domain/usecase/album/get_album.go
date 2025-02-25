package album

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
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

func (g GetAlbum) Handle(request requests.GetAlbumRequest) (album model.Album, e *wrapper.ErrorCode) {
	if len(request.SongsOrderBy) == 0 {
		request.SongsOrderBy = []string{"album_track_no"}
	}

	err := g.repository.GetWithAssociations(&album, request.ID, request.SongsOrderBy)
	if err != nil {
		return album, wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(album).IsZero() {
		return album, wrapper.NotFoundError(errors.New("album not found"))
	}
	return album, nil
}
