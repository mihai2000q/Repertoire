package repository

import (
	"repertoire/server/model"

	"github.com/stretchr/testify/mock"

	"github.com/google/uuid"
)

type SongArrangementRepositoryMock struct {
	mock.Mock
}

func (s *SongArrangementRepositoryMock) GetWithAssociations(arrangement *model.SongArrangement, id uuid.UUID) error {
	args := s.Called(arrangement, id)

	if len(args) > 1 {
		*arrangement = *args.Get(1).(*model.SongArrangement)
	}

	return args.Error(0)
}

func (s *SongArrangementRepositoryMock) GetAllBySong(arrangements *[]model.SongArrangement, songID uuid.UUID) error {
	args := s.Called(arrangements, songID)

	if len(args) > 1 {
		*arrangements = *args.Get(1).(*[]model.SongArrangement)
	}

	return args.Error(0)
}

func (s *SongArrangementRepositoryMock) CountBySong(count *int64, songID uuid.UUID) error {
	args := s.Called(count, songID)

	if len(args) > 1 {
		*count = *args.Get(1).(*int64)
	}

	return args.Error(0)
}

func (s *SongArrangementRepositoryMock) Create(arrangement *model.SongArrangement) error {
	args := s.Called(arrangement)
	return args.Error(0)
}

func (s *SongArrangementRepositoryMock) UpdateWithAssociations(arrangement *model.SongArrangement) error {
	args := s.Called(arrangement)
	return args.Error(0)
}

func (s *SongArrangementRepositoryMock) UpdateAllWithAssociations(arrangements *[]model.SongArrangement) error {
	args := s.Called(arrangements)
	return args.Error(0)
}

func (s *SongArrangementRepositoryMock) Delete(id uuid.UUID) error {
	args := s.Called(id)
	return args.Error(0)
}
