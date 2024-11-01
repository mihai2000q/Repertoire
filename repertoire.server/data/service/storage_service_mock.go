package service

import (
	"github.com/stretchr/testify/mock"
	"mime/multipart"
)

type StorageServiceMock struct {
	mock.Mock
}

func (s *StorageServiceMock) Upload(token string, fileHeader *multipart.FileHeader, filePath string) error {
	args := s.Called(token, fileHeader, filePath)
	return args.Error(0)
}

func (s *StorageServiceMock) Delete(token string, filePath string) error {
	args := s.Called(token, filePath)
	return args.Error(0)
}
