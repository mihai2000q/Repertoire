package service

import (
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song"
	"repertoire/server/domain/usecase/song/section"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type SongService interface {
	Get(id uuid.UUID) (model.Song, *wrapper.ErrorCode)
	GetAll(request requests.GetSongsRequest, token string) (wrapper.WithTotalCount[model.Song], *wrapper.ErrorCode)
	GetGuitarTunings(token string) ([]model.GuitarTuning, *wrapper.ErrorCode)
	Create(request requests.CreateSongRequest, token string) *wrapper.ErrorCode
	Update(request requests.UpdateSongRequest) *wrapper.ErrorCode
	Delete(id uuid.UUID) *wrapper.ErrorCode

	GetSectionTypes(token string) ([]model.SongSectionType, *wrapper.ErrorCode)
	CreateSection(request requests.CreateSongSectionRequest) *wrapper.ErrorCode
	MoveSection(request requests.MoveSongSectionRequest) *wrapper.ErrorCode
	UpdateSection(request requests.UpdateSongSectionRequest) *wrapper.ErrorCode
	DeleteSection(id uuid.UUID, songID uuid.UUID) *wrapper.ErrorCode
}

type songService struct {
	getSong             song.GetSong
	getAllSongs         song.GetAllSongs
	getGuitarTunings    song.GetGuitarTunings
	createSong          song.CreateSong
	updateSong          song.UpdateSong
	deleteSong          song.DeleteSong
	getSongSectionTypes section.GetSongSectionTypes
	createSongSection   section.CreateSongSection
	moveSongSection     section.MoveSongSection
	updateSongSection   section.UpdateSongSection
	deleteSongSection   section.DeleteSongSection
}

func NewSongService(
	getSong song.GetSong,
	getAllSongs song.GetAllSongs,
	getGuitarTunings song.GetGuitarTunings,
	createSong song.CreateSong,
	updateSong song.UpdateSong,
	deleteSong song.DeleteSong,
	getSongSectionTypes section.GetSongSectionTypes,
	createSongSection section.CreateSongSection,
	moveSongSection section.MoveSongSection,
	updateSongSection section.UpdateSongSection,
	deleteSongSection section.DeleteSongSection,
) SongService {
	return &songService{
		getSong:             getSong,
		getAllSongs:         getAllSongs,
		getGuitarTunings:    getGuitarTunings,
		createSong:          createSong,
		updateSong:          updateSong,
		deleteSong:          deleteSong,
		getSongSectionTypes: getSongSectionTypes,
		createSongSection:   createSongSection,
		moveSongSection:     moveSongSection,
		updateSongSection:   updateSongSection,
		deleteSongSection:   deleteSongSection,
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

func (s *songService) Create(request requests.CreateSongRequest, token string) *wrapper.ErrorCode {
	return s.createSong.Handle(request, token)
}

func (s *songService) Update(request requests.UpdateSongRequest) *wrapper.ErrorCode {
	return s.updateSong.Handle(request)
}

func (s *songService) Delete(id uuid.UUID) *wrapper.ErrorCode {
	return s.deleteSong.Handle(id)
}

// Sections

func (s *songService) GetSectionTypes(token string) ([]model.SongSectionType, *wrapper.ErrorCode) {
	return s.getSongSectionTypes.Handle(token)
}

func (s *songService) CreateSection(request requests.CreateSongSectionRequest) *wrapper.ErrorCode {
	return s.createSongSection.Handle(request)
}

func (s *songService) MoveSection(request requests.MoveSongSectionRequest) *wrapper.ErrorCode {
	return s.moveSongSection.Handle(request)
}

func (s *songService) UpdateSection(request requests.UpdateSongSectionRequest) *wrapper.ErrorCode {
	return s.updateSongSection.Handle(request)
}

func (s *songService) DeleteSection(id uuid.UUID, songID uuid.UUID) *wrapper.ErrorCode {
	return s.deleteSongSection.Handle(id, songID)
}
