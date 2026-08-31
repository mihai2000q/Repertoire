package playlist

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
	db.Create(&Playlists)
	db.Create(&Songs)
	db.Create(&SongParts)
	db.Create(&SongArrangements)
	// Default Song Arrangements
	db.Model(&model.Song{ID: Songs[0].ID}).Update("default_arrangement_id", SongArrangements[0].ID)
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

var Artists = []model.Artist{
	{
		ID:     uuid.New(),
		Name:   "Test Artist 1",
		UserID: Users[0].ID,
	},
	{
		ID:     uuid.New(),
		Name:   "Test Artist 2",
		UserID: Users[0].ID,
	},
}

var Albums = []model.Album{
	{
		ID:     uuid.New(),
		Title:  "Test Album 1",
		UserID: Users[0].ID,
	},
	{
		ID:     uuid.New(),
		Title:  "Test Album 2",
		UserID: Users[0].ID,
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
	{
		ID:     uuid.New(),
		Title:  "Test Playlist 3",
		UserID: Users[0].ID,
	},
}

var Songs = []model.Song{
	{
		ID:       uuid.New(),
		Title:    "Test Song 1",
		UserID:   Users[0].ID,
		AlbumID:  &[]uuid.UUID{Albums[0].ID}[0],
		ArtistID: &[]uuid.UUID{Artists[0].ID}[0],
	},
	{
		ID:       uuid.New(),
		Title:    "Test Song 2",
		UserID:   Users[0].ID,
		AlbumID:  &[]uuid.UUID{Albums[0].ID}[0],
		ArtistID: &[]uuid.UUID{Artists[0].ID}[0],
	},
	{
		ID:       uuid.New(),
		Title:    "Test Song 3",
		UserID:   Users[0].ID,
		AlbumID:  &[]uuid.UUID{Albums[1].ID}[0],
		ArtistID: &[]uuid.UUID{Artists[1].ID}[0],
	},
	{
		ID:       uuid.New(),
		Title:    "Test Song 4",
		UserID:   Users[0].ID,
		AlbumID:  &[]uuid.UUID{Albums[1].ID}[0],
		ArtistID: &[]uuid.UUID{Artists[1].ID}[0],
	},
	{
		ID:       uuid.New(),
		Title:    "Test Song 5",
		UserID:   Users[0].ID,
		AlbumID:  &[]uuid.UUID{Albums[1].ID}[0],
		ArtistID: &[]uuid.UUID{Artists[1].ID}[0],
	},
}

var SongParts = []model.SongPart{
	{
		ID:              uuid.New(),
		Name:            "Test Song Part 1",
		SongOrder:       0,
		SongID:          Songs[0].ID,
		Rehearsals:      15,
		Confidence:      10,
		ConfidenceScore: 12,
		RehearsalsScore: 45,
		Progress:        5,
		History: []model.SongPartHistory{
			{
				ID:       uuid.New(),
				From:     0,
				To:       15,
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
		ID:              uuid.New(),
		Name:            "Test Song Part 2",
		SongOrder:       1,
		SongID:          Songs[0].ID,
		Rehearsals:      20,
		Confidence:      30,
		ConfidenceScore: 25,
		RehearsalsScore: 45,
		Progress:        10,
		History: []model.SongPartHistory{
			{
				ID:       uuid.New(),
				From:     0,
				To:       20,
				Property: model.RehearsalsProperty,
			},
			{
				ID:       uuid.New(),
				From:     0,
				To:       30,
				Property: model.ConfidenceProperty,
			},
		},
	},

	{
		ID:              uuid.New(),
		Name:            "Test Song Part 1",
		SongOrder:       0,
		SongID:          Songs[1].ID,
		Rehearsals:      15,
		RehearsalsScore: 45,
		Progress:        1,
		History: []model.SongPartHistory{
			{
				ID:       uuid.New(),
				From:     0,
				To:       15,
				Property: model.RehearsalsProperty,
			},
		},
	},
}

var SongArrangements = []model.SongArrangement{
	{
		ID:     uuid.New(),
		Name:   "Test SongArrangement 1",
		SongID: Songs[0].ID,
		Order:  0,
		PartOccurrences: []model.SongPartOccurrences{
			{PartID: SongParts[0].ID, Occurrences: 7},
			{PartID: SongParts[1].ID, Occurrences: 0},
		},
	},
	{
		ID:     uuid.New(),
		Name:   "Test SongArrangement 2",
		SongID: Songs[0].ID,
		Order:  1,
		PartOccurrences: []model.SongPartOccurrences{
			{PartID: SongParts[0].ID, Occurrences: 5},
			{PartID: SongParts[1].ID, Occurrences: 11},
		},
	},
}

var PlaylistsSongs = []model.PlaylistSong{
	// Playlist 1
	{
		ID:          uuid.New(),
		PlaylistID:  Playlists[0].ID,
		SongID:      Songs[0].ID,
		SongTrackNo: 1,
	},
	{
		ID:          uuid.New(),
		PlaylistID:  Playlists[0].ID,
		SongID:      Songs[1].ID,
		SongTrackNo: 2,
	},
	{
		ID:          uuid.New(),
		PlaylistID:  Playlists[0].ID,
		SongID:      Songs[2].ID,
		SongTrackNo: 3,
	},
	{
		ID:          uuid.New(),
		PlaylistID:  Playlists[0].ID,
		SongID:      Songs[3].ID,
		SongTrackNo: 4,
	},
	{
		ID:          uuid.New(),
		PlaylistID:  Playlists[0].ID,
		SongID:      Songs[0].ID,
		SongTrackNo: 5,
	},

	// Playlist 2
	{
		ID:          uuid.New(),
		PlaylistID:  Playlists[1].ID,
		SongID:      Songs[0].ID,
		SongTrackNo: 1,
	},
	// Playlist 3
	{
		ID:          uuid.New(),
		PlaylistID:  Playlists[2].ID,
		SongID:      Songs[0].ID,
		SongTrackNo: 1,
	},
}
