package model

import (
	"gorm.io/gorm"
	"repertoire/server/internal"
	"time"

	"github.com/google/uuid"
)

type Artist struct {
	ID       uuid.UUID          `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Name     string             `gorm:"size:100; not null" json:"name"`
	ImageURL *internal.FilePath `json:"imageUrl"`
	Albums   []Album            `json:"albums"`
	Songs    []Song             `json:"songs"`

	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`
	UserID    uuid.UUID `gorm:"foreignKey:UserID; references:ID; notnull" json:"-"`
}

func (a *Artist) BeforeSave(*gorm.DB) error {
	a.ImageURL = a.ImageURL.StripNullableURL()
	return nil
}

func (a *Artist) AfterFind(*gorm.DB) error {
	a.ImageURL = a.ImageURL.ToNullableFullURL()
	return nil
}
