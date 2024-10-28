package service

import (
	"github.com/google/uuid"
	"repertoire/api/requests"
	"repertoire/domain/usecase/song"
	"repertoire/domain/usecase/song/section"
	"repertoire/model"
	"repertoire/utils/wrapper"
)

type SongService interface {
	Get(id uuid.UUID) (model.Song, *wrapper.ErrorCode)
	GetAll(request requests.GetSongsRequest, token string) (wrapper.WithTotalCount[model.Song], *wrapper.ErrorCode)
	GetGuitarTunings(token string) ([]model.GuitarTuning, *wrapper.ErrorCode)
	GetSectionTypes(token string) ([]model.SongSectionType, *wrapper.ErrorCode)
	Create(request requests.CreateSongRequest, token string) *wrapper.ErrorCode
	Update(request requests.UpdateSongRequest) *wrapper.ErrorCode
	Delete(id uuid.UUID) *wrapper.ErrorCode
}

type songService struct {
	getSong             song.GetSong
	getAllSongs         song.GetAllSongs
	getGuitarTunings    song.GetGuitarTunings
	getSongSectionTypes section.GetSongSectionTypes
	createSong          song.CreateSong
	updateSong          song.UpdateSong
	deleteSong          song.DeleteSong
}

func NewSongService(
	getSong song.GetSong,
	getAllSongs song.GetAllSongs,
	getGuitarTunings song.GetGuitarTunings,
	getSongSectionTypes section.GetSongSectionTypes,
	createSong song.CreateSong,
	updateSong song.UpdateSong,
	deleteSong song.DeleteSong,
) SongService {
	return &songService{
		getSong:             getSong,
		getAllSongs:         getAllSongs,
		getGuitarTunings:    getGuitarTunings,
		getSongSectionTypes: getSongSectionTypes,
		createSong:          createSong,
		updateSong:          updateSong,
		deleteSong:          deleteSong,
	}
}

func (s *songService) Get(id uuid.UUID) (model.Song, *wrapper.ErrorCode) {
	return s.getSong.Handle(id)
}

func (s *songService) GetAll(request requests.GetSongsRequest, token string) (wrapper.WithTotalCount[model.Song], *wrapper.ErrorCode) {
	return s.getAllSongs.Handle(request, token)
}

func (s *songService) GetGuitarTunings(token string) ([]model.GuitarTuning, *wrapper.ErrorCode) {
	return s.getGuitarTunings.Handle(token)
}

func (s *songService) GetSectionTypes(token string) ([]model.SongSectionType, *wrapper.ErrorCode) {
	return s.getSongSectionTypes.Handle(token)
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
