package model

import (
	"repertoire/utils"
)

type Artist struct {
	Name   string  `gorm:"size:100; not null" json:"name"`
	Albums []Album `json:"-"`
	Songs  []Song  `json:"-"`
	utils.BaseUserModel
}
