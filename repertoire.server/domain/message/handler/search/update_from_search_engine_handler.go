package search

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"repertoire/server/data/logger"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
)

type UpdateFromSearchEngineHandler struct {
	name                string
	topic               topics.Topic
	logger              *logger.Logger
	searchEngineService service.SearchEngineService
}

func NewUpdateFromSearchEngineHandler(
	logger *logger.Logger,
	searchEngineService service.SearchEngineService,
) UpdateFromSearchEngineHandler {
	return UpdateFromSearchEngineHandler{
		name:                "update_from_search_engine_handler",
		topic:               topics.UpdateFromSearchEngineTopic,
		logger:              logger,
		searchEngineService: searchEngineService,
	}
}

func (u UpdateFromSearchEngineHandler) Handle(msg *message.Message) error {
	var documents []any
	err := json.Unmarshal(msg.Payload, &documents)
	if err != nil {
		return err
	}

	err = u.searchEngineService.Update(documents)
	if err != nil {
		return err
	}
	u.logger.Debug("Search engine updated " + string(rune(len(documents))) + " documents")
	return nil
}

func (u UpdateFromSearchEngineHandler) GetName() string {
	return u.name
}

func (u UpdateFromSearchEngineHandler) GetTopic() topics.Topic {
	return u.topic
}
