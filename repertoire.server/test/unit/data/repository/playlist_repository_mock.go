package repository

import (
	"repertoire/server/model"

	"github.com/stretchr/testify/mock"

	"github.com/google/uuid"
)

type PlaylistRepositoryMock struct {
	mock.Mock
}

func (p *PlaylistRepositoryMock) Get(playlist *model.Playlist, id uuid.UUID) error {
	args := p.Called(playlist, id)

	if len(args) > 1 {
		*playlist = *args.Get(1).(*model.Playlist)
	}

	return args.Error(0)
}

func (p *PlaylistRepositoryMock) GetPlaylistSongs(playlistSongs *[]model.PlaylistSong, id uuid.UUID) error {
	args := p.Called(playlistSongs, id)

	if len(args) > 1 {
		*playlistSongs = *args.Get(1).(*[]model.PlaylistSong)
	}

	return args.Error(0)
}

func (p *PlaylistRepositoryMock) GetWithAssociations(playlist *model.Playlist, id uuid.UUID) error {
	args := p.Called(playlist, id)

	if len(args) > 1 {
		*playlist = *args.Get(1).(*model.Playlist)
	}

	return args.Error(0)
}

func (p *PlaylistRepositoryMock) GetAllByUser(
	playlists *[]model.Playlist,
	userID uuid.UUID,
	currentPage *int,
	pageSize *int,
	orderBy []string,
	searchBy []string,
) error {
	args := p.Called(playlists, userID, currentPage, pageSize, orderBy, searchBy)

	if len(args) > 1 {
		*playlists = *args.Get(1).(*[]model.Playlist)
	}

	return args.Error(0)
}

func (p *PlaylistRepositoryMock) GetAllByUserCount(count *int64, userID uuid.UUID, searchBy []string) error {
	args := p.Called(count, userID, searchBy)

	if len(args) > 1 {
		*count = *args.Get(1).(*int64)
	}

	return args.Error(0)
}

func (p *PlaylistRepositoryMock) CountSongs(count *int64, id uuid.UUID) error {
	args := p.Called(count, id)

	if len(args) > 1 {
		*count = *args.Get(1).(*int64)
	}

	return args.Error(0)
}

func (p *PlaylistRepositoryMock) Create(playlist *model.Playlist) error {
	args := p.Called(playlist)
	return args.Error(0)
}

func (p *PlaylistRepositoryMock) AddSongs(playlistSongs *[]model.PlaylistSong) error {
	args := p.Called(playlistSongs)
	return args.Error(0)
}

func (p *PlaylistRepositoryMock) Update(playlist *model.Playlist) error {
	args := p.Called(playlist)
	return args.Error(0)
}

func (p *PlaylistRepositoryMock) UpdateAllPlaylistSongs(playlistSongs *[]model.PlaylistSong) error {
	args := p.Called(playlistSongs)
	return args.Error(0)
}

func (p *PlaylistRepositoryMock) Delete(id uuid.UUID) error {
	args := p.Called(id)
	return args.Error(0)
}

func (p *PlaylistRepositoryMock) RemoveSongs(playlistSongs *[]model.PlaylistSong) error {
	args := p.Called(playlistSongs)
	return args.Error(0)
}
