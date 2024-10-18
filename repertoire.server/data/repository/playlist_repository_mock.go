package repository

import (
	"repertoire/models"

	"github.com/stretchr/testify/mock"

	"github.com/google/uuid"
)

type PlaylistRepositoryMock struct {
	mock.Mock
}

func (p *PlaylistRepositoryMock) Get(playlist *models.Playlist, id uuid.UUID) error {
	args := p.Called(playlist, id)

	if len(args) > 1 {
		*playlist = *args.Get(1).(*models.Playlist)
	}

	return args.Error(0)
}

func (p *PlaylistRepositoryMock) GetAllByUser(
	playlists *[]models.Playlist,
	userId uuid.UUID,
	currentPage *int,
	pageSize *int,
) error {
	args := p.Called(playlists, userId, currentPage, pageSize)

	if len(args) > 1 {
		*playlists = *args.Get(1).(*[]models.Playlist)
	}

	return args.Error(0)
}

func (p *PlaylistRepositoryMock) GetAllByUserCount(count *int64, userId uuid.UUID) error {
	args := p.Called(count, userId)

	if len(args) > 1 {
		*count = *args.Get(1).(*int64)
	}

	return args.Error(0)
}

func (p *PlaylistRepositoryMock) Create(playlist *models.Playlist) error {
	args := p.Called(playlist)
	return args.Error(0)
}

func (p *PlaylistRepositoryMock) Update(playlist *models.Playlist) error {
	args := p.Called(playlist)
	return args.Error(0)
}

func (p *PlaylistRepositoryMock) Delete(id uuid.UUID) error {
	args := p.Called(id)
	return args.Error(0)
}
