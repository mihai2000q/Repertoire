package repository

import (
	"repertoire/server/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type SongPartRepositoryMock struct {
	mock.Mock
}

func (s *SongPartRepositoryMock) GetWithSong(part *model.SongPart, id uuid.UUID) error {
	args := s.Called(part, id)

	if len(args) > 1 {
		*part = *args.Get(1).(*model.SongPart)
	}

	return args.Error(0)
}

func (s *SongPartRepositoryMock) GetAllByIDs(parts *[]model.SongPart, ids []uuid.UUID) error {
	args := s.Called(parts, ids)

	if len(args) > 1 {
		*parts = *args.Get(1).(*[]model.SongPart)
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

func (s *SongPartRepositoryMock) CountBySectionIDs(partIDs []uuid.UUID) (map[uuid.UUID]int64, error) {
	args := s.Called(partIDs)
	return args.Get(0).(map[uuid.UUID]int64), args.Error(1)
}

func (s *SongPartRepositoryMock) Create(part *model.SongPart) error {
	args := s.Called(part)
	return args.Error(0)
}

func (s *SongPartRepositoryMock) Update(part *model.SongPart) error {
	args := s.Called(part)
	return args.Error(0)
}

func (s *SongPartRepositoryMock) UpdateAll(parts *[]model.SongPart) error {
	args := s.Called(parts)
	return args.Error(0)
}

func (s *SongPartRepositoryMock) UpdateWithAssociations(part *model.SongPart) error {
	args := s.Called(part)
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
