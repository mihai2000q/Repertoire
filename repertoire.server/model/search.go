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

func NewArtistSearch(
	id uuid.UUID,
	name string,
	imageUrl *internal.FilePath,
	updatedAt time.Time,
	userID uuid.UUID,
) ArtistSearch {
	return ArtistSearch{
		ImageUrl:  imageUrl,
		Name:      name,
		UpdatedAt: updatedAt,
		SearchBase: SearchBase{
			ID:     "artist-" + id.String(),
			Type:   enums.Artist,
			UserID: userID,
		},
	}
}

func (a *Artist) ToSearch() ArtistSearch {
	return NewArtistSearch(
		a.ID,
		a.Name,
		a.ImageURL.StripURL(),
		a.UpdatedAt,
		a.UserID,
	)
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

func (a *Album) ToSearch() AlbumSearch {
	search := AlbumSearch{
		ImageUrl:  a.ImageURL.StripURL(),
		Title:     a.Title,
		UpdatedAt: a.UpdatedAt,
		SearchBase: SearchBase{
			ID:     "album-" + a.ID.String(),
			Type:   enums.Album,
			UserID: a.UserID,
		},
	}

	if a.Artist != nil {
		search.Artist = a.Artist.ToAlbumSearch()
	}

	return search
}

func (a *Artist) ToAlbumSearch() *AlbumArtistSearch {
	if a == nil {
		return nil
	}
	return &AlbumArtistSearch{
		ID:       a.ID,
		Name:     a.Name,
		ImageUrl: a.ImageURL.StripURL(),
	}
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

func (s *Song) ToSearch() SongSearch {
	search := SongSearch{
		ImageUrl:  s.ImageURL.StripURL(),
		Title:     s.Title,
		UpdatedAt: s.UpdatedAt,
		SearchBase: SearchBase{
			ID:     "song-" + s.ID.String(),
			Type:   enums.Song,
			UserID: s.UserID,
		},
	}

	if s.Artist != nil {
		search.Artist = s.Artist.ToSongSearch()
	}

	if s.Album != nil {
		search.Album = s.Album.ToSongSearch()
	}

	return search
}

func (a *Artist) ToSongSearch() *SongArtistSearch {
	if a == nil {
		return nil
	}
	return &SongArtistSearch{
		ID:       a.ID,
		Name:     a.Name,
		ImageUrl: a.ImageURL.StripURL(),
	}
}

func (a *Album) ToSongSearch() *SongAlbumSearch {
	if a == nil {
		return nil
	}
	return &SongAlbumSearch{
		ID:       a.ID,
		Title:    a.Title,
		ImageUrl: a.ImageURL.StripURL(),
	}
}

// Playlist

type PlaylistSearch struct {
	ImageUrl  *internal.FilePath `json:"imageUrl"`
	Title     string             `json:"title"`
	UpdatedAt time.Time          `json:"updatedAt"`
	SearchBase
}

func (p *Playlist) ToSearch() PlaylistSearch {
	return PlaylistSearch{
		ImageUrl:  p.ImageURL.StripURL(),
		Title:     p.Title,
		UpdatedAt: p.UpdatedAt,
		SearchBase: SearchBase{
			ID:     "playlist-" + p.ID.String(),
			Type:   enums.Playlist,
			UserID: p.UserID,
		},
	}
}
