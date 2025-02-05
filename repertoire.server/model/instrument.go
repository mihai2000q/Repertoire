package model

import (
	"github.com/google/uuid"
)

type Instrument struct {
	ID           uuid.UUID     `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Name         string        `gorm:"size:30" json:"name"`
	Order        uint          `gorm:"not null" json:"-"`
	SongSections []SongSection `gorm:"constraint:OnDelete:SET NULL" json:"-"`

	UserID uuid.UUID `gorm:"foreignKey:UserID; references:ID; notnull" json:"-"`
}

var DefaultInstruments = []string{
	"Voice", "Piano", "Keyboard", "Drums", "Electric Guitar", "Acoustic Guitar",
	"Bass", "Ukulele", "Violin", "Saxophone", "Flute", "Harp",
}
