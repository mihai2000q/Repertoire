package models

import (
	"time"

	"github.com/google/uuid"
)

type Song struct {
	ID         uuid.UUID `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Title      string    `gorm:"size:100; not null" json:"title"`
	IsRecorded *bool     `json:"isRecorded"`
	CreatedAt  time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`
	UserID     uuid.UUID `gorm:"foreignKey:UserID; references: ID; type:uuid; not null" json:"-"`
}
