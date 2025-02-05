package model

import (
	"gorm.io/gorm"
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
	InstrumentTypes  []InstrumentType  `json:"-"`
	GuitarTunings    []GuitarTuning    `json:"-"`
	BandMemberRoles  []BandMemberRole  `json:"-"`

	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`
}

func (u *User) BeforeSave(*gorm.DB) error {
	u.ProfilePictureURL = u.ProfilePictureURL.StripURL()
	return nil
}

func (u *User) AfterFind(*gorm.DB) error {
	u.ProfilePictureURL = u.ProfilePictureURL.ToFullURL(&u.UpdatedAt)
	return nil
}
