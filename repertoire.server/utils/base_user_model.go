package utils

import (
	"github.com/google/uuid"
)

type BaseUserModel struct {
	UserID uuid.UUID `gorm:"foreignKey:UserID; references:ID; notnull" json:"-"`
	BaseModel
}
