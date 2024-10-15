package models

import (
	"time"

	"github.com/google/uuid"
)

type Playlist struct {
	ID          uuid.UUID `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Title       string    `gorm:"size:100; not null" json:"title"`
	Description string    `gorm:"not null" json:"description"`
	CreatedAt   time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`
	UserID      uuid.UUID `gorm:"not null" json:"-"`
	Songs       []Song    `gorm:"many2many:playlist_song" json:"-"`
}
