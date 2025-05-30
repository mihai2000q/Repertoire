package search

import (
	"encoding/json"
	"errors"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/assert"
	"repertoire/server/domain/message/handler/storage"
	"repertoire/server/internal/wrapper"
	"repertoire/server/test/unit/data/logger"
	"repertoire/server/test/unit/data/service"
	"testing"
)

func TestDeleteDirectoriesStorageHandler_WhenDeleteDirectoriesFails_ShouldReturnError(t *testing.T) {
	// given
	storageService := new(service.StorageServiceMock)
	_uut := storage.NewDeleteDirectoriesStorageHandler(logger.NewLoggerMock(), storageService)

	directories := []string{"some_directory", "some_other_directory"}

	internalError := errors.New("internal error")
	storageService.On("DeleteDirectories", directories).
		Return(wrapper.InternalServerError(internalError)).
		Once()

	// when
	payload, _ := json.Marshal(&directories)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, err.Error(), internalError.Error())

	storageService.AssertExpectations(t)
}

func TestDeleteDirectoriesStorageHandler_WhenSuccessful_ShouldDeleteDirectories(t *testing.T) {
	// given
	storageService := new(service.StorageServiceMock)
	_uut := storage.NewDeleteDirectoriesStorageHandler(logger.NewLoggerMock(), storageService)

	directories := []string{"dir1/dir2/file.exe", "some_file.png", "an_image.jpeg"}

	storageService.On("DeleteDirectories", directories).
		Return(nil).
		Once()

	// when
	payload, _ := json.Marshal(&directories)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.NoError(t, err)

	storageService.AssertExpectations(t)
}
