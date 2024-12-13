package model

import (
	"math"
	"repertoire/server/internal"
	"repertoire/server/internal/enums"
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"
)

type Song struct {
	ID             uuid.UUID          `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Title          string             `gorm:"size:100; not null" json:"title"`
	Description    string             `gorm:"not null" json:"description"`
	ReleaseDate    *time.Time         `json:"releaseDate"`
	ImageURL       *internal.FilePath `json:"imageUrl"`
	IsRecorded     bool               `json:"isRecorded"`
	Bpm            *uint              `json:"bpm"`
	Difficulty     *enums.Difficulty  `json:"difficulty"`
	SongsterrLink  *string            `json:"songsterrLink"`
	YoutubeLink    *string            `json:"youtubeLink"`
	AlbumTrackNo   *uint              `json:"albumTrackNo"`
	LastTimePlayed *time.Time         `json:"lastTimePlayed"`

	Rehearsals *uint   `gorm:"-" json:"rehearsals"`
	Confidence *uint   `gorm:"-" json:"confidence"`
	Progress   *uint64 `gorm:"-" json:"progress"`

	AlbumID        *uuid.UUID     `json:"-"`
	ArtistID       *uuid.UUID     `json:"-"`
	GuitarTuningID *uuid.UUID     `json:"-"`
	Artist         *Artist        `json:"artist"`
	Album          *Album         `json:"album"`
	GuitarTuning   *GuitarTuning  `json:"guitarTuning"`
	Sections       []SongSection  `gorm:"constraint:OnDelete:CASCADE" json:"sections"`
	Playlists      []Playlist     `gorm:"many2many:playlist_songs" json:"playlists"`
	PlaylistSongs  []PlaylistSong `gorm:"foreignKey:SongID; constraint:OnDelete:CASCADE" json:"-"`

	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`
	UserID    uuid.UUID `gorm:"foreignKey:UserID; references:ID; notnull" json:"-"`
	playlistSongMetadata
}

type playlistSongMetadata struct {
	PlaylistTrackNo   uint      `gorm:"-" json:"playlistTrackNo"`
	PlaylistCreatedAt time.Time `gorm:"-" json:"playlistCreatedAt"`
}

func (s *Song) BeforeSave(*gorm.DB) error {
	s.ImageURL = s.ImageURL.StripNullableURL()
	return nil
}

func (s *Song) AfterFind(*gorm.DB) error {
	s.ImageURL = s.ImageURL.ToNullableFullURL()

	if len(s.Sections) == 0 {
		return nil
	}

	rehearsals := uint(0)
	confidence := uint(0)
	progress := uint64(0)
	for _, section := range s.Sections {
		rehearsals += section.Rehearsals
		confidence += section.Confidence
		progress += section.Progress
	}
	rehearsals = uint(math.Round(float64(int(rehearsals) / len(s.Sections))))
	confidence = uint(math.Round(float64(int(confidence) / len(s.Sections))))
	progress = uint64(math.Round(float64(int(progress) / len(s.Sections))))

	s.Rehearsals = &rehearsals
	s.Confidence = &confidence
	s.Progress = &progress

	return nil
}

// Guitar Tuning

type GuitarTuning struct {
	ID    uuid.UUID `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Name  string    `gorm:"size:16; not null" json:"name"`
	Order uint      `gorm:"not null" json:"-"`
	Songs []Song    `gorm:"constraint:OnDelete:SET NULL" json:"-"`

	UserID uuid.UUID `gorm:"foreignKey:UserID; references:ID; notnull" json:"-"`
}

var DefaultGuitarTunings = []string{
	"E Standard", "Eb Standard", "D Standard", "C# Standard", "C Standard", "B Standard", "A# Standard", "A Standard",
	"Drop D", "Drop C#", "Drop C", "Drop B", "Drop A#", "Drop A",
}

// Song Sections

type SongSection struct {
	ID    uuid.UUID `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Name  string    `gorm:"size:30" json:"name"`
	Order uint      `gorm:"not null" json:"-"`

	Rehearsals      uint   `gorm:"not null" json:"rehearsals"`
	Confidence      uint   `gorm:"not null; size:100" json:"confidence"`
	RehearsalsScore uint64 `gorm:"not null" json:"rehearsalsScore"`
	ConfidenceScore uint   `gorm:"not null" json:"confidenceScore"`
	Progress        uint64 `gorm:"not null" json:"progress"`

	SongID            uuid.UUID       `gorm:"not null" json:"-"`
	SongSectionTypeID uuid.UUID       `gorm:"not null" json:"-"`
	Song              Song            `json:"-"`
	SongSectionType   SongSectionType `json:"songSectionType"`

	History []SongSectionHistory `gorm:"constraint:OnDelete:CASCADE" json:"-"`

	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`
}

type SongSectionType struct {
	ID       uuid.UUID     `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Name     string        `gorm:"size:16" json:"name"`
	Order    uint          `gorm:"not null" json:"-"`
	Sections []SongSection `gorm:"constraint:OnDelete:CASCADE" json:"-"`

	UserID uuid.UUID `gorm:"foreignKey:UserID; references:ID; notnull" json:"-"`
}

type SongSectionHistory struct {
	ID            uuid.UUID           `gorm:"primaryKey; type:uuid; <-:create"`
	Property      SongSectionProperty `gorm:"size:255; not null"`
	From          uint                `gorm:"not null"`
	To            uint                `gorm:"not null"`
	SongSectionID uuid.UUID           `gorm:"not null"`

	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create"`
}

type SongSectionProperty string

const (
	ConfidenceProperty SongSectionProperty = "Confidence"
	RehearsalsProperty SongSectionProperty = "Rehearsals"
)

var DefaultSongSectionConfidence uint = 0

var DefaultSongSectionTypes = []string{"Intro", "Verse", "Chorus", "Interlude", "Breakdown", "Solo", "Riff", "Outro"}
