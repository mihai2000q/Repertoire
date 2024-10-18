package models

import (
	"time"

	"github.com/google/uuid"
)

type Album struct {
	ID        uuid.UUID  `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Title     string     `gorm:"size:100; not null" json:"title"`
	CreatedAt time.Time  `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"default:current_timestamp; not null" json:"updatedAt"`
	UserID    uuid.UUID  `gorm:"foreignKey:UserID; references:ID; not null" json:"-"`
	ArtistID  *uuid.UUID `json:"-"`
	Artist    Artist     `json:"-"`
	Songs     []Song     `json:"-"`
}
