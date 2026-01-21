package service

import (
	"mime/multipart"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type SongService interface {
	AddPerfectRehearsal(request requests.AddPerfectSongRehearsalRequest) *wrapper.ErrorCode
	AddPerfectRehearsals(request requests.AddPerfectSongRehearsalsRequest) *wrapper.ErrorCode
	AddPartialRehearsal(request requests.AddPartialSongRehearsalRequest) *wrapper.ErrorCode
	BulkDelete(request requests.BulkDeleteSongsRequest) *wrapper.ErrorCode
	Create(request requests.CreateSongRequest, token string) (uuid.UUID, *wrapper.ErrorCode)
	DeleteImage(id uuid.UUID) *wrapper.ErrorCode
	Delete(id uuid.UUID) *wrapper.ErrorCode
	GetAll(request requests.GetSongsRequest, token string) (wrapper.WithTotalCount[model.EnhancedSong], *wrapper.ErrorCode)
	Get(id uuid.UUID) (model.Song, *wrapper.ErrorCode)
	GetFiltersMetadata(
		request requests.GetSongFiltersMetadataRequest,
		token string,
	) (model.SongFiltersMetadata, *wrapper.ErrorCode)
	SaveImage(file *multipart.FileHeader, songID uuid.UUID) *wrapper.ErrorCode
	Update(request requests.UpdateSongRequest) *wrapper.ErrorCode
	UpdateSettings(request requests.UpdateSongSettingsRequest) *wrapper.ErrorCode

	GetGuitarTunings(token string) ([]model.GuitarTuning, *wrapper.ErrorCode)
	GetInstruments(token string) ([]model.Instrument, *wrapper.ErrorCode)
}

type songService struct {
	addPerfectSongRehearsal  song.AddPerfectSongRehearsal
	addPerfectSongRehearsals song.AddPerfectSongRehearsals
	addPartialSongRehearsal  song.AddPartialSongRehearsal
	bulkDeleteSongs          song.BulkDeleteSongs
	createSong               song.CreateSong
	deleteImageFromSong      song.DeleteImageFromSong
	deleteSong               song.DeleteSong
	getAllSongs              song.GetAllSongs
	getSong                  song.GetSong
	getSongFiltersMetadata   song.GetSongFiltersMetadata
	saveImageToSong          song.SaveImageToSong
	updateSong               song.UpdateSong
	updateSongSettings       song.UpdateSongSettings

	getGuitarTunings song.GetGuitarTunings
	getInstruments   song.GetInstruments
}

func NewSongService(
	addPerfectSongRehearsal song.AddPerfectSongRehearsal,
	addPerfectSongRehearsals song.AddPerfectSongRehearsals,
	addPartialSongRehearsal song.AddPartialSongRehearsal,
	bulkDeleteSongs song.BulkDeleteSongs,
	createSong song.CreateSong,
	deleteImageFromSong song.DeleteImageFromSong,
	deleteSong song.DeleteSong,
	getAllSongs song.GetAllSongs,
	getSong song.GetSong,
	getSongFiltersMetadata song.GetSongFiltersMetadata,
	saveImageToSong song.SaveImageToSong,
	updateSong song.UpdateSong,
	updateSongSettings song.UpdateSongSettings,

	getGuitarTunings song.GetGuitarTunings,
	getInstruments song.GetInstruments,
) SongService {
	return &songService{
		addPerfectSongRehearsal:  addPerfectSongRehearsal,
		addPerfectSongRehearsals: addPerfectSongRehearsals,
		addPartialSongRehearsal:  addPartialSongRehearsal,
		bulkDeleteSongs:          bulkDeleteSongs,
		createSong:               createSong,
		deleteImageFromSong:      deleteImageFromSong,
		deleteSong:               deleteSong,
		getAllSongs:              getAllSongs,
		getSong:                  getSong,
		getSongFiltersMetadata:   getSongFiltersMetadata,
		saveImageToSong:          saveImageToSong,
		updateSong:               updateSong,
		updateSongSettings:       updateSongSettings,

		getGuitarTunings: getGuitarTunings,
		getInstruments:   getInstruments,
	}
}

func (s *songService) AddPerfectRehearsal(request requests.AddPerfectSongRehearsalRequest) *wrapper.ErrorCode {
	return s.addPerfectSongRehearsal.Handle(request)
}

func (s *songService) AddPerfectRehearsals(request requests.AddPerfectSongRehearsalsRequest) *wrapper.ErrorCode {
	return s.addPerfectSongRehearsals.Handle(request)
}

func (s *songService) AddPartialRehearsal(request requests.AddPartialSongRehearsalRequest) *wrapper.ErrorCode {
	return s.addPartialSongRehearsal.Handle(request)
}

func (s *songService) BulkDelete(request requests.BulkDeleteSongsRequest) *wrapper.ErrorCode {
	return s.bulkDeleteSongs.Handle(request)
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

func (s *songService) GetAll(request requests.GetSongsRequest, token string) (wrapper.WithTotalCount[model.EnhancedSong], *wrapper.ErrorCode) {
	return s.getAllSongs.Handle(request, token)
}

func (s *songService) Get(id uuid.UUID) (model.Song, *wrapper.ErrorCode) {
	return s.getSong.Handle(id)
}

func (s *songService) GetFiltersMetadata(
	request requests.GetSongFiltersMetadataRequest,
	token string,
) (model.SongFiltersMetadata, *wrapper.ErrorCode) {
	return s.getSongFiltersMetadata.Handle(request, token)
}

func (s *songService) SaveImage(file *multipart.FileHeader, songID uuid.UUID) *wrapper.ErrorCode {
	return s.saveImageToSong.Handle(file, songID)
}

func (s *songService) Update(request requests.UpdateSongRequest) *wrapper.ErrorCode {
	return s.updateSong.Handle(request)
}

func (s *songService) UpdateSettings(request requests.UpdateSongSettingsRequest) *wrapper.ErrorCode {
	return s.updateSongSettings.Handle(request)
}

func (s *songService) GetGuitarTunings(token string) ([]model.GuitarTuning, *wrapper.ErrorCode) {
	return s.getGuitarTunings.Handle(token)
}

func (s *songService) GetInstruments(token string) ([]model.Instrument, *wrapper.ErrorCode) {
	return s.getInstruments.Handle(token)
}
