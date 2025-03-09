package search

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
)

type UpdateFromSearchEngineHandler struct {
	name                string
	topic               topics.Topic
	searchEngineService service.SearchEngineService
}

func NewUpdateFromSearchEngineHandler(searchEngineService service.SearchEngineService) UpdateFromSearchEngineHandler {
	return UpdateFromSearchEngineHandler{
		name:                "update_from_search_engine_handler",
		topic:               topics.UpdateFromSearchEngineTopic,
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
	return err
}

func (u UpdateFromSearchEngineHandler) GetName() string {
	return u.name
}

func (u UpdateFromSearchEngineHandler) GetTopic() topics.Topic {
	return u.topic
}
