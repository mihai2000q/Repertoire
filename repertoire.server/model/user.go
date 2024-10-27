package model

import (
	"repertoire/utils"
)

type User struct {
	Name            string            `gorm:"size:100; not null" json:"name"`
	Email           string            `gorm:"size:256; unique; not null" json:"email"`
	Password        string            `gorm:"not null" json:"-"`
	Albums          []Album           `json:"-"`
	Artists         []Artist          `json:"-"`
	Playlists       []Playlist        `json:"-"`
	Songs           []Song            `json:"-"`
	utils.BaseModel
}
