package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SongPart struct {
	ID        uuid.UUID `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Name      string    `gorm:"size:30" json:"name"`
	SongOrder uint      `gorm:"not null" json:"-"`

	Rehearsals      uint   `gorm:"not null" json:"rehearsals"`
	Confidence      uint   `gorm:"not null" json:"confidence"`
	RehearsalsScore uint64 `gorm:"not null" json:"rehearsalsScore"`
	ConfidenceScore uint   `gorm:"not null" json:"confidenceScore"`
	Progress        uint64 `gorm:"not null" json:"progress"`

	SongID       uuid.UUID  `gorm:"not null; index: idx_song_parts_song_id" json:"-"`
	InstrumentID *uuid.UUID `json:"-"`

	Song       Song        `json:"-"`
	Instrument *Instrument `json:"instrument"`

	SectionParts           []SongSectionPart     `gorm:"foreignKey:PartID; constraint:OnDelete:CASCADE" json:"-"`
	History                []SongPartHistory     `gorm:"foreignKey:PartID; constraint:OnDelete:CASCADE" json:"-"`
	ArrangementOccurrences []SongPartOccurrences `gorm:"foreignKey:PartID; constraint:OnDelete:CASCADE" json:"-"`

	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`

	songPartDerivedFields
}

type songPartDerivedFields struct {
	BandMembers []BandMember `gorm:"-" json:"bandMembers"`
}

type SongPartHistory struct {
	ID       uuid.UUID        `gorm:"primaryKey; type:uuid; <-:create"`
	Property SongPartProperty `gorm:"size:255; not null"`
	From     uint             `gorm:"not null"`
	To       uint             `gorm:"not null"`
	PartID   uuid.UUID        `gorm:"not null; index:idx_song_part_histories_part_id"`

	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create"`
}

func (s *SongPart) AfterFind(*gorm.DB) error {
	if len(s.SectionParts) == 0 {
		s.BandMembers = []BandMember{}
		return nil
	}

	seenBm := make(map[uuid.UUID]bool)
	s.BandMembers = make([]BandMember, 0)

	for _, sp := range s.SectionParts {
		if sp.BandMember != nil && !seenBm[sp.BandMember.ID] {
			seenBm[sp.BandMember.ID] = true
			sp.BandMember.ImageURL = sp.BandMember.ImageURL.ToFullURL()
			s.BandMembers = append(s.BandMembers, *sp.BandMember)
		}
	}
	return nil
}

type SongPartProperty string

const (
	ConfidenceProperty SongPartProperty = "Confidence"
	RehearsalsProperty SongPartProperty = "Rehearsals"
)

var DefaultSongPartConfidence uint = 0
