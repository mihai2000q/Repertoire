package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `gorm:"primaryKey,type:uuid; <-:create" json:"id"`
	Name      string    `gorm:"size:100; not null" json:"name"`
	Email     string    `gorm:"size:256;unique; not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	CreatedAt time.Time `gorm:"not null; <-:create" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
}
