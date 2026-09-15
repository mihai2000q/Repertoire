package artist

import (
	"repertoire/server/internal"
	"repertoire/server/internal/date"
	"repertoire/server/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedData(db *gorm.DB) {
	db.Create(&Users)
	db.Create(&Artists)
	db.Create(&Albums)
	db.Create(&Songs)
	db.Create(&SongParts)
	db.Create(&SongArrangements)
	// Default Song Arrangements
	db.Model(&model.Song{ID: Artists[1].Songs[0].ID}).Update("default_arrangement_id", SongArrangements[0].ID)
	db.Model(&model.Song{ID: Artists[1].Songs[1].ID}).Update("default_arrangement_id", SongArrangements[2].ID)
}

var Users = []model.User{
	{
		ID:       uuid.New(),
		Name:     "John Doe",
		Email:    "johndoe@gmail.com",
		Password: "",
		BandMemberRoles: []model.BandMemberRole{
			{
				ID:    uuid.New(),
				Name:  "Guitarist",
				Order: 0,
			},
			{
				ID:    uuid.New(),
				Name:  "Vocalist",
				Order: 1,
			},
			{
				ID:    uuid.New(),
				Name:  "Manager",
				Order: 2,
			},
		},
	},
}

var Artists = []model.Artist{
	{
		ID:       uuid.New(),
		Name:     "Arduino",
		UserID:   Users[0].ID,
		ImageURL: &[]internal.FilePath{"userId/Some image path/somewhere.jpeg"}[0],
		IsBand:   true,
		BandMembers: []model.BandMember{
			{
				ID:       uuid.New(),
				Name:     "Member 1",
				Order:    0,
				ImageURL: &[]internal.FilePath{"userId/Some image path/somewhere.jpeg"}[0],
				Roles: []model.BandMemberRole{
					Users[0].BandMemberRoles[0],
					Users[0].BandMemberRoles[1],
				},
			},
			{
				ID:    uuid.New(),
				Name:  "Member 2",
				Order: 1,
				Roles: []model.BandMemberRole{Users[0].BandMemberRoles[1]},
			},
			{
				ID:    uuid.New(),
				Name:  "Member 3",
				Order: 2,
				Roles: []model.BandMemberRole{Users[0].BandMemberRoles[0]},
			},
		},
		Albums: []model.Album{
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
			{
				ID:     uuid.New(),
				Title:  "Test Album 3",
				UserID: Users[0].ID,
			},
			{
				ID:     uuid.New(),
				Title:  "Test Album 4",
				UserID: Users[0].ID,
			},
		},
	},
	{
		ID:     uuid.New(),
		Name:   "Metal",
		UserID: Users[0].ID,
		IsBand: false,
		Albums: []model.Album{
			{
				ID:     uuid.New(),
				Title:  "Test Album 1",
				UserID: Users[0].ID,
			},
		},
		Songs: []model.Song{
			{
				ID:         uuid.New(),
				Title:      "Test Song 1",
				UserID:     Users[0].ID,
				Rehearsals: 10,
				Confidence: 10,
				Progress:   5,
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
		},
	},
	{
		ID:     uuid.New(),
		Name:   "Rock",
		UserID: Users[0].ID,
		IsBand: true,
		Songs: []model.Song{
			{
				ID:     uuid.New(),
				Title:  "Test Song",
				UserID: Users[0].ID,
			},
		},
	},
}

var Albums = []model.Album{
	{
		ID:          uuid.New(),
		Title:       "Test Album 1",
		ReleaseDate: &[]date.Date{date.Date(time.Now())}[0],
		UserID:      Users[0].ID,
	},
	{
		ID:     uuid.New(),
		Title:  "Test Album 2",
		UserID: Users[0].ID,
		Songs: []model.Song{
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
		},
	},
}

var Songs = []model.Song{
	{
		ID:     uuid.New(),
		Title:  "Test Song 1",
		UserID: Users[0].ID,
	},
	{
		ID:    uuid.New(),
		Title: "Test Song 2",
		Album: &model.Album{
			ID:     uuid.New(),
			Title:  "Some Album",
			UserID: Users[0].ID,
		},
		UserID: Users[0].ID,
	},
	{
		ID:     uuid.New(),
		Title:  "Test Song 3",
		UserID: Users[0].ID,
	},

	{
		ID:       uuid.New(),
		Title:    "Test Song 1",
		ArtistID: &[]uuid.UUID{Artists[0].ID}[0],
		UserID:   Users[0].ID,
	},
	{
		ID:       uuid.New(),
		Title:    "Test Song 2",
		ArtistID: &[]uuid.UUID{Artists[0].ID}[0],
		UserID:   Users[0].ID,
	},
	{
		ID:       uuid.New(),
		Title:    "Test Song 3",
		ArtistID: &[]uuid.UUID{Artists[0].ID}[0],
		UserID:   Users[0].ID,
	},
	{
		ID:       uuid.New(),
		Title:    "Test Song 4",
		ArtistID: &[]uuid.UUID{Artists[0].ID}[0],
		UserID:   Users[0].ID,
	},

	{
		ID:       uuid.New(),
		Title:    "Album 2 Song 1",
		UserID:   Users[0].ID,
		ArtistID: &[]uuid.UUID{Artists[0].ID}[0],
		AlbumID:  &[]uuid.UUID{Artists[0].Albums[1].ID}[0],
	},
	{
		ID:       uuid.New(),
		Title:    "Album 2 Song 2",
		UserID:   Users[0].ID,
		ArtistID: &[]uuid.UUID{Artists[0].ID}[0],
		AlbumID:  &[]uuid.UUID{Artists[0].Albums[1].ID}[0],
	},
	{
		ID:       uuid.New(),
		Title:    "Album 4 Song 1",
		UserID:   Users[0].ID,
		ArtistID: &[]uuid.UUID{Artists[0].ID}[0],
		AlbumID:  &[]uuid.UUID{Artists[0].Albums[3].ID}[0],
	},
}

var SongParts = []model.SongPart{
	{
		ID:              uuid.New(),
		Name:            "Test Song 1 Part 1",
		SongOrder:       0,
		Rehearsals:      16,
		Confidence:      10,
		ConfidenceScore: 12,
		RehearsalsScore: 45,
		Progress:        5,
		SongID:          Artists[1].Songs[0].ID,
		History: []model.SongPartHistory{
			{
				ID:       uuid.New(),
				From:     0,
				To:       16,
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
		Name:            "Test Song 1 Part 2",
		SongOrder:       1,
		Rehearsals:      20,
		Confidence:      30,
		ConfidenceScore: 25,
		RehearsalsScore: 45,
		Progress:        10,
		SongID:          Artists[1].Songs[0].ID,
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
		ID:        uuid.New(),
		Name:      "Test Song 2 Part 1",
		SongOrder: 0,
		SongID:    Artists[1].Songs[1].ID,
	},

	{
		ID:        uuid.New(),
		Name:      "Test Song 3 Part 1",
		SongOrder: 0,
		SongID:    Artists[2].Songs[0].ID,
	},
}

var SongArrangements = []model.SongArrangement{
	{
		ID:     uuid.New(),
		Name:   "Test SongArrangement 1",
		SongID: Artists[1].Songs[0].ID,
		Order:  0,
		PartOccurrences: []model.SongPartOccurrences{
			{PartID: SongParts[0].ID, Occurrences: 7},
			{PartID: SongParts[1].ID, Occurrences: 0},
		},
	},
	{
		ID:     uuid.New(),
		Name:   "Test SongArrangement 2",
		SongID: Artists[1].Songs[0].ID,
		Order:  1,
		PartOccurrences: []model.SongPartOccurrences{
			{PartID: SongParts[0].ID, Occurrences: 5},
			{PartID: SongParts[1].ID, Occurrences: 11},
		},
	},

	{
		ID:     uuid.New(),
		Name:   "Test SongArrangement 1",
		SongID: Artists[1].Songs[1].ID,
		Order:  0,
		PartOccurrences: []model.SongPartOccurrences{
			{PartID: SongParts[2].ID, Occurrences: 7},
		},
	},
}
