package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid; <-:create" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"unique; not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	CreatedAt time.Time `gorm:"not null; <-:create" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
}
