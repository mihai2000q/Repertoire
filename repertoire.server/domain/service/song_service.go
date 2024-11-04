package service

import (
	"mime/multipart"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song"
	"repertoire/server/domain/usecase/song/section"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type SongService interface {
	Create(request requests.CreateSongRequest, token string) (uuid.UUID, *wrapper.ErrorCode)
	DeleteImage(id uuid.UUID) *wrapper.ErrorCode
	Delete(id uuid.UUID) *wrapper.ErrorCode
	GetAll(request requests.GetSongsRequest, token string) (wrapper.WithTotalCount[model.Song], *wrapper.ErrorCode)
	GetGuitarTunings(token string) ([]model.GuitarTuning, *wrapper.ErrorCode)
	Get(id uuid.UUID) (model.Song, *wrapper.ErrorCode)
	SaveImage(file *multipart.FileHeader, songID uuid.UUID, token string) *wrapper.ErrorCode
	Update(request requests.UpdateSongRequest) *wrapper.ErrorCode

	CreateSection(request requests.CreateSongSectionRequest) *wrapper.ErrorCode
	DeleteSection(id uuid.UUID, songID uuid.UUID) *wrapper.ErrorCode
	GetSectionTypes(token string) ([]model.SongSectionType, *wrapper.ErrorCode)
	MoveSection(request requests.MoveSongSectionRequest) *wrapper.ErrorCode
	UpdateSection(request requests.UpdateSongSectionRequest) *wrapper.ErrorCode
}

type songService struct {
	createSong          song.CreateSong
	deleteImageFromSong song.DeleteImageFromSong
	deleteSong          song.DeleteSong
	getAllSongs         song.GetAllSongs
	getGuitarTunings    song.GetGuitarTunings
	getSong             song.GetSong
	saveImageToSong     song.SaveImageToSong
	updateSong          song.UpdateSong

	createSongSection   section.CreateSongSection
	deleteSongSection   section.DeleteSongSection
	getSongSectionTypes section.GetSongSectionTypes
	moveSongSection     section.MoveSongSection
	updateSongSection   section.UpdateSongSection
}

func NewSongService(
	createSong song.CreateSong,
	deleteImageFromSong song.DeleteImageFromSong,
	deleteSong song.DeleteSong,
	getAllSongs song.GetAllSongs,
	getGuitarTunings song.GetGuitarTunings,
	getSong song.GetSong,
	saveImageToSong song.SaveImageToSong,
	updateSong song.UpdateSong,

	createSongSection section.CreateSongSection,
	deleteSongSection section.DeleteSongSection,
	getSongSectionTypes section.GetSongSectionTypes,
	moveSongSection section.MoveSongSection,
	updateSongSection section.UpdateSongSection,
) SongService {
	return &songService{
		createSong:          createSong,
		deleteImageFromSong: deleteImageFromSong,
		deleteSong:          deleteSong,
		getAllSongs:         getAllSongs,
		getGuitarTunings:    getGuitarTunings,
		getSong:             getSong,
		saveImageToSong:     saveImageToSong,
		updateSong:          updateSong,

		createSongSection:   createSongSection,
		deleteSongSection:   deleteSongSection,
		getSongSectionTypes: getSongSectionTypes,
		moveSongSection:     moveSongSection,
		updateSongSection:   updateSongSection,
	}
}

func (s *songService) Create(request requests.CreateSongRequest, token string) (uuid.UUID, *wrapper.ErrorCode) {
	return s.createSong.Handle(request, token)
}

func (s *songService) Delete(id uuid.UUID) *wrapper.ErrorCode {
	return s.deleteSong.Handle(id)
}

func (s *songService) DeleteImage(id uuid.UUID) *wrapper.ErrorCode {
	return s.deleteImageFromSong.Handle(id)
}

func (s *songService) GetAll(request requests.GetSongsRequest, token string) (wrapper.WithTotalCount[model.Song], *wrapper.ErrorCode) {
	return s.getAllSongs.Handle(request, token)
}

func (s *songService) GetGuitarTunings(token string) ([]model.GuitarTuning, *wrapper.ErrorCode) {
	return s.getGuitarTunings.Handle(token)
}

func (s *songService) Get(id uuid.UUID) (model.Song, *wrapper.ErrorCode) {
	return s.getSong.Handle(id)
}

func (s *songService) SaveImage(file *multipart.FileHeader, songID uuid.UUID, token string) *wrapper.ErrorCode {
	return s.saveImageToSong.Handle(file, songID, token)
}

func (s *songService) Update(request requests.UpdateSongRequest) *wrapper.ErrorCode {
	return s.updateSong.Handle(request)
}

// Sections

func (s *songService) CreateSection(request requests.CreateSongSectionRequest) *wrapper.ErrorCode {
	return s.createSongSection.Handle(request)
}

func (s *songService) DeleteSection(id uuid.UUID, songID uuid.UUID) *wrapper.ErrorCode {
	return s.deleteSongSection.Handle(id, songID)
}

func (s *songService) GetSectionTypes(token string) ([]model.SongSectionType, *wrapper.ErrorCode) {
	return s.getSongSectionTypes.Handle(token)
}

func (s *songService) MoveSection(request requests.MoveSongSectionRequest) *wrapper.ErrorCode {
	return s.moveSongSection.Handle(request)
}

func (s *songService) UpdateSection(request requests.UpdateSongSectionRequest) *wrapper.ErrorCode {
	return s.updateSongSection.Handle(request)
}
