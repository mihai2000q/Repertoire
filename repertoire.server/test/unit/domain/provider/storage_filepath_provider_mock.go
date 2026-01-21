package provider

import (
	"mime/multipart"
	"repertoire/server/model"

	"github.com/google/uuid"

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

func (s *StorageFilePathProviderMock) GetBandMemberImagePath(file *multipart.FileHeader, artist model.BandMember) string {
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

func (s *StorageFilePathProviderMock) GetUserDirectoryPath(userID uuid.UUID) string {
	args := s.Called(userID)
	return args.String(0)
}

func (s *StorageFilePathProviderMock) GetAlbumDirectoryPath(album model.Album) string {
	args := s.Called(album)
	return args.String(0)
}

func (s *StorageFilePathProviderMock) GetArtistDirectoryPath(artist model.Artist) string {
	args := s.Called(artist)
	return args.String(0)
}

func (s *StorageFilePathProviderMock) GetPlaylistDirectoryPath(playlist model.Playlist) string {
	args := s.Called(playlist)
	return args.String(0)
}

func (s *StorageFilePathProviderMock) GetSongDirectoryPath(song model.Song) string {
	args := s.Called(song)
	return args.String(0)
}

func (s *StorageFilePathProviderMock) HasAlbumFiles(album model.Album) bool {
	args := s.Called(album)
	return args.Bool(0)
}

func (s *StorageFilePathProviderMock) HasArtistFiles(artist model.Artist) bool {
	args := s.Called(artist)
	return args.Bool(0)
}

func (s *StorageFilePathProviderMock) HasPlaylistFiles(playlist model.Playlist) bool {
	args := s.Called(playlist)
	return args.Bool(0)
}

func (s *StorageFilePathProviderMock) HasSongFiles(song model.Song) bool {
	args := s.Called(song)
	return args.Bool(0)
}
