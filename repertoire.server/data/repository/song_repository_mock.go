package repository

import (
	"repertoire/model"

	"github.com/stretchr/testify/mock"

	"github.com/google/uuid"
)

type SongRepositoryMock struct {
	mock.Mock
}

func (s *SongRepositoryMock) Get(song *model.Song, id uuid.UUID) error {
	args := s.Called(song, id)

	if len(args) > 1 {
		*song = *args.Get(1).(*model.Song)
	}

	return args.Error(0)
}

func (s *SongRepositoryMock) GetWithSections(song *model.Song, id uuid.UUID) error {
	args := s.Called(song, id)

	if len(args) > 1 {
		*song = *args.Get(1).(*model.Song)
	}

	return args.Error(0)
}

func (s *SongRepositoryMock) GetWithAssociations(song *model.Song, id uuid.UUID) error {
	args := s.Called(song, id)

	if len(args) > 1 {
		*song = *args.Get(1).(*model.Song)
	}

	return args.Error(0)
}

func (s *SongRepositoryMock) GetAllByUser(
	songs *[]model.Song,
	userID uuid.UUID,
	currentPage *int,
	pageSize *int,
	orderBy string,
) error {
	args := s.Called(songs, userID, currentPage, pageSize, orderBy)

	if len(args) > 1 {
		*songs = *args.Get(1).(*[]model.Song)
	}

	return args.Error(0)
}

func (s *SongRepositoryMock) GetAllByUserCount(count *int64, userID uuid.UUID) error {
	args := s.Called(count, userID)

	if len(args) > 1 {
		*count = *args.Get(1).(*int64)
	}

	return args.Error(0)
}

func (s *SongRepositoryMock) GetGuitarTunings(tunings *[]model.GuitarTuning, userID uuid.UUID) error {
	args := s.Called(tunings, userID)

	if len(args) > 1 {
		*tunings = *args.Get(1).(*[]model.GuitarTuning)
	}

	return args.Error(0)
}

func (s *SongRepositoryMock) Create(song *model.Song) error {
	args := s.Called(song)
	return args.Error(0)
}

func (s *SongRepositoryMock) Update(song *model.Song) error {
	args := s.Called(song)
	return args.Error(0)
}

func (s *SongRepositoryMock) Delete(id uuid.UUID) error {
	args := s.Called(id)
	return args.Error(0)
}

// Sections

func (s *SongRepositoryMock) GetSection(section *model.SongSection, id uuid.UUID) error {
	args := s.Called(section, id)

	if len(args) > 1 {
		*section = *args.Get(1).(*model.SongSection)
	}

	return args.Error(0)
}

func (s *SongRepositoryMock) GetSectionTypes(tunings *[]model.SongSectionType, userID uuid.UUID) error {
	args := s.Called(tunings, userID)

	if len(args) > 1 {
		*tunings = *args.Get(1).(*[]model.SongSectionType)
	}

	return args.Error(0)
}

func (s *SongRepositoryMock) CountSectionsBySong(count *int64, songID uuid.UUID) error {
	args := s.Called(count, songID)

	if len(args) > 1 {
		*count = *args.Get(1).(*int64)
	}

	return args.Error(0)
}

func (s *SongRepositoryMock) CreateSection(section *model.SongSection) error {
	args := s.Called(section)
	return args.Error(0)
}

func (s *SongRepositoryMock) UpdateSection(section *model.SongSection) error {
	args := s.Called(section)
	return args.Error(0)
}

func (s *SongRepositoryMock) DeleteSection(id uuid.UUID) error {
	args := s.Called(id)
	return args.Error(0)
}
