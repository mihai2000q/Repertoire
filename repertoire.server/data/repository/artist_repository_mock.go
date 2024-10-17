package repository

import (
	"repertoire/models"

	"github.com/stretchr/testify/mock"

	"github.com/google/uuid"
)

type ArtistRepositoryMock struct {
	mock.Mock
}

func (a *ArtistRepositoryMock) Get(artist *models.Artist, id uuid.UUID) error {
	args := a.Called(artist, id)

	if len(args) > 1 {
		*artist = *args.Get(1).(*models.Artist)
	}

	return args.Error(0)
}

func (a *ArtistRepositoryMock) GetAllByUser(
	artists *[]models.Artist,
	userId uuid.UUID,
	currentPage *int,
	pageSize *int,
) error {
	args := a.Called(artists, userId, currentPage, pageSize)

	if len(args) > 1 {
		*artists = *args.Get(1).(*[]models.Artist)
	}

	return args.Error(0)
}

func (a *ArtistRepositoryMock) Create(artist *models.Artist) error {
	args := a.Called(artist)
	return args.Error(0)
}

func (a *ArtistRepositoryMock) Update(artist *models.Artist) error {
	args := a.Called(artist)
	return args.Error(0)
}

func (a *ArtistRepositoryMock) Delete(id uuid.UUID) error {
	args := a.Called(id)
	return args.Error(0)
}
