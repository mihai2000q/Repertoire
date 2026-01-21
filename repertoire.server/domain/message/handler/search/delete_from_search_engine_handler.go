package search

import (
	"encoding/json"
	"repertoire/server/data/logger"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"strconv"

	"github.com/ThreeDotsLabs/watermill/message"
)

type DeleteFromSearchEngineHandler struct {
	name                     string
	topic                    topics.Topic
	logger                   *logger.Logger
	searchEngineService      service.SearchEngineService
	searchTaskTrackerService service.SearchTaskTrackerService
}

func NewDeleteFromSearchEngineHandler(
	logger *logger.Logger,
	searchEngineService service.SearchEngineService,
	searchTaskTrackerService service.SearchTaskTrackerService,
) DeleteFromSearchEngineHandler {
	return DeleteFromSearchEngineHandler{
		name:                     "delete_from_search_engine_handler",
		topic:                    topics.DeleteFromSearchEngineTopic,
		logger:                   logger,
		searchEngineService:      searchEngineService,
		searchTaskTrackerService: searchTaskTrackerService,
	}
}

func (d DeleteFromSearchEngineHandler) Handle(msg *message.Message) error {
	var ids []string
	err := json.Unmarshal(msg.Payload, &ids)
	if err != nil {
		return err
	}

	document, err := d.searchEngineService.GetDocument(ids[0])
	if err != nil {
		return err
	}

	taskID, err := d.searchEngineService.Delete(ids)
	if err != nil {
		return err
	}
	d.searchTaskTrackerService.Track(strconv.FormatInt(taskID, 10), document["userId"].(string))
	d.logger.Debug("Search engine deleted " + strconv.Itoa(len(ids)) + " documents")
	return nil
}

func (d DeleteFromSearchEngineHandler) GetName() string {
	return d.name
}

func (d DeleteFromSearchEngineHandler) GetTopic() topics.Topic {
	return d.topic
}
