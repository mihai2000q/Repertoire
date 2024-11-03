package repository

import (
	"repertoire/server/model"

	"github.com/stretchr/testify/mock"

	"github.com/google/uuid"
)

type AlbumRepositoryMock struct {
	mock.Mock
}

func (a *AlbumRepositoryMock) Get(album *model.Album, id uuid.UUID) error {
	args := a.Called(album, id)

	if len(args) > 1 {
		*album = *args.Get(1).(*model.Album)
	}

	return args.Error(0)
}

func (a *AlbumRepositoryMock) GetWithSongs(album *model.Album, id uuid.UUID) error {
	args := a.Called(album, id)

	if len(args) > 1 {
		*album = *args.Get(1).(*model.Album)
	}

	return args.Error(0)
}

func (a *AlbumRepositoryMock) GetWithAssociations(album *model.Album, id uuid.UUID) error {
	args := a.Called(album, id)

	if len(args) > 1 {
		*album = *args.Get(1).(*model.Album)
	}

	return args.Error(0)
}

func (a *AlbumRepositoryMock) GetAllByUser(
	albums *[]model.Album,
	userID uuid.UUID,
	currentPage *int,
	pageSize *int,
	orderBy string,
) error {
	args := a.Called(albums, userID, currentPage, pageSize, orderBy)

	if len(args) > 1 {
		*albums = *args.Get(1).(*[]model.Album)
	}

	return args.Error(0)
}

func (a *AlbumRepositoryMock) GetAllByUserCount(count *int64, userID uuid.UUID) error {
	args := a.Called(count, userID)

	if len(args) > 1 {
		*count = *args.Get(1).(*int64)
	}

	return args.Error(0)
}

func (a *AlbumRepositoryMock) CountSongs(count *int64, albumID *uuid.UUID) error {
	args := a.Called(count, albumID)

	if len(args) > 1 {
		*count = *args.Get(1).(*int64)
	}

	return args.Error(0)
}

func (a *AlbumRepositoryMock) Create(album *model.Album) error {
	args := a.Called(album)
	return args.Error(0)
}

func (a *AlbumRepositoryMock) Update(album *model.Album) error {
	args := a.Called(album)
	return args.Error(0)
}

func (a *AlbumRepositoryMock) UpdateWithAssociations(album *model.Album) error {
	args := a.Called(album)
	return args.Error(0)
}

func (a *AlbumRepositoryMock) Delete(id uuid.UUID) error {
	args := a.Called(id)
	return args.Error(0)
}

func (a *AlbumRepositoryMock) RemoveSong(album *model.Album, song *model.Song) error {
	args := a.Called(album, song)
	return args.Error(0)
}
