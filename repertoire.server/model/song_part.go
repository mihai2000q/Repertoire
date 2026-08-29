package model

import (
	"time"

	"github.com/google/uuid"
)

type SongPart struct {
	ID           uuid.UUID `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Name         string    `gorm:"size:30" json:"name"`
	SongOrder    uint      `gorm:"not null" json:"-"`
	SectionOrder *uint     `json:"-"`

	Rehearsals      uint   `gorm:"not null" json:"rehearsals"`
	Confidence      uint   `gorm:"not null" json:"confidence"`
	RehearsalsScore uint64 `gorm:"not null" json:"rehearsalsScore"`
	ConfidenceScore uint   `gorm:"not null" json:"confidenceScore"`
	Progress        uint64 `gorm:"not null" json:"progress"`

	SongID       uuid.UUID  `gorm:"not null; index: idx_song_instrument_parts_song_id" json:"-"`
	SectionID    *uuid.UUID `gorm:"not null; index: idx_song_instrument_parts_song_section_id" json:"-"`
	BandMemberID *uuid.UUID `json:"-"`
	InstrumentID *uuid.UUID `json:"-"`

	Section    SongSection `json:"-"`
	BandMember *BandMember `json:"bandMember"`
	Instrument *Instrument `json:"instrument"`

	History                []SongPartHistory     `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	ArrangementOccurrences []SongPartOccurrences `gorm:"foreignKey:PartID; constraint:OnDelete:CASCADE" json:"-"`

	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`
}

type SongPartHistory struct {
	ID         uuid.UUID        `gorm:"primaryKey; type:uuid; <-:create"`
	Property   SongPartProperty `gorm:"size:255; not null"`
	From       uint             `gorm:"not null"`
	To         uint             `gorm:"not null"`
	SongPartID uuid.UUID        `gorm:"not null; index:idx_song_instrument_part_histories_song_instrument_part_id"`

	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create"`
}

type SongPartProperty string

const (
	ConfidenceProperty SongPartProperty = "Confidence"
	RehearsalsProperty SongPartProperty = "Rehearsals"
)

var DefaultSongPartConfidence uint = 0
