package model

import (
	"github.com/google/uuid"
	"repertoire/utils"
)

type Song struct {
	Title          string        `gorm:"size:100; not null" json:"title"`
	IsRecorded     bool          `json:"isRecorded"`
	Rehearsals     uint          `json:"rehearsals"`
	Bpm            *uint         `json:"bpm"`
	SongsterrLink  *string       `json:"songsterrLink"`
	AlbumID        *uuid.UUID    `json:"-"`
	ArtistID       *uuid.UUID    `json:"-"`
	GuitarTuningID *uuid.UUID    `json:"-"`
	Album          Album         `json:"-"`
	Artist         Artist        `json:"-"`
	GuitarTuning   GuitarTuning  `json:"-"`
	Sections       []SongSection `json:"-"`
	Playlist       []Playlist    `gorm:"many2many:playlist_song" json:"-"`
	utils.BaseUserModel
}

type GuitarTuning struct {
	Name  string `gorm:"size:30; not null" json:"name"`
	Songs []Song `json:"-"`
	utils.BaseUserModel
}

type SongSection struct {
	Name   string          `gorm:"size:50" json:"name"`
	SongID uuid.UUID       `gorm:"not null" json:"-"`
	Song   Song            `json:"-"`
	TypeId uuid.UUID       `gorm:"not null" json:"-"`
	Type   SongSectionType `json:"-"`
	utils.BaseModel
}

type SongSectionType struct {
	Name     string        `gorm:"size:16" json:"name"`
	Sections []SongSection `json:"-"`
	utils.BaseUserModel
}
