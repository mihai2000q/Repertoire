package model

import (
	"gorm.io/gorm"
	"repertoire/server/internal"
	"time"

	"github.com/google/uuid"
)

type Playlist struct {
	ID          uuid.UUID          `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Title       string             `gorm:"size:100; not null" json:"title"`
	Description string             `gorm:"not null" json:"description"`
	ImageURL    *internal.FilePath `json:"imageUrl"`
	Songs       []Song             `gorm:"many2many:playlist_song" json:"songs"`

	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`
	UserID    uuid.UUID `gorm:"foreignKey:UserID; references:ID; notnull" json:"-"`
}

func (p *Playlist) BeforeSave(*gorm.DB) error {
	p.ImageURL = p.ImageURL.StripNullableURL()
	return nil
}

func (p *Playlist) AfterFind(*gorm.DB) error {
	p.ImageURL = p.ImageURL.ToNullableFullURL()
	return nil
}
