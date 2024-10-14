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

type SongService interface {
	Get(id uuid.UUID) (song models.Song, e *utils.ErrorCode)
	GetAll(request requests.GetSongsRequest) (songs []models.Song, e *utils.ErrorCode)
	Create(request requests.CreateSongRequest, token string) *utils.ErrorCode
	Update(request requests.UpdateSongRequest) *utils.ErrorCode
	Delete(id uuid.UUID) *utils.ErrorCode
}

type songService struct {
	repository repository.SongRepository
	jwtService service.JwtService
}

func NewSongService(repository repository.SongRepository, jwtService service.JwtService) SongService {
	return &songService{
		repository: repository,
		jwtService: jwtService,
	}
}

func (s *songService) Get(id uuid.UUID) (song models.Song, e *utils.ErrorCode) {
	err := s.repository.Get(&song, id)
	if err != nil {
		return song, utils.InternalServerError(err)
	}
	if song.ID == uuid.Nil {
		return song, utils.NotFoundError(errors.New("song not found"))
	}
	return song, nil
}

func (s *songService) GetAll(request requests.GetSongsRequest) (songs []models.Song, e *utils.ErrorCode) {
	err := s.repository.GetAllByUser(&songs, request.UserID)
	if err != nil {
		return songs, utils.InternalServerError(err)
	}
	return songs, nil
}

func (s *songService) Create(request requests.CreateSongRequest, token string) *utils.ErrorCode {
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

func (s *songService) Update(request requests.UpdateSongRequest) *utils.ErrorCode {
	var song models.Song
	err := s.repository.Get(&song, request.ID)
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

func (s *songService) Delete(id uuid.UUID) *utils.ErrorCode {
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
