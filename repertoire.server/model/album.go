package model

import (
	"gorm.io/gorm"
	"repertoire/server/internal"
	"time"

	"github.com/google/uuid"
)

type Album struct {
	ID          uuid.UUID          `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Title       string             `gorm:"size:100; not null" json:"title"`
	ReleaseDate *time.Time         `json:"releaseDate"`
	ImageURL    *internal.FilePath `json:"imageUrl"`
	ArtistID    *uuid.UUID         `json:"-"`
	Artist      *Artist            `json:"artist"`
	Songs       []Song             `gorm:"constraint:OnDelete:SET NULL" json:"songs"`

	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`
	UserID    uuid.UUID `gorm:"foreignKey:UserID; references:ID; notnull" json:"-"`
}

func (a *Album) BeforeSave(*gorm.DB) error {
	a.ImageURL = a.ImageURL.StripNullableURL()
	return nil
}

func (a *Album) AfterFind(*gorm.DB) error {
	a.ImageURL = a.ImageURL.ToNullableFullURL()
	// When Joins instead of Preload, AfterFind Hook is not used
	if a.Artist != nil {
		a.Artist.ImageURL = a.Artist.ImageURL.ToNullableFullURL()
	}
	return nil
}
