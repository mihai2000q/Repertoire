package model

import (
	"repertoire/utils"
)

type Playlist struct {
	Title       string `gorm:"size:100; not null" json:"title"`
	Description string `gorm:"not null" json:"description"`
	Songs       []Song `gorm:"many2many:playlist_song" json:"-"`
	utils.BaseUserModel
}
