package topics

import (
	"repertoire/server/internal/message/queues"
)

type Topic string

const (
	AlbumCreatedTopic Topic = "album_created_topic"
	AlbumDeletedTopic Topic = "album_deleted_topic"

	ArtistCreatedTopic Topic = "artist_created_topic"

	PlaylistCreatedTopic Topic = "playlist_created_topic"
	PlaylistDeletedTopic Topic = "playlist_deleted_topic"

	SongCreatedTopic Topic = "song_created_topic"
	SongDeletedTopic Topic = "song_deleted_topic"
	SongUpdatedTopic Topic = "song_updated_topic"

	AddToSearchEngineTopic      Topic = "add_to_search_engine_topic"
	DeleteFromSearchEngineTopic Topic = "delete_from_search_engine_topic"
	UpdateFromSearchEngineTopic Topic = "update_from_search_engine_topic"
)

var TopicToQueueMap = map[Topic]queues.Queue{
	AlbumCreatedTopic: queues.MainQueue,
	AlbumDeletedTopic: queues.MainQueue,

	ArtistCreatedTopic: queues.MainQueue,

	PlaylistCreatedTopic: queues.MainQueue,
	PlaylistDeletedTopic: queues.MainQueue,

	SongCreatedTopic: queues.MainQueue,
	SongDeletedTopic: queues.MainQueue,

	AddToSearchEngineTopic:      queues.SearchQueue,
	DeleteFromSearchEngineTopic: queues.SearchQueue,
	UpdateFromSearchEngineTopic: queues.SearchQueue,
}
