package provider

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"mime/multipart"
)

type StorageFilePathProviderMock struct {
	mock.Mock
}

func (s *StorageFilePathProviderMock) GetSongImagePathAndURL(file *multipart.FileHeader, songID uuid.UUID) (string, string) {
	args := s.Called(file, songID)
	return args.String(0), args.String(1)
}
