package model

import (
	"gorm.io/gorm"
	"repertoire/server/internal"
	"time"

	"github.com/google/uuid"
)

type EnhancedPlaylist struct {
	Playlist
	SongsCount float64 `gorm:"->" json:"songsCount"`
}

type Playlist struct {
	ID            uuid.UUID          `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Title         string             `gorm:"size:100; not null" json:"title"`
	Description   string             `gorm:"not null" json:"description"`
	ImageURL      *internal.FilePath `json:"imageUrl"`
	Songs         []Song             `gorm:"many2many:playlist_songs" json:"songs"`
	PlaylistSongs []PlaylistSong     `gorm:"foreignKey:PlaylistID; constraint:OnDelete:CASCADE" json:"-"`

	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`
	UserID    uuid.UUID `gorm:"foreignKey:UserID; references:ID; notnull" json:"userId"`
}

type PlaylistSong struct {
	ID          uuid.UUID `gorm:"primaryKey; type:uuid; <-:create"`
	PlaylistID  uuid.UUID `gorm:"type:uuid; not null; <-:create"`
	SongID      uuid.UUID `gorm:"type:uuid; not null; <-:create"`
	SongTrackNo uint      `gorm:"not null"`
	CreatedAt   time.Time `gorm:"default:current_timestamp; not null; <-:create"`

	Playlist Playlist
	Song     Song
}

func (p *Playlist) BeforeSave(*gorm.DB) error {
	p.ImageURL = p.ImageURL.StripURL()
	return nil
}

func (p *Playlist) AfterFind(*gorm.DB) error {
	p.ImageURL = p.ImageURL.ToFullURL()
	return nil
}
