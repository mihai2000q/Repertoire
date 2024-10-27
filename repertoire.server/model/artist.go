package model

import (
	"github.com/google/uuid"
	"time"
)

type Artist struct {
	ID        uuid.UUID `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`
	UserID    uuid.UUID `gorm:"foreignKey:UserID; references:ID; notnull" json:"-"`

	Name   string  `gorm:"size:100; not null" json:"name"`
	Albums []Album `json:"-"`
	Songs  []Song  `json:"-"`
}
