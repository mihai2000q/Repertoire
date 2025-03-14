package user

import (
	"github.com/google/uuid"
)

func GetSearchDocuments() []any {
	var documents []any
	for _, search := range Searches {
		documents = append(documents, search)
	}
	return documents
}

var UserSearchID = uuid.New()

var Searches = []map[string]any{
	{
		"id":     "artist-some-id",
		"userId": UserSearchID,
	},
	{
		"id":     "album-some-id",
		"userId": UserSearchID,
	},
	{
		"id":     "song-some-id",
		"userId": UserSearchID,
	},
	{
		"id":     "playlist-some-id",
		"userId": UserSearchID,
	},
}
