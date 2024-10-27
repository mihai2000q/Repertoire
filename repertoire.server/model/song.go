package model

import (
	"github.com/google/uuid"
	"repertoire/utils"
)

type Song struct {
	ID             uuid.UUID    `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Title          string       `gorm:"size:100; not null" json:"title"`
	IsRecorded     bool         `json:"isRecorded"`
	Reharsals      uint         `json:"reharsals"`
	Bpm            *uint        `json:"bpm"`
	SongsterrLink  *string      `json:"songsterrLink"`
	CreatedAt      time.Time    `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt      time.Time    `gorm:"default:current_timestamp; not null" json:"updatedAt"`
	UserID         uuid.UUID    `gorm:"not null" json:"-"`
	AlbumID        *uuid.UUID   `json:"-"`
	Album          Album        `json:"-"`
	ArtistID       *uuid.UUID   `json:"-"`
	Artist         Artist       `json:"-"`
	GuitarTuningID *uuid.UUID   `json:"-"`
	GuitarTuning   GuitarTuning `json:"-"`
	Playlist       []Playlist   `gorm:"many2many:playlist_song" json:"-"`
	utils.BaseUserModel
type GuitarTuning struct {
	ID        uuid.UUID `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Name      string    `gorm:"size:30; not null" json:"name"`
	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`
	UserID    uuid.UUID `gorm:"not null" json:"-"`
	Songs     []Song    `json:"-"`
}
}
