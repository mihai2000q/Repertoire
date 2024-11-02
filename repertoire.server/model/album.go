package model

import (
	"time"

	"github.com/google/uuid"
)

type Album struct {
	ID          uuid.UUID  `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Title       string     `gorm:"size:100; not null" json:"title"`
	ReleaseDate *time.Time `json:"releaseDate"`
	ArtistID    *uuid.UUID `json:"-"`
	Artist      *Artist    `json:"artist"`
	Songs       []Song     `json:"songs"`

	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`
	UserID    uuid.UUID `gorm:"foreignKey:UserID; references:ID; notnull" json:"-"`
}
