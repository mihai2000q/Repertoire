package model

import (
	"github.com/google/uuid"
	"repertoire/utils"
)

type Album struct {
	Title    string     `gorm:"size:100; not null" json:"title"`
	ArtistID *uuid.UUID `json:"-"`
	Artist   Artist     `json:"-"`
	Songs    []Song     `json:"-"`
	utils.BaseUserModel
}
