package service

import (
	"github.com/google/uuid"
	"repertoire/api/requests"
	"repertoire/domain/usecases/song"
	"repertoire/models"
	"repertoire/utils/wrapper"
)

type SongService interface {
	Get(id uuid.UUID) (song models.Song, e *wrapper.ErrorCode)
	GetAll(request requests.GetSongsRequest) (songs []models.Song, e *wrapper.ErrorCode)
	Create(request requests.CreateSongRequest, token string) *wrapper.ErrorCode
	Update(request requests.UpdateSongRequest) *wrapper.ErrorCode
	Delete(id uuid.UUID) *wrapper.ErrorCode
}

type songService struct {
	getSong     song.GetSong
	getAllSongs song.GetAllSongs
	createSong  song.CreateSong
	updateSong  song.UpdateSong
	deleteSong  song.DeleteSong
}

func NewSongService(
	getSong song.GetSong,
	getAllSongs song.GetAllSongs,
	createSong song.CreateSong,
	updateSong song.UpdateSong,
	deleteSong song.DeleteSong,
) SongService {
	return &songService{
		getSong:     getSong,
		getAllSongs: getAllSongs,
		createSong:  createSong,
		updateSong:  updateSong,
		deleteSong:  deleteSong,
	}
}

func (s *songService) Get(id uuid.UUID) (models.Song, *wrapper.ErrorCode) {
	return s.getSong.Handle(id)
}

func (s *songService) GetAll(request requests.GetSongsRequest) (songs []models.Song, e *wrapper.ErrorCode) {
	return s.getAllSongs.Handle(request)
}

func (s *songService) Create(request requests.CreateSongRequest, token string) *wrapper.ErrorCode {
	return s.createSong.Handle(request, token)
}

func (s *songService) Update(request requests.UpdateSongRequest) *wrapper.ErrorCode {
	return s.updateSong.Handle(request)
}

func (s *songService) Delete(id uuid.UUID) *wrapper.ErrorCode {
	return s.deleteSong.Handle(id)
}
