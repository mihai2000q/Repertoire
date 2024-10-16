package repository

import (
	"repertoire/models"

	"github.com/stretchr/testify/mock"

	"github.com/google/uuid"
)

type AlbumRepositoryMock struct {
	mock.Mock
}

func (s *AlbumRepositoryMock) Get(album *models.Album, id uuid.UUID) error {
	args := s.Called(album, id)

	if len(args) > 1 {
		*album = *args.Get(1).(*models.Album)
	}

	return args.Error(0)
}

func (s *AlbumRepositoryMock) GetAllByUser(albums *[]models.Album, userId uuid.UUID) error {
	args := s.Called(albums, userId)

	if len(args) > 1 {
		*albums = *args.Get(1).(*[]models.Album)
	}

	return args.Error(0)
}

func (s *AlbumRepositoryMock) Create(album *models.Album) error {
	args := s.Called(album)
	return args.Error(0)
}

func (s *AlbumRepositoryMock) Update(album *models.Album) error {
	args := s.Called(album)
	return args.Error(0)
}

func (s *AlbumRepositoryMock) Delete(id uuid.UUID) error {
	args := s.Called(id)
	return args.Error(0)
}
