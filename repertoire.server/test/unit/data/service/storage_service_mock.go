package service

import (
	"mime/multipart"
	"repertoire/server/internal"
	"repertoire/server/internal/httperror"

	"github.com/stretchr/testify/mock"
)

type StorageServiceMock struct {
	mock.Mock
}

func (s *StorageServiceMock) Upload(fileHeader *multipart.FileHeader, filePath string) error {
	args := s.Called(fileHeader, filePath)
	return args.Error(0)
}

func (s *StorageServiceMock) DeleteFile(filePath internal.FilePath) *httperror.ErrorCode {
	args := s.Called(filePath)

	var errCode *httperror.ErrorCode
	if a := args.Get(0); a != nil {
		errCode = a.(*httperror.ErrorCode)
	}

	return errCode
}

func (s *StorageServiceMock) DeleteDirectories(directoryPaths []string) *httperror.ErrorCode {
	args := s.Called(directoryPaths)

	var errCode *httperror.ErrorCode
	if a := args.Get(0); a != nil {
		errCode = a.(*httperror.ErrorCode)
	}

	return errCode
}
