package search

import (
	"encoding/json"
	"repertoire/server/data/logger"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"strconv"

	"github.com/ThreeDotsLabs/watermill/message"
)

type UpdateFromSearchEngineHandler struct {
	name                     string
	topic                    topics.Topic
	logger                   *logger.Logger
	searchEngineService      service.SearchEngineService
	searchTaskTrackerService service.SearchTaskTrackerService
}

func NewUpdateFromSearchEngineHandler(
	logger *logger.Logger,
	searchEngineService service.SearchEngineService,
	searchTaskTrackerService service.SearchTaskTrackerService,
) UpdateFromSearchEngineHandler {
	return UpdateFromSearchEngineHandler{
		name:                     "update_from_search_engine_handler",
		topic:                    topics.UpdateFromSearchEngineTopic,
		logger:                   logger,
		searchEngineService:      searchEngineService,
		searchTaskTrackerService: searchTaskTrackerService,
	}
}

func (u UpdateFromSearchEngineHandler) Handle(msg *message.Message) error {
	var documents []map[string]any // the documents can also be of type []any, because they will be unmarshalled
	err := json.Unmarshal(msg.Payload, &documents)
	if err != nil {
		return err
	}

	taskID, err := u.searchEngineService.Update(documents)
	if err != nil {
		return err
	}
	u.searchTaskTrackerService.Track(strconv.FormatInt(taskID, 10), documents[0]["userId"].(string))
	u.logger.Debug("Search engine updated " + strconv.Itoa(len(documents)) + " documents")
	return nil
}

func (u UpdateFromSearchEngineHandler) GetName() string {
	return u.name
}

func (u UpdateFromSearchEngineHandler) GetTopic() topics.Topic {
	return u.topic
}
