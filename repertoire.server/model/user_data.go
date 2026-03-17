package model

import "github.com/google/uuid"

type BandMemberRole struct {
	ID          uuid.UUID    `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Name        string       `gorm:"size:24; not null" json:"name"`
	Order       uint         `gorm:"not null" json:"-"`
	BandMembers []BandMember `gorm:"many2many:band_member_has_roles" json:"-"`

	UserID uuid.UUID `gorm:"foreignKey:UserID; references:ID; not null; <-:create; index:idx_band_member_roles_user_id" json:"-"`
}

var DefaultBandMemberRoles = []string{
	"Vocalist", "Lead Guitarist", "Rhythm Guitarist", "Bassist", "Drummer", "Pianist", "Keyboardist", "Backing Vocalist",
}

type GuitarTuning struct {
	ID    uuid.UUID `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Name  string    `gorm:"size:16; not null" json:"name"`
	Order uint      `gorm:"not null" json:"-"`
	Songs []Song    `gorm:"constraint:OnDelete:SET NULL" json:"-"`

	UserID uuid.UUID `gorm:"foreignKey:UserID; references:ID; not null; <-:create; index:idx_guitar_tunings_user_id" json:"-"`
}

var DefaultGuitarTunings = []string{
	"E Standard", "Eb Standard", "D Standard", "C# Standard", "C Standard", "B Standard", "A# Standard", "A Standard",
	"Drop D", "Drop C#", "Drop C", "Drop B", "Drop A#", "Drop A",
}

type Instrument struct {
	ID           uuid.UUID      `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Name         string         `gorm:"size:30" json:"name"`
	Order        uint           `gorm:"not null" json:"-"`
	SongSections []SongSection  `gorm:"constraint:OnDelete:SET NULL" json:"-"`
	SongSettings []SongSettings `gorm:"foreignKey:DefaultInstrumentID; references:ID; constraint:OnDelete:SET NULL" json:"-"`

	UserID uuid.UUID `gorm:"foreignKey:UserID; references:ID; not null; <-:create; index:idx_instruments_user_id" json:"-"`
}

var DefaultInstruments = []string{
	"Voice", "Piano", "Keyboard", "Drums", "Electric Guitar", "Acoustic Guitar",
	"Bass", "Ukulele", "Violin", "Saxophone", "Flute", "Harp",
}

type SongSectionType struct {
	ID       uuid.UUID     `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Name     string        `gorm:"size:16" json:"name"`
	Order    uint          `gorm:"not null" json:"-"`
	Sections []SongSection `gorm:"constraint:OnDelete:CASCADE" json:"-"`

	UserID uuid.UUID `gorm:"foreignKey:UserID; references:ID; not null; <-:create; index:idx_song_section_types_user_id" json:"-"`
}

var DefaultSongSectionTypes = []string{
	"Intro", "Verse", "Pre-Chorus", "Chorus", "Interlude",
	"Bridge", "Breakdown", "Solo", "Riff", "Outro",
}
