package service

import (
	"mime/multipart"

	"github.com/stretchr/testify/mock"
)

type StorageServiceMock struct {
	mock.Mock
}

func (s *StorageServiceMock) Upload(fileHeader *multipart.FileHeader, filePath string) error {
	args := s.Called(fileHeader, filePath)
	return args.Error(0)
}

func (s *StorageServiceMock) Delete(filePath string) error {
	args := s.Called(filePath)
	return args.Error(0)
}
