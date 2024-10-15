package service

import (
	"github.com/google/uuid"
	"repertoire/api/requests"
	"repertoire/domain/usecases/song"
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

func (s *songService) Get(id uuid.UUID) (models.Song, *utils.ErrorCode) {
	return s.getSong.Handle(id)
}

func (s *songService) GetAll(request requests.GetSongsRequest) (songs []models.Song, e *utils.ErrorCode) {
	return s.getAllSongs.Handle(request)
}

func (s *songService) Create(request requests.CreateSongRequest, token string) *utils.ErrorCode {
	return s.createSong.Handle(request, token)
}

func (s *songService) Update(request requests.UpdateSongRequest) *utils.ErrorCode {
	return s.updateSong.Handle(request)
}

func (s *songService) Delete(id uuid.UUID) *utils.ErrorCode {
	return s.deleteSong.Handle(id)
}
