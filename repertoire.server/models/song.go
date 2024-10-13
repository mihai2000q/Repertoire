package models

import (
	"github.com/google/uuid"
	"time"
)

type Song struct {
	ID         uuid.UUID  `gorm:"type:uuid; <-:create" json:"id"`
	Title      string     `gorm:"not null" json:"title"`
	IsRecorded bool 	  `json: "isRecorded`
	CreatedAt  time.Time  `gorm:"not null; <-:create" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"not null" json:"updated_at"`
}
