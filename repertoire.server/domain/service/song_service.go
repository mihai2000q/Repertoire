package service

import (
	"repertoire/api/requests/song"
	"repertoire/data/repository"
	"repertoire/models"
	"repertoire/utils"
	"time"

	"github.com/google/uuid"
)

type SongService struct {
	repository repository.SongRepository
}

func NewSongService(repository repository.SongRepository) SongService {
	return SongService{
		repository: repository,
	}
}

func (s *SongService) Create(request song.CreateSongRequest) *utils.ErrorCode {
	song := models.Song{
		ID:         uuid.New(),
		Title:      request.Title,
		IsRecorded: request.IsRecorded,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
	err := s.repository.Create(&song)
	if err != nil {
		return utils.InternalServerError(err)
	}
	return nil
}
