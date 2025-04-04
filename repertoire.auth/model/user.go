package model

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	Email    string    `gorm:"size:256; unique; not null" json:"email"`
	Password string    `gorm:"not null" json:"-"`
	Name     string    `gorm:"size:100; not null" json:"name"`
}
