package models

import (
	"time"

	"github.com/google/uuid"
)

type Artist struct {
	ID        uuid.UUID `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Name      string    `gorm:"size:100; not null" json:"name"`
	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`
	UserID    uuid.UUID `gorm:"foreignKey:UserID; references:ID; not null" json:"-"`
	Albums    []Album   `json:"-"`
	Songs     []Song    `json:"-"`
}
