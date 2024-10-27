package model

import (
	"github.com/google/uuid"
	"time"
)

type Song struct {
	ID        uuid.UUID `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`
	UserID    uuid.UUID `gorm:"foreignKey:UserID; references:ID; notnull" json:"-"`

	Title         string  `gorm:"size:100; not null" json:"title"`
	Description   string  `gorm:"not null" json:"description"`
	IsRecorded    bool    `json:"isRecorded"`
	Rehearsals    uint    `json:"rehearsals"`
	Bpm           *uint   `json:"bpm"`
	SongsterrLink *string `json:"songsterrLink"`

	AlbumID        *uuid.UUID    `json:"-"`
	ArtistID       *uuid.UUID    `json:"-"`
	GuitarTuningID *uuid.UUID    `json:"-"`
	Album          Album         `json:"-"`
	Artist         Artist        `json:"-"`
	GuitarTuning   GuitarTuning  `json:"-"`
	Sections       []SongSection `json:"-"`
	Playlist       []Playlist    `gorm:"many2many:playlist_song" json:"-"`
}

type GuitarTuning struct {
	ID        uuid.UUID `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`
	UserID    uuid.UUID `gorm:"foreignKey:UserID; references:ID; notnull" json:"-"`

	Name  string `gorm:"size:16; not null" json:"name"`
	Songs []Song `json:"-"`
}

type SongSection struct {
	ID        uuid.UUID `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`

	Name              string          `gorm:"size:30" json:"name"`
	SongID            uuid.UUID       `gorm:"not null" json:"-"`
	SongSectionTypeID uuid.UUID       `gorm:"not null" json:"-"`
	Song              Song            `json:"-"`
	SongSectionType   SongSectionType `json:"-"`
}

type SongSectionType struct {
	ID        uuid.UUID `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`
	UserID    uuid.UUID `gorm:"foreignKey:UserID; references:ID; notnull" json:"-"`

	Name     string        `gorm:"size:16" json:"name"`
	Sections []SongSection `json:"-"`
}
