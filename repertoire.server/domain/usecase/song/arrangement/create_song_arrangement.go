package arrangement

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type CreateSongArrangement struct {
	songArrangementRepository repository.SongArrangementRepository
	songRepository            repository.SongRepository
}

func NewCreateSongArrangement(
	songArrangementRepository repository.SongArrangementRepository,
	songRepository repository.SongRepository,
) CreateSongArrangement {
	return CreateSongArrangement{
		songArrangementRepository: songArrangementRepository,
		songRepository:            songRepository,
	}
}

func (c CreateSongArrangement) Handle(request requests.CreateSongArrangementRequest) (uuid.UUID, *httperror.ErrorCode) {
	var arrangementsCount int64
	if err := c.songArrangementRepository.CountBySong(&arrangementsCount, request.SongID); err != nil {
		return uuid.Nil, httperror.DatabaseError(err)
	}

	var song model.Song
	if err := c.songRepository.GetWithParts(&song, request.SongID); err != nil {
		return uuid.Nil, httperror.DatabaseError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return uuid.Nil, httperror.NotFoundError(errors.New("song not found"))
	}

	arrangement := model.SongArrangement{
		ID:     uuid.New(),
		Name:   request.Name,
		Order:  uint(arrangementsCount),
		SongID: request.SongID,
	}

	c.CreatePartOccurrences(&arrangement, song.Parts)

	if err := c.songArrangementRepository.Create(&arrangement); err != nil {
		return uuid.Nil, httperror.DatabaseError(err)
	}

	return arrangement.ID, nil
}

func (c CreateSongArrangement) CreatePartOccurrences(arrangement *model.SongArrangement, parts []model.SongPart) {
	var occurrences []model.SongPartOccurrences
	for _, part := range parts {
		occurrence := model.SongPartOccurrences{
			Occurrences:   0,
			PartID:        part.ID,
			ArrangementID: arrangement.ID,
		}
		occurrences = append(occurrences, occurrence)
	}
	arrangement.PartOccurrences = occurrences
}
