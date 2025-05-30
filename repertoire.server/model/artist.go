package model

import (
	"gorm.io/gorm"
	"repertoire/server/internal"
	"time"

	"github.com/google/uuid"
)

type EnhancedArtist struct {
	Artist
	BandMembersCount float64    `gorm:"->" json:"bandMembersCount"`
	AlbumsCount      float64    `gorm:"->" json:"albumsCount"`
	SongsCount       float64    `gorm:"->" json:"songsCount"`
	Rehearsals       float64    `gorm:"->" json:"rehearsals"`
	Confidence       float64    `gorm:"->" json:"confidence"`
	Progress         float64    `gorm:"->" json:"progress"`
	LastTimePlayed   *time.Time `gorm:"->" json:"lastTimePlayed"`
}

type Artist struct {
	ID       uuid.UUID          `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Name     string             `gorm:"size:100; not null" json:"name"`
	IsBand   bool               `gorm:"not null" json:"isBand"`
	ImageURL *internal.FilePath `json:"imageUrl"`

	Albums      []Album      `gorm:"constraint:OnDelete:SET NULL" json:"albums"`
	Songs       []Song       `gorm:"constraint:OnDelete:SET NULL" json:"songs"`
	BandMembers []BandMember `gorm:"constraint:OnDelete:CASCADE" json:"bandMembers"`

	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`
	UserID    uuid.UUID `gorm:"foreignKey:UserID; references:ID; notnull" json:"userId"`
}

func (a *Artist) BeforeSave(*gorm.DB) error {
	a.ImageURL = a.ImageURL.StripURL()
	return nil
}

func (a *Artist) AfterFind(*gorm.DB) error {
	a.ImageURL = a.ImageURL.ToFullURL()
	return nil
}

// Band Member

type BandMember struct {
	ID       uuid.UUID          `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Name     string             `gorm:"size:100; not null" json:"name"`
	Order    uint               `gorm:"not null" json:"-"`
	Color    *string            `gorm:"size:7" json:"color"`
	ImageURL *internal.FilePath `json:"imageUrl"`

	ArtistID     uuid.UUID        `gorm:"not null" json:"-"`
	Artist       Artist           `json:"-"`
	Roles        []BandMemberRole `gorm:"many2many:band_member_has_roles" json:"roles"`
	SongSections []SongSection    `gorm:"constraint:OnDelete:SET NULL" json:"-"`
	SongSettings []SongSettings   `gorm:"foreignKey:DefaultBandMemberID; references:ID; constraint:OnDelete:SET NULL" json:"-"`

	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`
}

func (b *BandMember) BeforeSave(*gorm.DB) error {
	b.ImageURL = b.ImageURL.StripURL()
	return nil
}

func (b *BandMember) AfterFind(*gorm.DB) error {
	b.ImageURL = b.ImageURL.ToFullURL()
	return nil
}
