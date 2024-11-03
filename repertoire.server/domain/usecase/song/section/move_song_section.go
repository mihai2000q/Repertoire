package section

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
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

	index, overIndex, err := c.getIndexes(song.Sections, request.ID, request.OverID)
	if err != nil {
		return wrapper.NotFoundError(err)
	}
	song.Sections = c.move(song.Sections, index, overIndex)

	err = c.songRepository.UpdateWithAssociations(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}

func (c MoveSongSection) getIndexes(sections []model.SongSection, id uuid.UUID, overID uuid.UUID) (int, int, error) {
	var index *int
	var overIndex *int
	for i := 0; i < len(sections); i++ {
		if sections[i].ID == id {
			index = &i
		} else if sections[i].ID == overID {
			overIndex = &i
		}
	}

	if index == nil {
		return -1, -1, errors.New("section not found")
	}
	if overIndex == nil {
		return -1, -1, errors.New("over section not found")
	}

	return *index, *overIndex, nil
}

func (c MoveSongSection) move(sections []model.SongSection, index int, overIndex int) []model.SongSection {
	if index < overIndex {
		for i := index + 1; i <= overIndex; i++ {
			sections[i].Order = uint(i - 1)
		}
	} else {
		for i := overIndex; i <= index; i++ {
			sections[i].Order = uint(i + 1)
		}
	}

	sections[index].Order = uint(overIndex)

	return sections
}
