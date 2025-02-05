package model

import (
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
	Rehearsals     float64            `gorm:"not null" json:"rehearsals"`
	Confidence     float64            `gorm:"not null" json:"confidence"`
	Progress       float64            `gorm:"not null" json:"progress"`

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
	s.ImageURL = s.ImageURL.StripURL()
	return nil
}

func (s *Song) AfterFind(*gorm.DB) error {
	s.ImageURL = s.ImageURL.ToFullURL(&s.UpdatedAt)
	// When Joins instead of Preload, AfterFind Hook is not used
	if s.Artist != nil {
		s.Artist.ImageURL = s.Artist.ImageURL.ToFullURL(&s.Artist.UpdatedAt)
	}
	if s.Album != nil {
		s.Album.ImageURL = s.Album.ImageURL.ToFullURL(&s.Album.UpdatedAt)
	}

	return nil
}

// Song Sections

type SongSection struct {
	ID          uuid.UUID `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Name        string    `gorm:"size:30" json:"name"`
	Order       uint      `gorm:"not null" json:"-"`
	Occurrences uint      `gorm:"not null" json:"occurrences"`

	Rehearsals      uint   `gorm:"not null" json:"rehearsals"`
	Confidence      uint   `gorm:"not null; size:100" json:"confidence"`
	RehearsalsScore uint64 `gorm:"not null" json:"rehearsalsScore"`
	ConfidenceScore uint   `gorm:"not null" json:"confidenceScore"`
	Progress        uint64 `gorm:"not null" json:"progress"`

	SongID            uuid.UUID  `gorm:"not null" json:"-"`
	SongSectionTypeID uuid.UUID  `gorm:"not null" json:"-"`
	InstrumentID      *uuid.UUID `json:"-"`
	BandMemberID      *uuid.UUID `json:"-"`

	Song            Song            `json:"-"`
	SongSectionType SongSectionType `json:"songSectionType"`
	BandMember      *BandMember     `json:"bandMember"`
	Instrument      *Instrument     `json:"instrument"`

	History []SongSectionHistory `gorm:"constraint:OnDelete:CASCADE" json:"-"`

	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`
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
