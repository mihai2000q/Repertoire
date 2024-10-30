package album

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type UpdateAlbum struct {
	repository repository.AlbumRepository
}

func NewUpdateAlbum(repository repository.AlbumRepository) UpdateAlbum {
	return UpdateAlbum{repository: repository}
}

func (u UpdateAlbum) Handle(request requests.UpdateAlbumRequest) *wrapper.ErrorCode {
	var album model.Album
	err := u.repository.Get(&album, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if album.ID == uuid.Nil {
		return wrapper.NotFoundError(errors.New("album not found"))
	}

	album.Title = request.Title

	err = u.repository.Update(&album)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
