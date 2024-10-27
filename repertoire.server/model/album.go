package model

import (
	"github.com/google/uuid"
	"time"
)

type Album struct {
	ID        uuid.UUID `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`
	UserID    uuid.UUID `gorm:"foreignKey:UserID; references:ID; notnull" json:"-"`

	Title    string     `gorm:"size:100; not null" json:"title"`
	ArtistID *uuid.UUID `json:"-"`
	Artist   Artist     `json:"-"`
	Songs    []Song     `json:"-"`
}
