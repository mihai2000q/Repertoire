package service

import (
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song/section"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type SongSectionService interface {
	BulkRehearsals(request requests.BulkRehearsalsSongSectionsRequest) *wrapper.ErrorCode
	BulkDelete(request requests.BulkDeleteSongSectionsRequest) *wrapper.ErrorCode
	Create(request requests.CreateSongSectionRequest) *wrapper.ErrorCode
	Delete(id uuid.UUID, songID uuid.UUID) *wrapper.ErrorCode
	Move(request requests.MoveSongSectionRequest) *wrapper.ErrorCode
	UpdateAll(request requests.UpdateAllSongSectionsRequest) *wrapper.ErrorCode
	Update(request requests.UpdateSongSectionRequest) *wrapper.ErrorCode

	GetTypes(token string) ([]model.SongSectionType, *wrapper.ErrorCode)
}

type songSectionService struct {
	bulkRehearsalsSongSections           section.BulkRehearsalsSongSections
	bulkDeleteSongSections               section.BulkDeleteSongSections
	createSongSection                    section.CreateSongSection
	deleteSongSection                    section.DeleteSongSection
	moveSongSection                      section.MoveSongSection
	updateAllSongSections                section.UpdateAllSongSections
	updateSongSection                    section.UpdateSongSection
	getSongSectionTypes                  section.GetSongSectionTypes
}

func NewSongSectionService(
	bulkRehearsalsSongSections section.BulkRehearsalsSongSections,
	bulkDeleteSongSections section.BulkDeleteSongSections,
	createSongSection section.CreateSongSection,
	deleteSongSection section.DeleteSongSection,
	moveSongSection section.MoveSongSection,
	updateAllSongSections section.UpdateAllSongSections,
	updateSongSection section.UpdateSongSection,

	getSongSectionTypes section.GetSongSectionTypes,
) SongSectionService {
	return &songSectionService{
		bulkRehearsalsSongSections:           bulkRehearsalsSongSections,
		bulkDeleteSongSections:               bulkDeleteSongSections,
		createSongSection:                    createSongSection,
		deleteSongSection:                    deleteSongSection,
		moveSongSection:                      moveSongSection,
		updateAllSongSections:                updateAllSongSections,
		updateSongSection:                    updateSongSection,

		getSongSectionTypes: getSongSectionTypes,
	}
}

func (s *songSectionService) BulkRehearsals(request requests.BulkRehearsalsSongSectionsRequest) *wrapper.ErrorCode {
	return s.bulkRehearsalsSongSections.Handle(request)
}

func (s *songSectionService) BulkDelete(request requests.BulkDeleteSongSectionsRequest) *wrapper.ErrorCode {
	return s.bulkDeleteSongSections.Handle(request)
}

func (s *songSectionService) Create(request requests.CreateSongSectionRequest) *wrapper.ErrorCode {
	return s.createSongSection.Handle(request)
}

func (s *songSectionService) Delete(id uuid.UUID, songID uuid.UUID) *wrapper.ErrorCode {
	return s.deleteSongSection.Handle(id, songID)
}

func (s *songSectionService) Move(request requests.MoveSongSectionRequest) *wrapper.ErrorCode {
	return s.moveSongSection.Handle(request)
}

func (s *songSectionService) UpdateAll(request requests.UpdateAllSongSectionsRequest) *wrapper.ErrorCode {
	return s.updateAllSongSections.Handle(request)
}

func (s *songSectionService) Update(request requests.UpdateSongSectionRequest) *wrapper.ErrorCode {
	return s.updateSongSection.Handle(request)
}

// Types

func (s *songSectionService) GetTypes(token string) ([]model.SongSectionType, *wrapper.ErrorCode) {
	return s.getSongSectionTypes.Handle(token)
}
