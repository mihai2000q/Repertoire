package service

import (
	"github.com/google/uuid"
	"repertoire/api/request"
	"repertoire/domain/usecase/song"
	"repertoire/model"
	"repertoire/utils/wrapper"
)

type SongService interface {
	Get(id uuid.UUID) (model.Song, *wrapper.ErrorCode)
	GetAll(request request.GetSongsRequest, token string) (wrapper.WithTotalCount[model.Song], *wrapper.ErrorCode)
	Create(request request.CreateSongRequest, token string) *wrapper.ErrorCode
	Update(request request.UpdateSongRequest) *wrapper.ErrorCode
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

func (s *songService) Get(id uuid.UUID) (model.Song, *wrapper.ErrorCode) {
	return s.getSong.Handle(id)
}

func (s *songService) GetAll(request request.GetSongsRequest, token string) (wrapper.WithTotalCount[model.Song], *wrapper.ErrorCode) {
	return s.getAllSongs.Handle(request, token)
}

func (s *songService) Create(request request.CreateSongRequest, token string) *wrapper.ErrorCode {
	return s.createSong.Handle(request, token)
}

func (s *songService) Update(request request.UpdateSongRequest) *wrapper.ErrorCode {
	return s.updateSong.Handle(request)
}

func (s *songService) Delete(id uuid.UUID) *wrapper.ErrorCode {
	return s.deleteSong.Handle(id)
}
