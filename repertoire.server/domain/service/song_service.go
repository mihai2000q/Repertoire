package service

import (
	"mime/multipart"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/pagination"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type SongService interface {
	AddCustomRehearsal(request requests.AddCustomSongRehearsalRequest) *httperror.ErrorCode
	AddCustomRehearsals(request requests.AddCustomSongRehearsalsRequest) *httperror.ErrorCode
	AddPerfectRehearsal(request requests.AddPerfectSongRehearsalRequest) *httperror.ErrorCode
	AddPerfectRehearsals(request requests.AddPerfectSongRehearsalsRequest) *httperror.ErrorCode
	BulkDelete(request requests.BulkDeleteSongsRequest) *httperror.ErrorCode
	Create(request requests.CreateSongRequest, token string) (uuid.UUID, *httperror.ErrorCode)
	DeleteImage(id uuid.UUID) *httperror.ErrorCode
	Delete(id uuid.UUID) *httperror.ErrorCode
	GetAll(request requests.GetSongsRequest, token string) (pagination.WithTotalCount[model.EnhancedSong], *httperror.ErrorCode)
	Get(id uuid.UUID) (model.Song, *httperror.ErrorCode)
	GetFiltersMetadata(
		request requests.GetSongFiltersMetadataRequest,
		token string,
	) (model.SongFiltersMetadata, *httperror.ErrorCode)
	SaveImage(file *multipart.FileHeader, songID uuid.UUID) *httperror.ErrorCode
	Update(request requests.UpdateSongRequest) *httperror.ErrorCode
	UpdateSettings(request requests.UpdateSongSettingsRequest) *httperror.ErrorCode

	GetGuitarTunings(token string) ([]model.GuitarTuning, *httperror.ErrorCode)
	GetInstruments(token string) ([]model.Instrument, *httperror.ErrorCode)
}

type songService struct {
	addCustomSongRehearsal   song.AddCustomSongRehearsal
	addCustomSongRehearsals  song.AddCustomSongRehearsals
	addPerfectSongRehearsal  song.AddPerfectSongRehearsal
	addPerfectSongRehearsals song.AddPerfectSongRehearsals
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
	addCustomSongRehearsal song.AddCustomSongRehearsal,
	addCustomSongRehearsals song.AddCustomSongRehearsals,
	addPerfectSongRehearsal song.AddPerfectSongRehearsal,
	addPerfectSongRehearsals song.AddPerfectSongRehearsals,
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
		addCustomSongRehearsal:   addCustomSongRehearsal,
		addCustomSongRehearsals:  addCustomSongRehearsals,
		addPerfectSongRehearsal:  addPerfectSongRehearsal,
		addPerfectSongRehearsals: addPerfectSongRehearsals,
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

func (s *songService) AddCustomRehearsal(request requests.AddCustomSongRehearsalRequest) *httperror.ErrorCode {
	return s.addCustomSongRehearsal.Handle(request)
}

func (s *songService) AddCustomRehearsals(request requests.AddCustomSongRehearsalsRequest) *httperror.ErrorCode {
	return s.addCustomSongRehearsals.Handle(request)
}

func (s *songService) AddPerfectRehearsal(request requests.AddPerfectSongRehearsalRequest) *httperror.ErrorCode {
	return s.addPerfectSongRehearsal.Handle(request)
}

func (s *songService) AddPerfectRehearsals(request requests.AddPerfectSongRehearsalsRequest) *httperror.ErrorCode {
	return s.addPerfectSongRehearsals.Handle(request)
}

func (s *songService) BulkDelete(request requests.BulkDeleteSongsRequest) *httperror.ErrorCode {
	return s.bulkDeleteSongs.Handle(request)
}

func (s *songService) Create(request requests.CreateSongRequest, token string) (uuid.UUID, *httperror.ErrorCode) {
	return s.createSong.Handle(request, token)
}

func (s *songService) Delete(id uuid.UUID) *httperror.ErrorCode {
	return s.deleteSong.Handle(id)
}

func (s *songService) DeleteImage(id uuid.UUID) *httperror.ErrorCode {
	return s.deleteImageFromSong.Handle(id)
}

func (s *songService) GetAll(request requests.GetSongsRequest, token string) (pagination.WithTotalCount[model.EnhancedSong], *httperror.ErrorCode) {
	return s.getAllSongs.Handle(request, token)
}

func (s *songService) Get(id uuid.UUID) (model.Song, *httperror.ErrorCode) {
	return s.getSong.Handle(id)
}

func (s *songService) GetFiltersMetadata(
	request requests.GetSongFiltersMetadataRequest,
	token string,
) (model.SongFiltersMetadata, *httperror.ErrorCode) {
	return s.getSongFiltersMetadata.Handle(request, token)
}

func (s *songService) SaveImage(file *multipart.FileHeader, songID uuid.UUID) *httperror.ErrorCode {
	return s.saveImageToSong.Handle(file, songID)
}

func (s *songService) Update(request requests.UpdateSongRequest) *httperror.ErrorCode {
	return s.updateSong.Handle(request)
}

func (s *songService) UpdateSettings(request requests.UpdateSongSettingsRequest) *httperror.ErrorCode {
	return s.updateSongSettings.Handle(request)
}

func (s *songService) GetGuitarTunings(token string) ([]model.GuitarTuning, *httperror.ErrorCode) {
	return s.getGuitarTunings.Handle(token)
}

func (s *songService) GetInstruments(token string) ([]model.Instrument, *httperror.ErrorCode) {
	return s.getInstruments.Handle(token)
}
