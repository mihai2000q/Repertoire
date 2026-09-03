package arrangement

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/reorder"
	"repertoire/server/model"
)

type MoveSongArrangement struct {
	songRepository repository.SongRepository
}

func NewMoveSongArrangement(repository repository.SongRepository) MoveSongArrangement {
	return MoveSongArrangement{
		songRepository: repository,
	}
}

func (c MoveSongArrangement) Handle(request requests.MoveSongArrangementRequest) *httperror.ErrorCode {
	var song model.Song
	if err := c.songRepository.GetWithArrangements(&song, request.SongID); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return httperror.NotFoundError(errors.New("song not found"))
	}

	errCode := reorder.MoveEntity(
		song.Arrangements,
		request.ID,
		request.OverID,
		&reorder.Config{
			EntityNotFoundMsg:     "arrangement not found",
			OverEntityNotFoundMsg: "over arrangement not found",
		},
	)
	if errCode != nil {
		return errCode
	}

	if err := c.songRepository.UpdateWithAssociations(&song); err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
