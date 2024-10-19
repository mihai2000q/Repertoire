package repository

import (
	"repertoire/model"

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

func (a *AlbumRepositoryMock) GetAllByUser(
	albums *[]model.Album,
	userId uuid.UUID,
	currentPage *int,
	pageSize *int,
	orderBy string,
) error {
	args := a.Called(albums, userId, currentPage, pageSize)

	if len(args) > 1 {
		*albums = *args.Get(1).(*[]model.Album)
	}

	return args.Error(0)
}

func (a *AlbumRepositoryMock) GetAllByUserCount(count *int64, userId uuid.UUID) error {
	args := a.Called(count, userId)

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

func (a *AlbumRepositoryMock) Delete(id uuid.UUID) error {
	args := a.Called(id)
	return args.Error(0)
}
