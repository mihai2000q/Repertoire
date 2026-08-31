package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SongSection struct {
	ID    uuid.UUID `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Name  string    `gorm:"size:30" json:"name"`
	Order uint      `gorm:"not null" json:"-"`

	SongID            uuid.UUID `gorm:"not null; index: idx_song_sections_song_id" json:"-"`
	SongSectionTypeID uuid.UUID `gorm:"not null" json:"-"`

	Song            Song            `json:"-"`
	SongSectionType SongSectionType `json:"songSectionType"`

	SectionParts []SongSectionPart `gorm:"foreignKey:SectionID; constraint:OnDelete:CASCADE" json:"-"`

	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`

	songSectionDerivedFields
}

type songSectionDerivedFields struct {
	Parts      []SongPart `gorm:"-" json:"parts"`
	Rehearsals float64    `gorm:"-" json:"rehearsals"`
	Confidence float64    `gorm:"-" json:"confidence"`
	Progress   float64    `gorm:"-" json:"progress"`
}

type SongSectionPart struct {
	PartID    uuid.UUID `gorm:"primaryKey; type:uuid"`
	SectionID uuid.UUID `gorm:"primaryKey; type:uuid"`
	Order     uint      `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create"`

	Part    SongPart    `gorm:"foreignKey:PartID; constraint:OnDelete:CASCADE"`
	Section SongSection `gorm:"foreignKey:SectionID; constraint:OnDelete:CASCADE"`
}

func (s *SongSection) AfterFind(*gorm.DB) error {
	if len(s.SectionParts) == 0 {
		s.Parts = []SongPart{}
		s.Rehearsals = 0
		s.Confidence = 0
		s.Progress = 0
		return nil
	}

	partsLen := float64(len(s.SectionParts))
	s.Parts = make([]SongPart, int(partsLen))
	var totalRehearsals, totalConfidence uint
	var totalProgress uint64
	for i, sp := range s.SectionParts {
		s.Parts[i] = sp.Part
		totalRehearsals += sp.Part.Rehearsals
		totalConfidence += sp.Part.Confidence
		totalProgress += sp.Part.Progress
	}

	s.Rehearsals = float64(totalRehearsals) / partsLen
	s.Confidence = float64(totalConfidence) / partsLen
	s.Progress = float64(totalProgress) / partsLen

	return nil
}
