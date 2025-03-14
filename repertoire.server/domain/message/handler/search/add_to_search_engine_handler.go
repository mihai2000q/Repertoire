package search

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
)

type AddToSearchEngineHandler struct {
	name                string
	topic               topics.Topic
	searchEngineService service.SearchEngineService
}

func NewAddToSearchEngineHandler(searchEngineService service.SearchEngineService) AddToSearchEngineHandler {
	return AddToSearchEngineHandler{
		name:                "add_to_search_engine_handler",
		topic:               topics.AddToSearchEngineTopic,
		searchEngineService: searchEngineService,
	}
}

func (a AddToSearchEngineHandler) Handle(msg *message.Message) error {
	var searches []any
	err := json.Unmarshal(msg.Payload, &searches)
	if err != nil {
		return err
	}

	err = a.searchEngineService.Add(searches)
	return err
}

func (a AddToSearchEngineHandler) GetName() string {
	return a.name
}

func (a AddToSearchEngineHandler) GetTopic() topics.Topic {
	return a.topic
}
