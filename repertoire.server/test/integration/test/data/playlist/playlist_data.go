package playlist

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"repertoire/server/internal"
	"repertoire/server/model"
)

func SeedData(db *gorm.DB) {
	db.Create(&Users)
	db.Create(&Playlists)
	db.Create(&Songs)
	db.Create(&PlaylistsSongs)
}

var Users = []model.User{
	{
		ID:       uuid.New(),
		Name:     "John Doe",
		Email:    "johndoe@gmail.com",
		Password: "",
	},
}

var Playlists = []model.Playlist{
	{
		ID:          uuid.New(),
		Title:       "Test Playlist",
		Description: "This is a test playlist",
		ImageURL:    &[]internal.FilePath{"userId/Some image path/somewhere.jpeg"}[0],
		UserID:      Users[0].ID,
	},
	{
		ID:     uuid.New(),
		Title:  "Test Playlist 2",
		UserID: Users[0].ID,
	},
}

var Songs = []model.Song{
	{
		ID:     uuid.New(),
		Title:  "Test Song 1",
		UserID: Users[0].ID,
	},
	{
		ID:     uuid.New(),
		Title:  "Test Song 2",
		UserID: Users[0].ID,
	},
	{
		ID:     uuid.New(),
		Title:  "Test Song 3",
		UserID: Users[0].ID,
	},
	{
		ID:     uuid.New(),
		Title:  "Test Song 4",
		UserID: Users[0].ID,
	},
	{
		ID:     uuid.New(),
		Title:  "Test Song 5",
		UserID: Users[0].ID,
	},
}

var PlaylistsSongs = []model.PlaylistSong{
	// Playlist 1
	{
		PlaylistID:  Playlists[0].ID,
		SongID:      Songs[0].ID,
		SongTrackNo: 1,
	},
	{
		PlaylistID:  Playlists[0].ID,
		SongID:      Songs[1].ID,
		SongTrackNo: 2,
	},
	{
		PlaylistID:  Playlists[0].ID,
		SongID:      Songs[2].ID,
		SongTrackNo: 3,
	},
	{
		PlaylistID:  Playlists[0].ID,
		SongID:      Songs[3].ID,
		SongTrackNo: 4,
	},

	// Playlist 2
	{
		PlaylistID:  Playlists[1].ID,
		SongID:      Songs[0].ID,
		SongTrackNo: 1,
	},
}