package user

import (
	"repertoire/server/internal"
	"repertoire/server/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
		SongSectionTypes: []model.SongSectionType{
			{
				ID:    uuid.New(),
				Name:  "Chorus",
				Order: 0,
			},
		},
		GuitarTunings: []model.GuitarTuning{
			{
				ID:    uuid.New(),
				Name:  "E Standard",
				Order: 0,
			},
		},
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
		ID:       uuid.New(),
		Title:    "Some Album",
		UserID:   Users[0].ID,
		ArtistID: &[]uuid.UUID{Artists[0].ID}[0],
	},
}

var Songs = []model.Song{
	{
		ID:             uuid.New(),
		Title:          "Some Song",
		GuitarTuningID: &[]uuid.UUID{Users[0].GuitarTunings[0].ID}[0],
		ArtistID:       &[]uuid.UUID{Artists[0].ID}[0],
		AlbumID:        &[]uuid.UUID{Albums[0].ID}[0],
		AlbumTrackNo:   &[]uint{1}[0],
		UserID:         Users[0].ID,
		Sections: []model.SongSection{
			{
				ID:                uuid.New(),
				Name:              "Chorus 1",
				SongSectionTypeID: Users[0].SongSectionTypes[0].ID,
			},
		},
	},
}

var Playlists = []model.Playlist{
	{
		ID:     uuid.New(),
		Title:  "Some Playlist",
		UserID: Users[0].ID,
	},
}
