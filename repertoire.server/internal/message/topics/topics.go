package topics

import (
	"repertoire/server/internal/message/queues"
)

type Topic string

const (
	AlbumCreatedTopic  Topic = "album_created_topic"
	AlbumDeletedTopic  Topic = "album_deleted_topic"
	AlbumsUpdatedTopic Topic = "albums_updated_topic"

	ArtistCreatedTopic Topic = "artist_created_topic"
	ArtistDeletedTopic Topic = "artist_deleted_topic"
	ArtistUpdatedTopic Topic = "artist_updated_topic"

	PlaylistCreatedTopic Topic = "playlist_created_topic"
	PlaylistDeletedTopic Topic = "playlist_deleted_topic"
	PlaylistUpdatedTopic Topic = "playlist_updated_topic"

	SongCreatedTopic  Topic = "song_created_topic"
	SongDeletedTopic  Topic = "song_deleted_topic"
	SongsUpdatedTopic Topic = "songs_updated_topic"

	UserDeletedTopic Topic = "user_deleted_topic"

	AddToSearchEngineTopic      Topic = "add_to_search_engine_topic"
	DeleteFromSearchEngineTopic Topic = "delete_from_search_engine_topic"
	UpdateFromSearchEngineTopic Topic = "update_from_search_engine_topic"
)

var TopicToQueueMap = map[Topic]queues.Queue{
	AlbumCreatedTopic:  queues.MainQueue,
	AlbumDeletedTopic:  queues.MainQueue,
	AlbumsUpdatedTopic: queues.MainQueue,

	ArtistCreatedTopic: queues.MainQueue,
	ArtistDeletedTopic: queues.MainQueue,
	ArtistUpdatedTopic: queues.MainQueue,

	PlaylistCreatedTopic: queues.MainQueue,
	PlaylistDeletedTopic: queues.MainQueue,
	PlaylistUpdatedTopic: queues.MainQueue,

	SongCreatedTopic: queues.MainQueue,
	SongDeletedTopic: queues.MainQueue,

	UserDeletedTopic: queues.MainQueue,

	AddToSearchEngineTopic:      queues.SearchQueue,
	DeleteFromSearchEngineTopic: queues.SearchQueue,
	UpdateFromSearchEngineTopic: queues.SearchQueue,
}
