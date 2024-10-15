package repository

import (
	"github.com/stretchr/testify/mock"
	"repertoire/models"

	"github.com/google/uuid"
)

type SongRepositoryMock struct {
	mock.Mock
}

func (s *SongRepositoryMock) Get(song *models.Song, id uuid.UUID) error {
	args := s.Called(song, id)

	if len(args) > 1 {
		*song = *args.Get(1).(*models.Song)
	}

	return args.Error(0)
}

func (s *SongRepositoryMock) GetAllByUser(songs *[]models.Song, userId uuid.UUID) error {
	args := s.Called(songs, userId)

	if len(args) > 1 {
		*songs = *args.Get(1).(*[]models.Song)
	}

	return args.Error(0)
}

func (s *SongRepositoryMock) Create(song *models.Song) error {
	args := s.Called(song)
	return args.Error(0)
}

func (s *SongRepositoryMock) Update(song *models.Song) error {
	args := s.Called(song)
	return args.Error(0)
}

func (s *SongRepositoryMock) Delete(id uuid.UUID) error {
	args := s.Called(id)
	return args.Error(0)
}