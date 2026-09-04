package album

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/reorder"
	"repertoire/server/model"
)

type MoveSongFromAlbum struct {
	albumRepository repository.AlbumRepository
}

func NewMoveSongFromAlbum(albumRepository repository.AlbumRepository) MoveSongFromAlbum {
	return MoveSongFromAlbum{albumRepository: albumRepository}
}

func (m MoveSongFromAlbum) Handle(request requests.MoveSongFromAlbumRequest) *httperror.ErrorCode {
	var album model.Album
	if err := m.albumRepository.GetWithSongs(&album, request.ID); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(album).IsZero() {
		return httperror.NotFoundError(errors.New("album not found"))
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

	if err := m.albumRepository.UpdateWithAssociations(&album); err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
