package section

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/reorder"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type MoveSongSection struct {
	songRepository repository.SongRepository
}

func NewMoveSongSection(repository repository.SongRepository) MoveSongSection {
	return MoveSongSection{
		songRepository: repository,
	}
}

func (c MoveSongSection) Handle(request requests.MoveSongSectionRequest) *wrapper.ErrorCode {
	var song model.Song
	err := c.songRepository.GetWithSections(&song, request.SongID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return wrapper.NotFoundError(errors.New("song not found"))
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

	err = c.songRepository.UpdateWithAssociations(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
