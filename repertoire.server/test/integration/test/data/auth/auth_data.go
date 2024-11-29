package auth

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"repertoire/server/model"
)

func SeedData(db *gorm.DB) {
	db.Create(&Users)
}

var Users = []model.User{
	{
		ID:       uuid.New(),
		Name:     "John Doe",
		Email:    "johndoe@gmail.com",
		Password: "$2a$10$EXl0YQUN4AHaV6ZuRGXMheQLoJo6Hb4iy/IdaruL/e0pIbibcvn5C",
	},
}
