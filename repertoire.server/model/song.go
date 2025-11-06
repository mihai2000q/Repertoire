package model

import (
	"repertoire/server/internal"
	"repertoire/server/internal/enums"
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"
)

type EnhancedSong struct {
	Song
	SectionsCount float64 `gorm:"->" json:"sectionsCount"`
	SolosCount    float64 `gorm:"->" json:"solosCount"`
	RiffsCount    float64 `gorm:"->" json:"riffsCount"`
}

type Song struct {
	ID             uuid.UUID          `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Title          string             `gorm:"size:100; not null" json:"title"`
	Description    string             `gorm:"not null" json:"description"`
	ReleaseDate    *internal.Date     `json:"releaseDate"`
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

	Settings       SongSettings   `gorm:"constraint:OnDelete:CASCADE" json:"settings"`
	AlbumID        *uuid.UUID     `json:"albumId"`
	ArtistID       *uuid.UUID     `json:"artistId"`
	GuitarTuningID *uuid.UUID     `json:"-"`
	Artist         *Artist        `json:"artist"`
	Album          *Album         `json:"album"`
	GuitarTuning   *GuitarTuning  `json:"guitarTuning"`
	Sections       []SongSection  `gorm:"constraint:OnDelete:CASCADE" json:"sections"`
	Playlists      []Playlist     `gorm:"many2many:playlist_songs" json:"playlists"`
	PlaylistSongs  []PlaylistSong `gorm:"foreignKey:SongID; constraint:OnDelete:CASCADE" json:"-"`

	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`
	UserID    uuid.UUID `gorm:"foreignKey:UserID; references:ID; notnull" json:"userId"`
	playlistSongMetadata
}

type playlistSongMetadata struct {
	PlaylistSongID    uuid.UUID `gorm:"-" json:"playlistSongId"`
	PlaylistTrackNo   uint      `gorm:"-" json:"playlistTrackNo"`
	PlaylistCreatedAt time.Time `gorm:"-" json:"playlistCreatedAt"`
}

func (s *Song) BeforeSave(*gorm.DB) error {
	s.ImageURL = s.ImageURL.StripURL()
	return nil
}

func (s *Song) AfterFind(*gorm.DB) error {
	s.ToFullImageURL()
	return nil
}

func (s *Song) ToFullImageURL() {
	s.ImageURL = s.ImageURL.ToFullURL()
	// When Joins instead of Preload, AfterFind Hook is not used
	if s.Artist != nil {
		s.Artist.ImageURL = s.Artist.ImageURL.ToFullURL()
	}
	if s.Album != nil {
		s.Album.ImageURL = s.Album.ImageURL.ToFullURL()
	}
	if s.Settings.DefaultBandMember != nil {
		s.Settings.DefaultBandMember.ImageURL =
			s.Settings.DefaultBandMember.ImageURL.ToFullURL()
	}
}

// Song Settings

type SongSettings struct {
	ID                  uuid.UUID   `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	DefaultInstrumentID *uuid.UUID  `json:"-"`
	DefaultInstrument   *Instrument `json:"defaultInstrument"`
	DefaultBandMemberID *uuid.UUID  `json:"-"`
	DefaultBandMember   *BandMember `json:"defaultBandMember"`

	SongID uuid.UUID `gorm:"not null" json:"-"`
}
