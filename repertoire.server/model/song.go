package model

import (
	"repertoire/server/internal"
	"repertoire/server/internal/enums"
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"
)

type Song struct {
	ID              uuid.UUID          `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Title           string             `gorm:"size:100; not null" json:"title"`
	Description     string             `gorm:"not null" json:"description"`
	IsRecorded      bool               `json:"isRecorded"`
	Bpm             *uint              `json:"bpm"`
	SongsterrLink   *string            `json:"songsterrLink"`
	ReleaseDate     *time.Time         `json:"releaseDate"`
	Difficulty      *enums.Difficulty  `json:"difficulty"`
	ImageURL        *internal.FilePath `json:"imageUrl"`
	AlbumTrackNo    *uint              `json:"albumTrackNo"`
	PlaylistTrackNo *uint              `json:"playlistTrackNo"`

	AlbumID        *uuid.UUID    `json:"-"`
	ArtistID       *uuid.UUID    `json:"-"`
	GuitarTuningID *uuid.UUID    `json:"-"`
	Album          *Album        `json:"album"`
	Artist         *Artist       `json:"artist"`
	GuitarTuning   *GuitarTuning `json:"guitarTuning"`
	Sections       []SongSection `json:"sections"`
	Playlists      []Playlist    `gorm:"many2many:playlist_song" json:"-"`

	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`
	UserID    uuid.UUID `gorm:"foreignKey:UserID; references:ID; notnull" json:"-"`
}

func (s *Song) AfterFind(*gorm.DB) error {
	s.ImageURL = s.ImageURL.ToNullableFullURL()
	return nil
}

type GuitarTuning struct {
	ID    uuid.UUID `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Name  string    `gorm:"size:16; not null" json:"name"`
	Order uint      `gorm:"not null" json:"-"`
	Songs []Song    `json:"-"`

	UserID uuid.UUID `gorm:"foreignKey:UserID; references:ID; notnull" json:"-"`
}

type SongSection struct {
	ID         uuid.UUID `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Name       string    `gorm:"size:30" json:"name"`
	Rehearsals uint      `json:"rehearsals"`
	Order      uint      `gorm:"not null" json:"-"`

	SongID            uuid.UUID       `gorm:"not null" json:"-"`
	SongSectionTypeID uuid.UUID       `gorm:"not null" json:"-"`
	Song              Song            `json:"-"`
	SongSectionType   SongSectionType `json:"songSectionType"`

	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`
}

type SongSectionType struct {
	ID       uuid.UUID     `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Name     string        `gorm:"size:16" json:"name"`
	Order    uint          `gorm:"not null" json:"-"`
	Sections []SongSection `json:"-"`

	UserID uuid.UUID `gorm:"foreignKey:UserID; references:ID; notnull" json:"-"`
}

var DefaultGuitarTuning = []string{
	"E Standard", "Eb Standard", "D Standard", "C# Standard", "C Standard", "B Standard", "A# Standard", "A Standard",
	"Drop D", "Drop C#", "Drop C", "Drop B", "Drop A#", "Drop A",
}
var DefaultSongSectionTypes = []string{"Intro", "Verse", "Chorus", "Interlude", "Breakdown", "Solo", "Riff", "Outro"}
