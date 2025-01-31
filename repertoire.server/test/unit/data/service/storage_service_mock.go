package service

import (
	"mime/multipart"
	"repertoire/server/internal"
	"repertoire/server/internal/wrapper"

	"github.com/stretchr/testify/mock"
)

type StorageServiceMock struct {
	mock.Mock
}

func (s *StorageServiceMock) Upload(fileHeader *multipart.FileHeader, filePath string) error {
	args := s.Called(fileHeader, filePath)
	return args.Error(0)
}

func (s *StorageServiceMock) DeleteFile(filePath internal.FilePath) *wrapper.ErrorCode {
	args := s.Called(filePath)

	var errCode *wrapper.ErrorCode
	if a := args.Get(0); a != nil {
		errCode = a.(*wrapper.ErrorCode)
	}

	return errCode
}

func (s *StorageServiceMock) DeleteDirectory(directoryPath string) *wrapper.ErrorCode {
	args := s.Called(directoryPath)

	var errCode *wrapper.ErrorCode
	if a := args.Get(0); a != nil {
		errCode = a.(*wrapper.ErrorCode)
	}

	return errCode
}
