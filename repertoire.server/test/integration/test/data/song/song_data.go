package album

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"repertoire/server/internal"
	"repertoire/server/internal/enums"
	"repertoire/server/model"
	"time"
)

func SeedData(db *gorm.DB) {
	db.Create(&Users)
	db.Create(&Playlists)
	db.Create(&Artists)
	db.Create(&Albums)
	db.Create(&Songs)
	db.Create(&PlaylistSongs)
}

var Playlists = []model.Playlist{
	{
		ID:     uuid.New(),
		Title:  "Test Playlist 1",
		UserID: Users[0].ID,
	},
	{
		ID:     uuid.New(),
		Title:  "Test Playlist 2",
		UserID: Users[0].ID,
	},
}

var Users = []model.User{
	{
		ID:       uuid.New(),
		Name:     "John Doe",
		Email:    "johndoe@gmail.com",
		Password: "",
		GuitarTunings: []model.GuitarTuning{
			{
				ID:    uuid.New(),
				Name:  "E Standard",
				Order: 0,
			},
			{
				ID:    uuid.New(),
				Name:  "Drop C",
				Order: 1,
			},
			{
				ID:    uuid.New(),
				Name:  "Drop B",
				Order: 2,
			},
		},
		SongSectionTypes: []model.SongSectionType{
			{
				ID:    uuid.New(),
				Name:  "Chorus",
				Order: 0,
			},
			{
				ID:    uuid.New(),
				Name:  "Solo",
				Order: 1,
			},
			{
				ID:    uuid.New(),
				Name:  "Verse",
				Order: 2,
			},
		},
	},
}

var Artists = []model.Artist{
	{
		ID:     uuid.New(),
		Name:   "Arduino",
		UserID: Users[0].ID,
	},
	{
		ID:     uuid.New(),
		Name:   "Metal",
		UserID: Users[0].ID,
	},
}

var Albums = []model.Album{
	{
		ID:          uuid.New(),
		Title:       "Test Album 1",
		ReleaseDate: &[]time.Time{time.Now()}[0],
		UserID:      Users[0].ID,
		ArtistID:    &[]uuid.UUID{Artists[0].ID}[0],
	},
	{
		ID:     uuid.New(),
		Title:  "Test Album 2",
		UserID: Users[0].ID,
	},
}

var Songs = []model.Song{
	{
		ID:            uuid.New(),
		Title:         "Test Song 1",
		Description:   "Some description",
		ReleaseDate:   &[]time.Time{time.Now()}[0],
		ImageURL:      &[]internal.FilePath{"userId/Some image path/somewhere.jpeg"}[0],
		IsRecorded:    true,
		Bpm:           &[]uint{123}[0],
		Difficulty:    &[]enums.Difficulty{enums.Easy}[0],
		SongsterrLink: &[]string{"https://songster.com/some-song"}[0],
		YoutubeLink:   &[]string{"https://youtube.com/some-song"}[0],

		GuitarTuningID: &[]uuid.UUID{Users[0].GuitarTunings[0].ID}[0],
		ArtistID:       &[]uuid.UUID{Artists[0].ID}[0],
		AlbumID:        &[]uuid.UUID{Albums[0].ID}[0],
		AlbumTrackNo:   &[]uint{1}[0],

		Sections: []model.SongSection{
			{
				ID:                uuid.New(),
				Name:              "Verse 1",
				SongSectionTypeID: Users[0].SongSectionTypes[2].ID,
				Order:             0,
				Confidence:        10,
				Rehearsals:        10,
				ConfidenceScore:   12,
				RehearsalsScore:   45,
				History: []model.SongSectionHistory{
					{
						ID:       uuid.New(),
						From:     0,
						To:       5,
						Property: model.RehearsalsProperty,
					},
					{
						ID:       uuid.New(),
						From:     5,
						To:       10,
						Property: model.RehearsalsProperty,
					},
					{
						ID:       uuid.New(),
						From:     0,
						To:       10,
						Property: model.ConfidenceProperty,
					},
				},
			},
			{
				ID:                uuid.New(),
				Name:              "Chorus 1",
				SongSectionTypeID: Users[0].SongSectionTypes[0].ID,
				Order:             1,
			},
			{
				ID:                uuid.New(),
				Name:              "Verse 2",
				SongSectionTypeID: Users[0].SongSectionTypes[2].ID,
				Order:             2,
			},
			{
				ID:                uuid.New(),
				Name:              "Solo",
				SongSectionTypeID: Users[0].SongSectionTypes[1].ID,
				Rehearsals:        10,
				Order:             3,
			},
		},
		UserID: Users[0].ID,
	},
	{
		ID:           uuid.New(),
		Title:        "Test Song 2",
		ArtistID:     &[]uuid.UUID{Artists[0].ID}[0],
		AlbumID:      &[]uuid.UUID{Albums[0].ID}[0],
		AlbumTrackNo: &[]uint{2}[0],
		UserID:       Users[0].ID,
	},
	{
		ID:           uuid.New(),
		Title:        "Test Song 3",
		ArtistID:     &[]uuid.UUID{Artists[0].ID}[0],
		AlbumID:      &[]uuid.UUID{Albums[0].ID}[0],
		AlbumTrackNo: &[]uint{3}[0],
		UserID:       Users[0].ID,
	},
	{
		ID:           uuid.New(),
		Title:        "Test Song 4",
		ArtistID:     &[]uuid.UUID{Artists[0].ID}[0],
		AlbumID:      &[]uuid.UUID{Albums[0].ID}[0],
		ImageURL:     &[]internal.FilePath{"userId/Some image path/somewhere.jpeg"}[0],
		AlbumTrackNo: &[]uint{4}[0],
		UserID:       Users[0].ID,
	},

	{
		ID:     uuid.New(),
		Title:  "Test Song 5 - No Album",
		UserID: Users[0].ID,
	},

	{
		ID:     uuid.New(),
		Title:  "Test Song 6 - In Playlist",
		UserID: Users[0].ID,
	},
	{
		ID:     uuid.New(),
		Title:  "Test Song 7 - In Playlist",
		UserID: Users[0].ID,
	},
	{
		ID:     uuid.New(),
		Title:  "Test Song 8 - In Playlist",
		UserID: Users[0].ID,
	},
	{
		ID:     uuid.New(),
		Title:  "Test Song 9 - In Playlist",
		UserID: Users[0].ID,
	},
}

var PlaylistSongs = []model.PlaylistSong{
	{
		PlaylistID:  Playlists[0].ID,
		SongID:      Songs[0].ID,
		SongTrackNo: 1,
	},

	{
		PlaylistID:  Playlists[1].ID,
		SongID:      Songs[5].ID,
		SongTrackNo: 1,
	},
	{
		PlaylistID:  Playlists[1].ID,
		SongID:      Songs[6].ID,
		SongTrackNo: 2,
	},
	{
		PlaylistID:  Playlists[1].ID,
		SongID:      Songs[7].ID,
		SongTrackNo: 3,
	},
	{
		PlaylistID:  Playlists[1].ID,
		SongID:      Songs[8].ID,
		SongTrackNo: 4,
	},
}
