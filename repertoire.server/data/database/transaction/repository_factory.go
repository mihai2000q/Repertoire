package transaction

import (
	"repertoire/server/data/database"
	"repertoire/server/data/repository"
)

type RepositoryFactory interface {
	NewArtistRepository() repository.ArtistRepository
	NewAlbumRepository() repository.AlbumRepository
	NewPlaylistRepository() repository.PlaylistRepository
	NewSongRepository() repository.SongRepository
	NewUserDataRepository() repository.UserDataRepository
	NewUserRepository() repository.UserRepository
}

type repositoryFactory struct {
	client database.Client
}

func (f repositoryFactory) NewArtistRepository() repository.ArtistRepository {
	return repository.NewArtistRepository(f.client)
}

func (f repositoryFactory) NewAlbumRepository() repository.AlbumRepository {
	return repository.NewAlbumRepository(f.client)
}

func (f repositoryFactory) NewPlaylistRepository() repository.PlaylistRepository {
	return repository.NewPlaylistRepository(f.client)
}

func (f repositoryFactory) NewSongRepository() repository.SongRepository {
	return repository.NewSongRepository(f.client)
}

func (f repositoryFactory) NewUserDataRepository() repository.UserDataRepository {
	return repository.NewUserDataRepository(f.client)
}

func (f repositoryFactory) NewUserRepository() repository.UserRepository {
	return repository.NewUserRepository(f.client)
}
