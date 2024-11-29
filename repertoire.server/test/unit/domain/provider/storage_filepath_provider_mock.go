package provider

import (
	"mime/multipart"
	"repertoire/server/model"

	"github.com/stretchr/testify/mock"
)

type StorageFilePathProviderMock struct {
	mock.Mock
}

func (s *StorageFilePathProviderMock) GetUserProfilePicturePath(file *multipart.FileHeader, user model.User) string {
	args := s.Called(file, user)
	return args.String(0)
}

func (s *StorageFilePathProviderMock) GetAlbumImagePath(file *multipart.FileHeader, album model.Album) string {
	args := s.Called(file, album)
	return args.String(0)
}

func (s *StorageFilePathProviderMock) GetArtistImagePath(file *multipart.FileHeader, artist model.Artist) string {
	args := s.Called(file, artist)
	return args.String(0)
}

func (s *StorageFilePathProviderMock) GetPlaylistImagePath(file *multipart.FileHeader, playlist model.Playlist) string {
	args := s.Called(file, playlist)
	return args.String(0)
}

func (s *StorageFilePathProviderMock) GetSongImagePath(file *multipart.FileHeader, song model.Song) string {
	args := s.Called(file, song)
	return args.String(0)
}
