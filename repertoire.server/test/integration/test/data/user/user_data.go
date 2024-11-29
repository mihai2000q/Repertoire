package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"repertoire/server/internal"
	"repertoire/server/model"
)

func SeedData(db *gorm.DB) {
	db.Create(&Users)
	db.Create(&Artists)
	db.Create(&Albums)
	db.Create(&Songs)
	db.Create(&Playlists)
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

var Artists = []model.Artist{
	{
		ID:     uuid.New(),
		Name:   "Some Artist",
		UserID: Users[0].ID,
	},
}

var Albums = []model.Album{
	{
		ID:     uuid.New(),
		Title:  "Some Album",
		UserID: Users[0].ID,
	},
}

var Songs = []model.Song{
	{
		ID:     uuid.New(),
		Title:  "Some Song",
		UserID: Users[0].ID,
	},
}

var Playlists = []model.Playlist{
	{
		ID:     uuid.New(),
		Title:  "Some Playlist",
		UserID: Users[0].ID,
	},
}
