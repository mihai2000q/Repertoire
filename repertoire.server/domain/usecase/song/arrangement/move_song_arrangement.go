package arrangement

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/reorder"
	"repertoire/server/internal/wrapper"
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

func (c MoveSongArrangement) Handle(request requests.MoveSongArrangementRequest) *wrapper.ErrorCode {
	var song model.Song
	err := c.songRepository.GetWithArrangements(&song, request.SongID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return wrapper.NotFoundError(errors.New("song not found"))
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

	err = c.songRepository.UpdateWithAssociations(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
