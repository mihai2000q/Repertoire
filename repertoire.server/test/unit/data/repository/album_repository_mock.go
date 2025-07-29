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

func (a *AlbumRepositoryMock) GetWithSongsAndArtist(album *model.Album, id uuid.UUID) error {
	args := a.Called(album, id)

	if len(args) > 1 {
		*album = *args.Get(1).(*model.Album)
	}

	return args.Error(0)
}

func (a *AlbumRepositoryMock) GetWithAssociations(album *model.Album, id uuid.UUID, songsOrderBy []string) error {
	args := a.Called(album, id, songsOrderBy)

	if len(args) > 1 {
		*album = *args.Get(1).(*model.Album)
	}

	return args.Error(0)
}

func (a *AlbumRepositoryMock) GetFiltersMetadata(
	metadata *model.AlbumFiltersMetadata,
	userID uuid.UUID,
	searchBy []string,
) error {
	args := a.Called(metadata, userID, searchBy)

	if len(args) > 1 {
		*metadata = *args.Get(1).(*model.AlbumFiltersMetadata)
	}

	return args.Error(0)
}

func (a *AlbumRepositoryMock) GetAllByIDsWithSongs(albums *[]model.Album, ids []uuid.UUID) error {
	args := a.Called(albums, ids)

	if len(args) > 1 {
		*albums = *args.Get(1).(*[]model.Album)
	}

	return args.Error(0)
}

func (a *AlbumRepositoryMock) GetAllByIDsWithSongsAndArtist(albums *[]model.Album, ids []uuid.UUID) error {
	args := a.Called(albums, ids)

	if len(args) > 1 {
		*albums = *args.Get(1).(*[]model.Album)
	}

	return args.Error(0)
}

func (a *AlbumRepositoryMock) GetAllByUser(
	albums *[]model.EnhancedAlbum,
	userID uuid.UUID,
	currentPage *int,
	pageSize *int,
	orderBy []string,
	searchBy []string,
) error {
	args := a.Called(albums, userID, currentPage, pageSize, orderBy, searchBy)

	if len(args) > 1 {
		*albums = *args.Get(1).(*[]model.EnhancedAlbum)
	}

	return args.Error(0)
}

func (a *AlbumRepositoryMock) GetAllByUserCount(count *int64, userID uuid.UUID, searchBy []string) error {
	args := a.Called(count, userID, searchBy)

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

func (a *AlbumRepositoryMock) UpdateAllWithSongs(albums *[]model.Album) error {
	args := a.Called(albums)
	return args.Error(0)
}

func (a *AlbumRepositoryMock) Delete(ids []uuid.UUID) error {
	args := a.Called(ids)
	return args.Error(0)
}

func (a *AlbumRepositoryMock) DeleteWithSongs(ids []uuid.UUID) error {
	args := a.Called(ids)
	return args.Error(0)
}

func (a *AlbumRepositoryMock) RemoveSongs(album *model.Album, songs *[]model.Song) error {
	args := a.Called(album, songs)
	return args.Error(0)
}
