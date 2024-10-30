package section

import (
	"github.com/google/uuid"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/model"
	"repertoire/server/utils/wrapper"
)

type CreateSongSection struct {
	songRepository repository.SongRepository
}

func NewCreateSongSection(repository repository.SongRepository) CreateSongSection {
	return CreateSongSection{
		songRepository: repository,
	}
}

func (c CreateSongSection) Handle(request requests.CreateSongSectionRequest) *wrapper.ErrorCode {
	var count int64
	err := c.songRepository.CountSectionsBySong(&count, request.SongID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	section := model.SongSection{
		ID:                uuid.New(),
		Name:              request.Name,
		SongSectionTypeID: request.TypeID,
		Order:             uint(count),
		SongID:            request.SongID,
	}
	err = c.songRepository.CreateSection(&section)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
