package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid; <-:create" json:"id"`
	Name      string    `json:"name"`
	Email     string    `gorm:"unique" json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `gorm:"<-:create" json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
