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

func (c CreateSongArrangement) Handle(request requests.CreateSongArrangementRequest) *wrapper.ErrorCode {
	var arrangementsCount int64
	err := c.songArrangementRepository.CountBySong(&arrangementsCount, request.SongID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	var song model.Song
	err = c.songRepository.GetWithSections(&song, request.SongID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return wrapper.NotFoundError(errors.New("song not found"))
	}

	arrangement := model.SongArrangement{
		ID:     uuid.New(),
		Name:   request.Name,
		Order:  uint(arrangementsCount),
		SongID: request.SongID,
	}

	c.CreateSectionOccurrences(&arrangement, song.Sections)

	err = c.songArrangementRepository.Create(&arrangement)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}

func (c CreateSongArrangement) CreateSectionOccurrences(arrangement *model.SongArrangement, sections []model.SongSection) {
	var occurrences []model.SongSectionOccurrences
	for _, section := range sections {
		occurrence := model.SongSectionOccurrences{
			ID:            uuid.New(),
			Occurrences:   0,
			Section:       section,
			ArrangementID: arrangement.ID,
		}
		occurrences = append(occurrences, occurrence)
	}
	arrangement.Occurrences = occurrences
}
