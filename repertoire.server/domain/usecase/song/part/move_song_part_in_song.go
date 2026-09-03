package part

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type MoveSongPartInSong struct {
	songRepository repository.SongRepository
}

func NewMoveSongPartInSong(songRepository repository.SongRepository) MoveSongPartInSong {
	return MoveSongPartInSong{
		songRepository: songRepository,
	}
}

func (c MoveSongPartInSong) Handle(request requests.MoveSongPartInSongRequest) *httperror.ErrorCode {
	var song model.Song
	err := c.songRepository.GetWithParts(&song, request.SongID)
	if err != nil {
		return httperror.InternalServerError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return httperror.NotFoundError(errors.New("song not found"))
	}

	index, overIndex, err := c.getIndexes(song.Parts, request.ID, request.OverID)
	if err != nil {
		return httperror.NotFoundError(err)
	}
	song.Parts = c.move(song.Parts, index, overIndex)

	err = c.songRepository.UpdateWithAssociations(&song)
	if err != nil {
		return httperror.InternalServerError(err)
	}

	return nil
}

func (c MoveSongPartInSong) getIndexes(parts []model.SongPart, id uuid.UUID, overID uuid.UUID) (int, int, error) {
	var index *int
	var overIndex *int
	for i := 0; i < len(parts); i++ {
		if parts[i].ID == id {
			index = &i
		} else if parts[i].ID == overID {
			overIndex = &i
		}
	}

	if index == nil {
		return -1, -1, errors.New("part not found")
	}
	if overIndex == nil {
		return -1, -1, errors.New("over part not found")
	}

	return *index, *overIndex, nil
}

func (c MoveSongPartInSong) move(parts []model.SongPart, index int, overIndex int) []model.SongPart {
	if index < overIndex {
		for i := index + 1; i <= overIndex; i++ {
			parts[i].SongOrder = uint(i - 1)
		}
	} else {
		for i := overIndex; i <= index; i++ {
			parts[i].SongOrder = uint(i + 1)
		}
	}

	parts[index].SongOrder = uint(overIndex)

	return parts
}
