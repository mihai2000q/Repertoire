package provider

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"mime/multipart"
)

type StorageFilePathProviderMock struct {
	mock.Mock
}

func (s *StorageFilePathProviderMock) GetSongImagePath(file *multipart.FileHeader, songID uuid.UUID) string {
	args := s.Called(file, songID)
	return args.String(0)
}
