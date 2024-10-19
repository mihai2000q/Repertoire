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

func (s *SongRepositoryMock) GetAllByUser(
	songs *[]model.Song,
	userId uuid.UUID,
	currentPage *int,
	pageSize *int,
	orderBy string,
) error {
	args := s.Called(songs, userId, currentPage, pageSize, orderBy)

	if len(args) > 1 {
		*songs = *args.Get(1).(*[]model.Song)
	}

	return args.Error(0)
}

func (s *SongRepositoryMock) GetAllByUserCount(count *int64, userId uuid.UUID) error {
	args := s.Called(count, userId)

	if len(args) > 1 {
		*count = *args.Get(1).(*int64)
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
