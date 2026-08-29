package repository

import (
	"repertoire/server/model"

	"github.com/stretchr/testify/mock"

	"github.com/google/uuid"
)

type SongSectionRepositoryMock struct {
	mock.Mock
}

func (s *SongSectionRepositoryMock) Get(section *model.SongSection, id uuid.UUID) error {
	args := s.Called(section, id)

	if len(args) > 1 {
		*section = *args.Get(1).(*model.SongSection)
	}

	return args.Error(0)
}

func (s *SongSectionRepositoryMock) GetWithParts(section *model.SongSection, id uuid.UUID) error {
	args := s.Called(section, id)

	if len(args) > 1 {
		*section = *args.Get(1).(*model.SongSection)
	}

	return args.Error(0)
}

func (s *SongSectionRepositoryMock) CountAllBySong(count *int64, songID uuid.UUID) error {
	args := s.Called(count, songID)

	if len(args) > 1 {
		*count = *args.Get(1).(*int64)
	}

	return args.Error(0)
}

func (s *SongSectionRepositoryMock) Create(section *model.SongSection) error {
	args := s.Called(section)
	return args.Error(0)
}

func (s *SongSectionRepositoryMock) Update(section *model.SongSection) error {
	args := s.Called(section)
	return args.Error(0)
}

func (s *SongSectionRepositoryMock) UpdateWithAssociations(section *model.SongSection) error {
	args := s.Called(section)
	return args.Error(0)
}

func (s *SongSectionRepositoryMock) Delete(ids []uuid.UUID) error {
	args := s.Called(ids)
	return args.Error(0)
}

// Types

func (s *SongSectionRepositoryMock) GetTypes(sectionTypes *[]model.SongSectionType, userID uuid.UUID) error {
	args := s.Called(sectionTypes, userID)

	if len(args) > 1 {
		*sectionTypes = *args.Get(1).(*[]model.SongSectionType)
	}

	return args.Error(0)
}
