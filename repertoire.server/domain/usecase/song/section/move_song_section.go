package section

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/reorder"
	"repertoire/server/model"
)

type MoveSongSection struct {
	songRepository repository.SongRepository
}

func NewMoveSongSection(songRepository repository.SongRepository) MoveSongSection {
	return MoveSongSection{
		songRepository: songRepository,
	}
}

func (c MoveSongSection) Handle(request requests.MoveSongSectionRequest) *httperror.ErrorCode {
	var song model.Song
	if err := c.songRepository.GetWithSections(&song, request.SongID); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return httperror.NotFoundError(errors.New("song not found"))
	}

	errCode := reorder.MoveEntity(
		song.Sections,
		request.ID,
		request.OverID,
		&reorder.Config{
			EntityNotFoundMsg:     "section not found",
			OverEntityNotFoundMsg: "over section not found",
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
