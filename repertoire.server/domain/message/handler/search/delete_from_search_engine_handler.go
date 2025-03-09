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

func (d DeleteFromSearchEngineHandler) Handle(msg *message.Message) error {
	var ids []string
	err := json.Unmarshal(msg.Payload, &ids)
	if err != nil {
		return err
	}

	err = d.searchEngineService.Delete(ids)
	return err
}

func (d DeleteFromSearchEngineHandler) GetName() string {
	return d.name
}

func (d DeleteFromSearchEngineHandler) GetTopic() topics.Topic {
	return d.topic
}
