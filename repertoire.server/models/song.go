package models

import (
	"time"

	"github.com/google/uuid"
)

type Song struct {
	ID         uuid.UUID `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Title      string    `gorm:"size:100; not null" json:"title"`
	IsRecorded *bool     `json:"isRecorded"`
	User       User      `gorm:"not null" json:"-"`
	CreatedAt  time.Time `gorm:"not null; <-:create" json:"created_at"`
	UpdatedAt  time.Time `gorm:"not null" json:"updated_at"`
}
