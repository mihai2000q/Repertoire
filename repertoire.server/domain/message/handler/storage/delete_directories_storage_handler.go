package storage

import (
	"encoding/json"
	"net/http"
	"repertoire/server/data/logger"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"

	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/zap"
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

	errCode := d.storageService.DeleteDirectories(directoryPaths)
	if errCode != nil && errCode.Code == http.StatusNotFound {
		d.logger.Debug("Directory not found", zap.Error(errCode.Error))
	} else if errCode != nil {
		return errCode.Error
	}

	return nil
}

func (d DeleteDirectoriesStorageHandler) GetName() string {
	return d.name
}

func (d DeleteDirectoriesStorageHandler) GetTopic() topics.Topic {
	return d.topic
}
