package repository

import (
	"repertoire/models"

	"github.com/stretchr/testify/mock"

	"github.com/google/uuid"
)

type ArtistRepositoryMock struct {
	mock.Mock
}

func (s *ArtistRepositoryMock) Get(artist *models.Artist, id uuid.UUID) error {
	args := s.Called(artist, id)

	if len(args) > 1 {
		*artist = *args.Get(1).(*models.Artist)
	}

	return args.Error(0)
}

func (s *ArtistRepositoryMock) GetAllByUser(artists *[]models.Artist, userId uuid.UUID) error {
	args := s.Called(artists, userId)

	if len(args) > 1 {
		*artists = *args.Get(1).(*[]models.Artist)
	}

	return args.Error(0)
}

func (s *ArtistRepositoryMock) Create(artist *models.Artist) error {
	args := s.Called(artist)
	return args.Error(0)
}

func (s *ArtistRepositoryMock) Update(artist *models.Artist) error {
	args := s.Called(artist)
	return args.Error(0)
}

func (s *ArtistRepositoryMock) Delete(id uuid.UUID) error {
	args := s.Called(id)
	return args.Error(0)
}
