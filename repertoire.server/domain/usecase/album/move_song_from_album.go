package album

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/reorder"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type MoveSongFromAlbum struct {
	repository repository.AlbumRepository
}

func NewMoveSongFromAlbum(repository repository.AlbumRepository) MoveSongFromAlbum {
	return MoveSongFromAlbum{repository: repository}
}

func (m MoveSongFromAlbum) Handle(request requests.MoveSongFromAlbumRequest) *wrapper.ErrorCode {
	var album model.Album
	err := m.repository.GetWithSongs(&album, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(album).IsZero() {
		return wrapper.NotFoundError(errors.New("album not found"))
	}

	errCode := reorder.MoveEntity(
		album.Songs,
		request.SongID,
		request.OverSongID,
		&reorder.Config{
			OrderField:            "AlbumTrackNo",
			StartOffset:           1,
			EntityNotFoundMsg:     "song not found",
			OverEntityNotFoundMsg: "over song not found",
		},
	)
	if errCode != nil {
		return errCode
	}

	err = m.repository.UpdateWithAssociations(&album)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
