package model

import (
	"github.com/google/uuid"
	"repertoire/server/internal"
	"repertoire/server/internal/enums"
	"time"
)

type SearchResult []string

type SearchBase struct {
	ID     string           `json:"id"`
	Type   enums.SearchType `json:"type"`
	UserID uuid.UUID        `json:"userId"`
}

// Artist

type ArtistSearch struct {
	ImageUrl  *internal.FilePath `json:"imageUrl"`
	Name      string             `json:"name"`
	UpdatedAt time.Time          `json:"updatedAt"`
	SearchBase
}

// Album

type AlbumSearch struct {
	ImageUrl  *internal.FilePath `json:"imageUrl"`
	Title     string             `json:"title"`
	UpdatedAt time.Time          `json:"updatedAt"`
	Artist    *AlbumArtistSearch `json:"artist"`
	SearchBase
}

type AlbumArtistSearch struct {
	ID        uuid.UUID          `json:"id"`
	ImageUrl  *internal.FilePath `json:"imageUrl"`
	Name      string             `json:"name"`
	UpdatedAt time.Time          `json:"updatedAt"`
}

// Song

type SongSearch struct {
	ImageUrl  *internal.FilePath `json:"imageUrl"`
	Title     string             `json:"title"`
	UpdatedAt time.Time          `json:"updatedAt"`
	Artist    *SongArtistSearch  `json:"artist"`
	Album     *SongAlbumSearch   `json:"album"`
	SearchBase
}

type SongAlbumSearch struct {
	ID        uuid.UUID          `json:"id"`
	ImageUrl  *internal.FilePath `json:"imageUrl"`
	Title     string             `json:"title"`
	UpdatedAt time.Time          `json:"updatedAt"`
}

type SongArtistSearch struct {
	ID        uuid.UUID          `json:"id"`
	ImageUrl  *internal.FilePath `json:"imageUrl"`
	Name      string             `json:"name"`
	UpdatedAt time.Time          `json:"updatedAt"`
}

// Playlist

type PlaylistSearch struct {
	ImageUrl  *internal.FilePath `json:"imageUrl"`
	Title     string             `json:"title"`
	UpdatedAt time.Time          `json:"updatedAt"`
	SearchBase
}
