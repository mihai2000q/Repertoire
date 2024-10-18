package repository

import (
	"repertoire/model"

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

func (a *ArtistRepositoryMock) GetAllByUser(
	artists *[]model.Artist,
	userId uuid.UUID,
	currentPage *int,
	pageSize *int,
) error {
	args := a.Called(artists, userId, currentPage, pageSize)

	if len(args) > 1 {
		*artists = *args.Get(1).(*[]model.Artist)
	}

	return args.Error(0)
}

func (a *ArtistRepositoryMock) GetAllByUserCount(count *int64, userId uuid.UUID) error {
	args := a.Called(count, userId)

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
