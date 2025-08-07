package transaction

import (
	"repertoire/server/data/repository"

	"github.com/stretchr/testify/mock"
)

type RepositoryFactoryMock struct {
	mock.Mock
}

func (m *RepositoryFactoryMock) NewArtistRepository() repository.ArtistRepository {
	args := m.Called()
	return args.Get(0).(repository.ArtistRepository)
}

func (m *RepositoryFactoryMock) NewAlbumRepository() repository.AlbumRepository {
	args := m.Called()
	return args.Get(0).(repository.AlbumRepository)
}

func (m *RepositoryFactoryMock) NewPlaylistRepository() repository.PlaylistRepository {
	args := m.Called()
	return args.Get(0).(repository.PlaylistRepository)
}

func (m *RepositoryFactoryMock) NewSongRepository() repository.SongRepository {
	args := m.Called()
	return args.Get(0).(repository.SongRepository)
}

func (m *RepositoryFactoryMock) NewUserDataRepository() repository.UserDataRepository {
	args := m.Called()
	return args.Get(0).(repository.UserDataRepository)
}

func (m *RepositoryFactoryMock) NewUserRepository() repository.UserRepository {
	args := m.Called()
	return args.Get(0).(repository.UserRepository)
}
