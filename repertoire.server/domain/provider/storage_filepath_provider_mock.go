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

func (s *StorageFilePathProviderMock) GetAlbumImagePath(file *multipart.FileHeader, album model.Album) string {
	args := s.Called(file, album)
	return args.String(0)
}

func (s *StorageFilePathProviderMock) GetSongImagePath(file *multipart.FileHeader, songID uuid.UUID) string {
	args := s.Called(file, songID)
	return args.String(0)
}
