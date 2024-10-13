package service

import (
	"errors"
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

func (s *SongService) Get(id uuid.UUID) (song models.Song, e *utils.ErrorCode) {
	err := s.repository.Get(&song, id)
	if err != nil {
		return song, utils.InternalServerError(err)
	}
	if song.ID == uuid.Nil {
		return song, utils.NotFoundError(errors.New("song not found"))
	}
	return song, nil
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

func (s *SongService) Delete(id uuid.UUID) *utils.ErrorCode {
	var song models.Song
	err := s.repository.Get(&song, id)
	if err != nil {
		return utils.InternalServerError(err)
	}
	if song.ID == uuid.Nil {
		return utils.NotFoundError(errors.New("song not found"))
	}

	err = s.repository.Delete(&song)
	if err != nil {
		return utils.InternalServerError(err)
	}

	return nil
}
