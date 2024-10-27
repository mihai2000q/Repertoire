package model

import (
	"github.com/google/uuid"
	"time"
)

type Song struct {
	ID            uuid.UUID `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Title         string    `gorm:"size:100; not null" json:"title"`
	Description   string    `gorm:"not null" json:"description"`
	IsRecorded    bool      `json:"isRecorded"`
	Bpm           *uint     `json:"bpm"`
	SongsterrLink *string   `json:"songsterrLink"`

	AlbumID        *uuid.UUID    `json:"-"`
	ArtistID       *uuid.UUID    `json:"-"`
	GuitarTuningID *uuid.UUID    `json:"-"`
	Album          Album         `json:"-"`
	Artist         Artist        `json:"-"`
	GuitarTuning   GuitarTuning  `json:"-"`
	Sections       []SongSection `json:"-"`
	Playlist       []Playlist    `gorm:"many2many:playlist_song" json:"-"`

	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`
	UserID    uuid.UUID `gorm:"foreignKey:UserID; references:ID; notnull" json:"-"`
}

type GuitarTuning struct {
	ID    uuid.UUID `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Name  string    `gorm:"size:16; not null" json:"name"`
	Order uint      `gorm:"not null" json:"-"`
	Songs []Song    `json:"-"`

	UserID uuid.UUID `gorm:"foreignKey:UserID; references:ID; notnull" json:"-"`
}

type SongSection struct {
	ID         uuid.UUID `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Name       string    `gorm:"size:30" json:"name"`
	Rehearsals uint      `json:"rehearsals"`

	SongID            uuid.UUID       `gorm:"not null" json:"-"`
	SongSectionTypeID uuid.UUID       `gorm:"not null" json:"-"`
	Song              Song            `json:"-"`
	SongSectionType   SongSectionType `json:"-"`

	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`
}

type SongSectionType struct {
	ID       uuid.UUID     `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Name     string        `gorm:"size:16" json:"name"`
	Order    uint          `gorm:"not null" json:"-"`
	Sections []SongSection `json:"-"`

	UserID uuid.UUID `gorm:"foreignKey:UserID; references:ID; notnull" json:"-"`
}

var DefaultGuitarTuning = []string{
	"E Standard", "Eb Standard", "D Standard", "C# Standard", "C Standard", "B Standard", "A# Standard", "A Standard",
	"Drop D", "Drop C#", "Drop C", "Drop B", "Drop A#", "Drop A",
}
var DefaultSongSectionTypes = []string{"Intro", "Verse", "Chorus", "Interlude", "Breakdown", "Solo", "Riff", "Outro"}
