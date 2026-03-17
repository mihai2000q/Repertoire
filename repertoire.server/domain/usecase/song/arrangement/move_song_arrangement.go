package arrangement

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
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

	index, overIndex, err := c.getIndexes(song.Arrangements, request.ID, request.OverID)
	if err != nil {
		return wrapper.NotFoundError(err)
	}
	song.Arrangements = c.move(song.Arrangements, index, overIndex)

	err = c.songRepository.UpdateWithAssociations(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}

func (c MoveSongArrangement) getIndexes(arrangements []model.SongArrangement, id uuid.UUID, overID uuid.UUID) (int, int, error) {
	var index *int
	var overIndex *int
	for i := 0; i < len(arrangements); i++ {
		if arrangements[i].ID == id {
			index = &i
		} else if arrangements[i].ID == overID {
			overIndex = &i
		}
	}

	if index == nil {
		return -1, -1, errors.New("arrangement not found")
	}
	if overIndex == nil {
		return -1, -1, errors.New("over arrangement not found")
	}

	return *index, *overIndex, nil
}

func (c MoveSongArrangement) move(arrangements []model.SongArrangement, index int, overIndex int) []model.SongArrangement {
	if index < overIndex {
		for i := index + 1; i <= overIndex; i++ {
			arrangements[i].Order = uint(i - 1)
		}
	} else {
		for i := overIndex; i <= index; i++ {
			arrangements[i].Order = uint(i + 1)
		}
	}

	arrangements[index].Order = uint(overIndex)

	return arrangements
}
