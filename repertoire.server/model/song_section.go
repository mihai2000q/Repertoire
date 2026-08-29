package model

import (
	"time"

	"github.com/google/uuid"
)

type SongSection struct {
	ID    uuid.UUID `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Name  string    `gorm:"size:30" json:"name"`
	Order uint      `gorm:"not null" json:"-"`

	Rehearsals float64 `gorm:"-" json:"rehearsals"`
	Confidence float64 `gorm:"-" json:"confidence"`
	Progress   float64 `gorm:"-" json:"progress"`

	SongID            uuid.UUID `gorm:"not null; index: idx_song_sections_song_id" json:"-"`
	SongSectionTypeID uuid.UUID `gorm:"not null" json:"-"`

	Song            Song            `json:"-"`
	SongSectionType SongSectionType `json:"songSectionType"`

	Parts []SongPart `gorm:"foreignKey:SectionID; constraint:OnDelete:CASCADE" json:"-"`

	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`
}
