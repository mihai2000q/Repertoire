package utils

import (
	"github.com/google/uuid"
	"time"
)

type BaseModel struct {
	ID        uuid.UUID `gorm:"primaryKey; type:uuid; <-:create" json:"id"`
	CreatedAt time.Time `gorm:"default:current_timestamp; not null; <-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:current_timestamp; not null" json:"updatedAt"`
}
