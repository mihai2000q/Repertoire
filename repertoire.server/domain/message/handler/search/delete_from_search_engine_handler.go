package search

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
)

type DeleteFromSearchEngineHandler struct {
	name                string
	topic               topics.Topic
	searchEngineService service.SearchEngineService
}

func NewDeleteFromSearchEngineHandler(searchEngineService service.SearchEngineService) DeleteFromSearchEngineHandler {
	return DeleteFromSearchEngineHandler{
		name:                "delete_from_search_engine_handler",
		topic:               topics.DeleteFromSearchEngineTopic,
		searchEngineService: searchEngineService,
	}
}

func (s DeleteFromSearchEngineHandler) Handle(msg *message.Message) error {
	var ids []string
	err := json.Unmarshal(msg.Payload, &ids)
	if err != nil {
		return err
	}

	err = s.searchEngineService.Delete(ids)
	return err
}

func (s DeleteFromSearchEngineHandler) GetName() string {
	return s.name
}

func (s DeleteFromSearchEngineHandler) GetTopic() topics.Topic {
	return s.topic
}
