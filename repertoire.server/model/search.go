package model

import (
	"github.com/google/uuid"
	"repertoire/server/internal"
	"repertoire/server/internal/enums"
)

type SearchBase struct {
	ID     uuid.UUID        `json:"id"`
	Type   enums.SearchType `json:"type"`
	UserID uuid.UUID        `json:"-"`
}

type ArtistSearch struct {
	ImageUrl *internal.FilePath `json:"imageUrl"`
	Name     string             `json:"name"`
	SearchBase
}

type AlbumSearch struct {
	ImageUrl *internal.FilePath `json:"imageUrl"`
	Title    string             `json:"title"`
	Artist   *ArtistSearch      `json:"artist"`
	SearchBase
}

type SongSearch struct {
	ImageUrl *internal.FilePath `json:"imageUrl"`
	Title    string             `json:"title"`
	Artist   *ArtistSearch      `json:"artist"`
	Album    *AlbumSearch       `json:"album"`
	SearchBase
}

type PlaylistSearch struct {
	ImageUrl *internal.FilePath `json:"imageUrl"`
	Title    string             `json:"title"`
	SearchBase
}
