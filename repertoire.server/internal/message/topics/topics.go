package topics

import (
	"repertoire/server/internal/message/queues"
)

type Topic string

const (
	ArtistCreatedTopic Topic = "artist_created_topic"

	AddToSearchEngineTopic      Topic = "add_to_search_engine_topic"
	DeleteFromSearchEngineTopic Topic = "delete_from_search_engine_topic"
)

var TopicToQueueMap = map[Topic]queues.Queue{
	ArtistCreatedTopic: queues.MainQueue,

	AddToSearchEngineTopic:      queues.SearchQueue,
	DeleteFromSearchEngineTopic: queues.SearchQueue,
}
