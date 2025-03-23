package search

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"repertoire/server/data/logger"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"strconv"
)

type DeleteFromSearchEngineHandler struct {
	name                string
	topic               topics.Topic
	logger              *logger.Logger
	searchEngineService service.SearchEngineService
}

func NewDeleteFromSearchEngineHandler(
	logger *logger.Logger,
	searchEngineService service.SearchEngineService,
) DeleteFromSearchEngineHandler {
	return DeleteFromSearchEngineHandler{
		name:                "delete_from_search_engine_handler",
		topic:               topics.DeleteFromSearchEngineTopic,
		logger:              logger,
		searchEngineService: searchEngineService,
	}
}

func (d DeleteFromSearchEngineHandler) Handle(msg *message.Message) error {
	var ids []string
	err := json.Unmarshal(msg.Payload, &ids)
	if err != nil {
		return err
	}

	err = d.searchEngineService.Delete(ids)
	if err != nil {
		return err
	}
	d.logger.Debug("Search engine deleted " + strconv.Itoa(len(ids)) + " documents")
	return nil
}

func (d DeleteFromSearchEngineHandler) GetName() string {
	return d.name
}

func (d DeleteFromSearchEngineHandler) GetTopic() topics.Topic {
	return d.topic
}
