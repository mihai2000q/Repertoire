package model

import (
	"repertoire/server/internal"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                uuid.UUID          `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Email             string             `gorm:"size:256; unique; not null" json:"email"`
	Password          string             `gorm:"not null" json:"-"`
	Name              string             `gorm:"size:100; not null" json:"name"`
	ProfilePictureURL *internal.FilePath `json:"profilePictureUrl"`

	Albums           []Album           `json:"-"`
	Artists          []Artist          `json:"-"`
	Playlists        []Playlist        `json:"-"`
	Songs            []Song            `json:"-"`
	SongSectionTypes []SongSectionType `json:"-"`
	GuitarTunings    []GuitarTuning    `json:"-"`

	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`
}
