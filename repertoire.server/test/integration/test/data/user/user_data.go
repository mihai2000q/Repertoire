package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"repertoire/server/internal"
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
		Password: "",
	},
	{
		ID:                uuid.New(),
		Name:              "Vivaldi",
		Email:             "vivaldi@gmail.com",
		Password:          "",
		ProfilePictureURL: &[]internal.FilePath{"userId/Some image path/somewhere.jpeg"}[0],
	},
}
