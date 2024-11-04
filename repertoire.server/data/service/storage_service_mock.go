package service

import (
	"mime/multipart"

	"github.com/stretchr/testify/mock"
)

type StorageServiceMock struct {
	mock.Mock
}

func (s *StorageServiceMock) Upload(token string, fileHeader *multipart.FileHeader, filePath string) error {
	args := s.Called(token, fileHeader, filePath)
	return args.Error(0)
}

func (s *StorageServiceMock) Delete(filePath string) error {
	args := s.Called(filePath)
	return args.Error(0)
}
