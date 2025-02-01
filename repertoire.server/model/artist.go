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
	IsBand   bool               `gorm:"not null" json:"isBand"`
	ImageURL *internal.FilePath `json:"imageUrl"`
	Albums   []Album            `gorm:"constraint:OnDelete:SET NULL" json:"albums"`
	Songs    []Song             `gorm:"constraint:OnDelete:SET NULL" json:"songs"`

	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`
	UserID    uuid.UUID `gorm:"foreignKey:UserID; references:ID; notnull" json:"-"`
}

func (a *Artist) BeforeSave(*gorm.DB) error {
	a.ImageURL = a.ImageURL.StripURL()
	return nil
}

func (a *Artist) AfterFind(*gorm.DB) error {
	a.ImageURL = a.ImageURL.ToFullURL(&a.UpdatedAt)
	return nil
}
