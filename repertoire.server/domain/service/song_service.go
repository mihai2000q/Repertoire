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
	AddPerfectRehearsal(request requests.AddPerfectSongRehearsalRequest) *wrapper.ErrorCode
	AddPartialRehearsal(request requests.AddPartialSongRehearsalRequest) *wrapper.ErrorCode
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
	GetSectionTypes(token string) ([]model.SongSectionType, *wrapper.ErrorCode)

	CreateSection(request requests.CreateSongSectionRequest) *wrapper.ErrorCode
	DeleteSection(id uuid.UUID, songID uuid.UUID) *wrapper.ErrorCode
	MoveSection(request requests.MoveSongSectionRequest) *wrapper.ErrorCode
	UpdateAllSections(request requests.UpdateAllSongSectionsRequest) *wrapper.ErrorCode
	UpdateSection(request requests.UpdateSongSectionRequest) *wrapper.ErrorCode
	UpdateSectionsOccurrences(request requests.UpdateSongSectionsOccurrencesRequest) *wrapper.ErrorCode
	UpdateSectionsPartialOccurrences(request requests.UpdateSongSectionsPartialOccurrencesRequest) *wrapper.ErrorCode
}

type songService struct {
	addPerfectSongRehearsal song.AddPerfectSongRehearsal
	addPartialSongRehearsal song.AddPartialSongRehearsal
	createSong              song.CreateSong
	deleteImageFromSong     song.DeleteImageFromSong
	deleteSong              song.DeleteSong
	getAllSongs             song.GetAllSongs
	getSong                 song.GetSong
	getSongFiltersMetadata  song.GetSongFiltersMetadata
	saveImageToSong         song.SaveImageToSong
	updateSong              song.UpdateSong
	updateSongSettings      song.UpdateSongSettings

	getGuitarTunings    song.GetGuitarTunings
	getInstruments      song.GetInstruments
	getSongSectionTypes section.GetSongSectionTypes

	createSongSection                    section.CreateSongSection
	deleteSongSection                    section.DeleteSongSection
	moveSongSection                      section.MoveSongSection
	updateAllSongSections                section.UpdateAllSongSections
	updateSongSection                    section.UpdateSongSection
	updateSongSectionsOccurrences        section.UpdateSongSectionsOccurrences
	updateSongSectionsPartialOccurrences section.UpdateSongSectionsPartialOccurrences
}

func NewSongService(
	addPerfectSongRehearsal song.AddPerfectSongRehearsal,
	addPartialSongRehearsal song.AddPartialSongRehearsal,
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
	getSongSectionTypes section.GetSongSectionTypes,

	createSongSection section.CreateSongSection,
	deleteSongSection section.DeleteSongSection,
	moveSongSection section.MoveSongSection,
	updateAllSongSections section.UpdateAllSongSections,
	updateSongSection section.UpdateSongSection,
	updateSongSectionsOccurrences section.UpdateSongSectionsOccurrences,
	updateSongSectionsPartialOccurrences section.UpdateSongSectionsPartialOccurrences,
) SongService {
	return &songService{
		addPerfectSongRehearsal: addPerfectSongRehearsal,
		addPartialSongRehearsal: addPartialSongRehearsal,
		createSong:              createSong,
		deleteImageFromSong:     deleteImageFromSong,
		deleteSong:              deleteSong,
		getAllSongs:             getAllSongs,
		getSong:                 getSong,
		getSongFiltersMetadata:  getSongFiltersMetadata,
		saveImageToSong:         saveImageToSong,
		updateSong:              updateSong,
		updateSongSettings:      updateSongSettings,

		getGuitarTunings:    getGuitarTunings,
		getInstruments:      getInstruments,
		getSongSectionTypes: getSongSectionTypes,

		createSongSection:                    createSongSection,
		deleteSongSection:                    deleteSongSection,
		moveSongSection:                      moveSongSection,
		updateAllSongSections:                updateAllSongSections,
		updateSongSection:                    updateSongSection,
		updateSongSectionsOccurrences:        updateSongSectionsOccurrences,
		updateSongSectionsPartialOccurrences: updateSongSectionsPartialOccurrences,
	}
}

func (s *songService) AddPerfectRehearsal(request requests.AddPerfectSongRehearsalRequest) *wrapper.ErrorCode {
	return s.addPerfectSongRehearsal.Handle(request)
}

func (s *songService) AddPartialRehearsal(request requests.AddPartialSongRehearsalRequest) *wrapper.ErrorCode {
	return s.addPartialSongRehearsal.Handle(request)
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

func (s *songService) GetSectionTypes(token string) ([]model.SongSectionType, *wrapper.ErrorCode) {
	return s.getSongSectionTypes.Handle(token)
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

func (s *songService) UpdateAllSections(request requests.UpdateAllSongSectionsRequest) *wrapper.ErrorCode {
	return s.updateAllSongSections.Handle(request)
}

func (s *songService) UpdateSection(request requests.UpdateSongSectionRequest) *wrapper.ErrorCode {
	return s.updateSongSection.Handle(request)
}

func (s *songService) UpdateSectionsOccurrences(request requests.UpdateSongSectionsOccurrencesRequest) *wrapper.ErrorCode {
	return s.updateSongSectionsOccurrences.Handle(request)
}

func (s *songService) UpdateSectionsPartialOccurrences(request requests.UpdateSongSectionsPartialOccurrencesRequest) *wrapper.ErrorCode {
	return s.updateSongSectionsPartialOccurrences.Handle(request)
}
