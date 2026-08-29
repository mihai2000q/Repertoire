package section

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type CreateSongSection struct {
	songSectionRepository repository.SongSectionRepository
}

func NewCreateSongSection(
	songSectionRepository repository.SongSectionRepository,
) CreateSongSection {
	return CreateSongSection{
		songSectionRepository: songSectionRepository,
	}
}

func (c CreateSongSection) Handle(request requests.CreateSongSectionRequest) *wrapper.ErrorCode {
	var sectionsCount int64
	err := c.songSectionRepository.CountAllBySong(&sectionsCount, request.SongID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	section := model.SongSection{
		ID:                uuid.New(),
		Name:              request.Name,
		SongSectionTypeID: request.TypeID,
		Order:             uint(sectionsCount),
		SongID:            request.SongID,
	}
	err = c.songSectionRepository.Create(&section)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
