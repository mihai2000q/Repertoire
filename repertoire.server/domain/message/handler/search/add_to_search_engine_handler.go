package search

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"repertoire/server/data/logger"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"strconv"
)

type AddToSearchEngineHandler struct {
	name                     string
	topic                    topics.Topic
	logger                   *logger.Logger
	searchEngineService      service.SearchEngineService
	searchTaskTrackerService service.SearchTaskTrackerService
}

func NewAddToSearchEngineHandler(
	logger *logger.Logger,
	searchEngineService service.SearchEngineService,
	searchTaskTrackerService service.SearchTaskTrackerService,
) AddToSearchEngineHandler {
	return AddToSearchEngineHandler{
		name:                     "add_to_search_engine_handler",
		topic:                    topics.AddToSearchEngineTopic,
		logger:                   logger,
		searchEngineService:      searchEngineService,
		searchTaskTrackerService: searchTaskTrackerService,
	}
}

func (a AddToSearchEngineHandler) Handle(msg *message.Message) error {
	var documents []map[string]any // the documents can also be of type []any, because they will be unmarshalled
	err := json.Unmarshal(msg.Payload, &documents)
	if err != nil {
		return err
	}

	taskID, err := a.searchEngineService.Add(documents)
	if err != nil {
		return err
	}
	a.searchTaskTrackerService.Track(strconv.FormatInt(taskID, 10), documents[0]["userId"].(string))
	a.logger.Debug("Search engine added " + strconv.Itoa(len(documents)) + " documents")
	return nil
}

func (a AddToSearchEngineHandler) GetName() string {
	return a.name
}

func (a AddToSearchEngineHandler) GetTopic() topics.Topic {
	return a.topic
}
