package model

import (
	"time"

	"github.com/google/uuid"
)

type GuitarTuning struct {
	ID        uuid.UUID `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Name      string    `gorm:"size:30; not null" json:"name"`
	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`
	UserID    uuid.UUID `gorm:"not null" json:"-"`
	Songs     []Song    `json:"-"`
}
