package service

import (
	"mime/multipart"
	"repertoire/server/internal"

	"github.com/stretchr/testify/mock"
)

type StorageServiceMock struct {
	mock.Mock
}

func (s *StorageServiceMock) Upload(fileHeader *multipart.FileHeader, filePath string) error {
	args := s.Called(fileHeader, filePath)
	return args.Error(0)
}

func (s *StorageServiceMock) DeleteFile(filePath internal.FilePath) error {
	args := s.Called(filePath)
	return args.Error(0)
}

func (s *StorageServiceMock) DeleteDirectory(directoryPath string) error {
	args := s.Called(directoryPath)
	return args.Error(0)
}
