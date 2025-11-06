package model

import (
	"time"

	"github.com/google/uuid"
)

type SongSection struct {
	ID                 uuid.UUID `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Name               string    `gorm:"size:30" json:"name"`
	Order              uint      `gorm:"not null" json:"-"`
	Occurrences        uint      `gorm:"not null" json:"occurrences"`
	PartialOccurrences uint      `gorm:"not null" json:"partialOccurrences"`

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
