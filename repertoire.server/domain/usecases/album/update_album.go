package album

import (
	"errors"
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/models"
	"repertoire/utils"

	"github.com/google/uuid"
)

type UpdateAlbum struct {
	repository repository.AlbumRepository
}

func NewUpdateAlbum(repository repository.AlbumRepository) UpdateAlbum {
	return UpdateAlbum{repository: repository}
}

func (u UpdateAlbum) Handle(request requests.UpdateAlbumRequest) *utils.ErrorCode {
	var album models.Album
	err := u.repository.Get(&album, request.ID)
	if err != nil {
		return utils.InternalServerError(err)
	}
	if album.ID == uuid.Nil {
		return utils.NotFoundError(errors.New("album not found"))
	}

	album.Title = request.Title

	err = u.repository.Update(&album)
	if err != nil {
		return utils.InternalServerError(err)
	}

	return nil
}
