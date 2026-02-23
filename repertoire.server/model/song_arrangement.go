package model

import "github.com/google/uuid"

type SongArrangement struct {
	ID                 uuid.UUID                `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Name               string                   `gorm:"size:30; not null" json:"name"`
	Order              uint                     `gorm:"not null" json:"-"`
	SongID             uuid.UUID                `gorm:"foreignKey:SongID; references:ID; not null" json:"-"`
	SectionOccurrences []SongSectionOccurrences `gorm:"foreignKey:ArrangementID; constraint:OnDelete:CASCADE" json:"sectionOccurrences"`
	DefaultSong        *Song                    `gorm:"foreignKey:DefaultArrangementID; references:ID; constraint:OnDelete:SET NULL" json:"-"`
}

type SongSectionOccurrences struct {
	Occurrences   uint        `gorm:"not null" json:"occurrences"`
	Section       SongSection `gorm:"not null" json:"section"`
	SectionID     uuid.UUID   `gorm:"primaryKey; type:uuid; <-:create;" json:"-"`
	ArrangementID uuid.UUID   `gorm:"primaryKey; type:uuid; <-:create;" json:"-"`
}

var DefaultSongArrangementName = "Perfect Rehearsal"
