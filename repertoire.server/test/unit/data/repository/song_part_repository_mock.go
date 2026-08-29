package repository

import (
	"repertoire/server/model"

	"github.com/stretchr/testify/mock"

	"github.com/google/uuid"
)

type SongPartRepositoryMock struct {
	mock.Mock
}

func (s *SongPartRepositoryMock) Get(section *model.SongPart, id uuid.UUID) error {
	args := s.Called(section, id)

	if len(args) > 1 {
		*section = *args.Get(1).(*model.SongPart)
	}

	return args.Error(0)
}

func (s *SongPartRepositoryMock) CountAllBySection(count *int64, songID uuid.UUID) error {
	args := s.Called(count, songID)

	if len(args) > 1 {
		*count = *args.Get(1).(*int64)
	}

	return args.Error(0)
}

func (s *SongPartRepositoryMock) CountAllBySong(count *int64, songID uuid.UUID) error {
	args := s.Called(count, songID)

	if len(args) > 1 {
		*count = *args.Get(1).(*int64)
	}

	return args.Error(0)
}

func (s *SongPartRepositoryMock) Create(section *model.SongPart) error {
	args := s.Called(section)
	return args.Error(0)
}

func (s *SongPartRepositoryMock) Update(section *model.SongPart) error {
	args := s.Called(section)
	return args.Error(0)
}

func (s *SongPartRepositoryMock) UpdateWithAssociations(section *model.SongPart) error {
	args := s.Called(section)
	return args.Error(0)
}

func (s *SongPartRepositoryMock) Delete(ids []uuid.UUID) error {
	args := s.Called(ids)
	return args.Error(0)
}

// History

func (s *SongPartRepositoryMock) GetHistory(
	history *[]model.SongPartHistory,
	partID uuid.UUID,
	property model.SongPartProperty,
) error {
	args := s.Called(history, partID, property)

	if len(args) > 1 {
		*history = *args.Get(1).(*[]model.SongPartHistory)
	}

	return args.Error(0)
}

func (s *SongPartRepositoryMock) CreateHistory(history *model.SongPartHistory) error {
	args := s.Called(history)
	return args.Error(0)
}
