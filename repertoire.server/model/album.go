package model

import (
	"gorm.io/gorm"
	"repertoire/server/internal"
	"time"

	"github.com/google/uuid"
)

type EnhancedAlbum struct {
	Album
	SongsCount     float64    `gorm:"->" json:"songsCount"`
	Rehearsals     float64    `gorm:"->" json:"rehearsals"`
	Confidence     float64    `gorm:"->" json:"confidence"`
	Progress       float64    `gorm:"->" json:"progress"`
	LastTimePlayed *time.Time `gorm:"->" json:"lastTimePlayed"`
}

type Album struct {
	ID          uuid.UUID          `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Title       string             `gorm:"size:100; not null" json:"title"`
	ReleaseDate *time.Time         `json:"releaseDate"`
	ImageURL    *internal.FilePath `json:"imageUrl"`
	ArtistID    *uuid.UUID         `json:"artistId"`
	Artist      *Artist            `json:"artist"`
	Songs       []Song             `gorm:"constraint:OnDelete:SET NULL" json:"songs"`

	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`
	UserID    uuid.UUID `gorm:"foreignKey:UserID; references:ID; notnull" json:"userId"`
}

func (a *Album) BeforeSave(*gorm.DB) error {
	a.ImageURL = a.ImageURL.StripURL()
	return nil
}

func (a *Album) AfterFind(*gorm.DB) error {
	a.ImageURL = a.ImageURL.ToFullURL()
	// When Joins instead of Preload, AfterFind Hook is not used
	if a.Artist != nil {
		a.Artist.ImageURL = a.Artist.ImageURL.ToFullURL()
	}
	return nil
}
