package storage

import (
	"encoding/json"
	"errors"
	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"net/http"
	"repertoire/server/data/logger"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
)

type DeleteDirectoriesStorageHandler struct {
	name           string
	topic          topics.Topic
	logger         *logger.Logger
	storageService service.StorageService
}

func NewDeleteDirectoriesStorageHandler(
	logger *logger.Logger,
	storageService service.StorageService,
) DeleteDirectoriesStorageHandler {
	return DeleteDirectoriesStorageHandler{
		name:           "delete_directories_storage_handler",
		topic:          topics.DeleteDirectoriesStorageTopic,
		logger:         logger,
		storageService: storageService,
	}
}

func (d DeleteDirectoriesStorageHandler) Handle(msg *watermillMessage.Message) error {
	var directoryPaths []string
	err := json.Unmarshal(msg.Payload, &directoryPaths)
	if err != nil {
		return err
	}

	var resultErrors []error
	for _, path := range directoryPaths {
		errCode := d.storageService.DeleteDirectory(path)
		if errCode != nil && errCode.Code == http.StatusNotFound {
			d.logger.Debug("Directory not found: " + path)
		} else if errCode != nil {
			resultErrors = append(resultErrors, errCode.Error)
		}
	}

	if len(resultErrors) > 0 {
		return errors.Join(resultErrors...)
	} else {
		return nil
	}
}

func (d DeleteDirectoriesStorageHandler) GetName() string {
	return d.name
}

func (d DeleteDirectoriesStorageHandler) GetTopic() topics.Topic {
	return d.topic
}
