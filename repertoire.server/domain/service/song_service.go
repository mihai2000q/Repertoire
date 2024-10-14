package service

import (
	"errors"
	"github.com/google/uuid"
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/models"
	"repertoire/utils"
)

type SongService struct {
	repository repository.SongRepository
	jwtService service.JwtService
}

func NewSongService(repository repository.SongRepository, jwtService service.JwtService) SongService {
	return SongService{
		repository: repository,
		jwtService: jwtService,
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

func (s *SongService) GetAll(request requests.GetSongsRequest) (songs []models.Song, e *utils.ErrorCode) {
	err := s.repository.GetAllByUser(songs, request.UserId)
	if err != nil {
		return songs, utils.InternalServerError(err)
	}
	return songs, nil
}

func (s *SongService) Create(request requests.CreateSongRequest, token string) *utils.ErrorCode {
	userId, errCode := s.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	song := models.Song{
		ID:         uuid.New(),
		Title:      request.Title,
		IsRecorded: request.IsRecorded,
		UserID:     userId,
	}
	err := s.repository.Create(&song)
	if err != nil {
		return utils.InternalServerError(err)
	}
	return nil
}

func (s *SongService) Update(request requests.UpdateSongRequest) *utils.ErrorCode {
	var song models.Song
	err := s.repository.Get(&song, request.Id)
	if err != nil {
		return utils.InternalServerError(err)
	}
	if song.ID == uuid.Nil {
		return utils.NotFoundError(errors.New("song not found"))
	}

	song.Title = request.Title
	song.IsRecorded = request.IsRecorded

	err = s.repository.Update(&song)
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
