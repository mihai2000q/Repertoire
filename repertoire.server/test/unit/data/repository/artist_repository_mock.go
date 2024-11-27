package repository

import (
	"repertoire/server/model"

	"github.com/stretchr/testify/mock"

	"github.com/google/uuid"
)

type ArtistRepositoryMock struct {
	mock.Mock
}

func (a *ArtistRepositoryMock) Get(artist *model.Artist, id uuid.UUID) error {
	args := a.Called(artist, id)

	if len(args) > 1 {
		*artist = *args.Get(1).(*model.Artist)
	}

	return args.Error(0)
}

func (a *ArtistRepositoryMock) GetWithAssociations(artist *model.Artist, id uuid.UUID) error {
	args := a.Called(artist, id)

	if len(args) > 1 {
		*artist = *args.Get(1).(*model.Artist)
	}

	return args.Error(0)
}

func (a *ArtistRepositoryMock) GetAllByUser(
	artists *[]model.Artist,
	userID uuid.UUID,
	currentPage *int,
	pageSize *int,
	orderBy []string,
	searchBy []string,
) error {
	args := a.Called(artists, userID, currentPage, pageSize, orderBy, searchBy)

	if len(args) > 1 {
		*artists = *args.Get(1).(*[]model.Artist)
	}

	return args.Error(0)
}

func (a *ArtistRepositoryMock) GetAllByUserCount(count *int64, userID uuid.UUID, searchBy []string) error {
	args := a.Called(count, userID, searchBy)

	if len(args) > 1 {
		*count = *args.Get(1).(*int64)
	}

	return args.Error(0)
}

func (a *ArtistRepositoryMock) Create(artist *model.Artist) error {
	args := a.Called(artist)
	return args.Error(0)
}

func (a *ArtistRepositoryMock) Update(artist *model.Artist) error {
	args := a.Called(artist)
	return args.Error(0)
}

func (a *ArtistRepositoryMock) Delete(id uuid.UUID) error {
	args := a.Called(id)
	return args.Error(0)
}
