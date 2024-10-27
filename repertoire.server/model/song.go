package model

import (
	"time"

	"github.com/google/uuid"
)

type Song struct {
	ID            uuid.UUID  `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Title         string     `gorm:"size:100; not null" json:"title"`
	IsRecorded    bool       `json:"isRecorded"`
	Reharsals     uint       `json:"reharsals"`
	Bpm           *uint      `json:"bpm"`
	SongsterrLink *string    `json:"songsterrLink"`
	CreatedAt     time.Time  `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt     time.Time  `gorm:"default:current_timestamp; not null" json:"updatedAt"`
	UserID        uuid.UUID  `gorm:"not null" json:"-"`
	AlbumID       *uuid.UUID `json:"-"`
	Album         Album      `json:"-"`
	ArtistID      *uuid.UUID `json:"-"`
	Artist        Artist     `json:"-"`
	Playlist      []Playlist `gorm:"many2many:playlist_song" json:"-"`
}
