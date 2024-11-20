package service

import (
	"mime/multipart"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song"
	"repertoire/server/domain/usecase/song/guitar/tuning"
	"repertoire/server/domain/usecase/song/section"
	"repertoire/server/domain/usecase/song/section/types"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type SongService interface {
	Create(request requests.CreateSongRequest, token string) (uuid.UUID, *wrapper.ErrorCode)
	DeleteImage(id uuid.UUID) *wrapper.ErrorCode
	Delete(id uuid.UUID) *wrapper.ErrorCode
	GetAll(request requests.GetSongsRequest, token string) (wrapper.WithTotalCount[model.Song], *wrapper.ErrorCode)
	Get(id uuid.UUID) (model.Song, *wrapper.ErrorCode)
	SaveImage(file *multipart.FileHeader, songID uuid.UUID) *wrapper.ErrorCode
	Update(request requests.UpdateSongRequest) *wrapper.ErrorCode

	CreateGuitarTuning(request requests.CreateGuitarTuningRequest, token string) *wrapper.ErrorCode
	MoveGuitarTuning(request requests.MoveGuitarTuningRequest, token string) *wrapper.ErrorCode
	DeleteGuitarTuning(id uuid.UUID, token string) *wrapper.ErrorCode
	GetGuitarTunings(token string) ([]model.GuitarTuning, *wrapper.ErrorCode)

	CreateSection(request requests.CreateSongSectionRequest) *wrapper.ErrorCode
	DeleteSection(id uuid.UUID, songID uuid.UUID) *wrapper.ErrorCode
	MoveSection(request requests.MoveSongSectionRequest) *wrapper.ErrorCode
	UpdateSection(request requests.UpdateSongSectionRequest) *wrapper.ErrorCode

	CreateSectionType(request requests.CreateSongSectionTypeRequest, token string) *wrapper.ErrorCode
	GetSectionTypes(token string) ([]model.SongSectionType, *wrapper.ErrorCode)
}

type songService struct {
	createSong          song.CreateSong
	deleteImageFromSong song.DeleteImageFromSong
	deleteSong          song.DeleteSong
	getAllSongs         song.GetAllSongs
	getSong             song.GetSong
	saveImageToSong     song.SaveImageToSong
	updateSong          song.UpdateSong

	createGuitarTuning tuning.CreateGuitarTuning
	deleteGuitarTuning tuning.DeleteGuitarTuning
	getGuitarTunings   tuning.GetGuitarTunings
	moveGuitarTuning   tuning.MoveGuitarTuning

	createSongSection section.CreateSongSection
	deleteSongSection section.DeleteSongSection
	moveSongSection   section.MoveSongSection
	updateSongSection section.UpdateSongSection

	createSongSectionType types.CreateSongSectionType
	getSongSectionTypes   types.GetSongSectionTypes
}

func NewSongService(
	createSong song.CreateSong,
	deleteImageFromSong song.DeleteImageFromSong,
	deleteSong song.DeleteSong,
	getAllSongs song.GetAllSongs,
	getSong song.GetSong,
	saveImageToSong song.SaveImageToSong,
	updateSong song.UpdateSong,

	createGuitarTuning tuning.CreateGuitarTuning,
	deleteGuitarTuning tuning.DeleteGuitarTuning,
	getGuitarTunings tuning.GetGuitarTunings,
	moveGuitarTuning tuning.MoveGuitarTuning,

	createSongSection section.CreateSongSection,
	deleteSongSection section.DeleteSongSection,
	moveSongSection section.MoveSongSection,
	updateSongSection section.UpdateSongSection,

	createSongSectionType types.CreateSongSectionType,
	getSongSectionTypes types.GetSongSectionTypes,
) SongService {
	return &songService{
		createSong:          createSong,
		deleteImageFromSong: deleteImageFromSong,
		deleteSong:          deleteSong,
		getAllSongs:         getAllSongs,
		getSong:             getSong,
		saveImageToSong:     saveImageToSong,
		updateSong:          updateSong,

		createGuitarTuning: createGuitarTuning,
		deleteGuitarTuning: deleteGuitarTuning,
		getGuitarTunings:   getGuitarTunings,
		moveGuitarTuning:   moveGuitarTuning,

		createSongSection: createSongSection,
		deleteSongSection: deleteSongSection,
		moveSongSection:   moveSongSection,
		updateSongSection: updateSongSection,

		createSongSectionType: createSongSectionType,
		getSongSectionTypes:   getSongSectionTypes,
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

func (s *songService) Get(id uuid.UUID) (model.Song, *wrapper.ErrorCode) {
	return s.getSong.Handle(id)
}

func (s *songService) SaveImage(file *multipart.FileHeader, songID uuid.UUID) *wrapper.ErrorCode {
	return s.saveImageToSong.Handle(file, songID)
}

func (s *songService) Update(request requests.UpdateSongRequest) *wrapper.ErrorCode {
	return s.updateSong.Handle(request)
}

// Guitar Tunings

func (s *songService) CreateGuitarTuning(request requests.CreateGuitarTuningRequest, token string) *wrapper.ErrorCode {
	return s.createGuitarTuning.Handle(request, token)
}

func (s *songService) DeleteGuitarTuning(id uuid.UUID, token string) *wrapper.ErrorCode {
	return s.deleteGuitarTuning.Handle(id, token)
}

func (s *songService) GetGuitarTunings(token string) ([]model.GuitarTuning, *wrapper.ErrorCode) {
	return s.getGuitarTunings.Handle(token)
}

func (s *songService) MoveGuitarTuning(request requests.MoveGuitarTuningRequest, token string) *wrapper.ErrorCode {
	return s.moveGuitarTuning.Handle(request, token)
}

// Sections

func (s *songService) CreateSection(request requests.CreateSongSectionRequest) *wrapper.ErrorCode {
	return s.createSongSection.Handle(request)
}

func (s *songService) DeleteSection(id uuid.UUID, songID uuid.UUID) *wrapper.ErrorCode {
	return s.deleteSongSection.Handle(id, songID)
}

func (s *songService) MoveSection(request requests.MoveSongSectionRequest) *wrapper.ErrorCode {
	return s.moveSongSection.Handle(request)
}

func (s *songService) UpdateSection(request requests.UpdateSongSectionRequest) *wrapper.ErrorCode {
	return s.updateSongSection.Handle(request)
}

// Section - Types

func (s *songService) CreateSectionType(
	request requests.CreateSongSectionTypeRequest,
	token string,
) *wrapper.ErrorCode {
	return s.createSongSectionType.Handle(request, token)
}

func (s *songService) GetSectionTypes(token string) ([]model.SongSectionType, *wrapper.ErrorCode) {
	return s.getSongSectionTypes.Handle(token)
}
