package album

import (
	"repertoire/server/internal"
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
		ReleaseDate: &[]internal.Date{internal.Date(time.Now())}[0],
		UserID:      Users[0].ID,
		ArtistID:    &[]uuid.UUID{Artists[0].ID}[0],
		Songs: []model.Song{
			{
				ID:           uuid.New(),
				Title:        "Test S1",
				UserID:       Users[0].ID,
				AlbumTrackNo: &[]uint{1}[0],
				ArtistID:     &[]uuid.UUID{Artists[0].ID}[0],
				Sections: []model.SongSection{
					{
						ID:                uuid.New(),
						Name:              "Test Song Section 1",
						Order:             0,
						Rehearsals:        15,
						Confidence:        10,
						ConfidenceScore:   12,
						RehearsalsScore:   45,
						Progress:          5,
						Occurrences:       2,
						SongSectionTypeID: Users[0].SongSectionTypes[0].ID,
						History: []model.SongSectionHistory{
							{
								ID:       uuid.New(),
								From:     0,
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
						Name:              "Test Song Section 2",
						Order:             1,
						Rehearsals:        20,
						Confidence:        30,
						ConfidenceScore:   25,
						RehearsalsScore:   45,
						Progress:          10,
						Occurrences:       5,
						SongSectionTypeID: Users[0].SongSectionTypes[0].ID,
						History: []model.SongSectionHistory{
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
				},
			},
			{
				ID:           uuid.New(),
				Title:        "Test S2",
				UserID:       Users[0].ID,
				AlbumTrackNo: &[]uint{2}[0],
				ArtistID:     &[]uuid.UUID{Artists[0].ID}[0],
				Sections: []model.SongSection{
					{
						ID:                uuid.New(),
						Name:              "Test Song Section 1",
						Order:             0,
						Occurrences:       2,
						SongSectionTypeID: Users[0].SongSectionTypes[0].ID,
					},
					{
						ID:                uuid.New(),
						Name:              "Test Song Section 2",
						Order:             1,
						Occurrences:       5,
						SongSectionTypeID: Users[0].SongSectionTypes[0].ID,
					},
					{
						ID:                uuid.New(),
						Name:              "Test Song Section 3",
						Order:             2,
						SongSectionTypeID: Users[0].SongSectionTypes[0].ID,
					},
				},
			},
			{
				ID:           uuid.New(),
				Title:        "Test S3",
				UserID:       Users[0].ID,
				AlbumTrackNo: &[]uint{3}[0],
				ArtistID:     &[]uuid.UUID{Artists[0].ID}[0],
			},
			{
				ID:           uuid.New(),
				Title:        "Test S4",
				UserID:       Users[0].ID,
				AlbumTrackNo: &[]uint{4}[0],
				ArtistID:     &[]uuid.UUID{Artists[0].ID}[0],
			},
		},
	},
	{
		ID:       uuid.New(),
		Title:    "Test Album 2",
		ImageURL: &[]internal.FilePath{"userId/Some image path/somewhere.jpeg"}[0],
		UserID:   Users[0].ID,
		Songs: []model.Song{
			{
				ID:           uuid.New(),
				Title:        "Test S1",
				UserID:       Users[0].ID,
				AlbumTrackNo: &[]uint{1}[0],
				Sections: []model.SongSection{
					{
						ID:                uuid.New(),
						Name:              "Test Song Section 1",
						Order:             0,
						Rehearsals:        15,
						RehearsalsScore:   45,
						Progress:          1,
						Occurrences:       2,
						SongSectionTypeID: Users[0].SongSectionTypes[0].ID,
						History: []model.SongSectionHistory{
							{
								ID:       uuid.New(),
								From:     0,
								To:       15,
								Property: model.RehearsalsProperty,
							},
						},
					},
				},
			},
		},
	},
	{
		ID:     uuid.New(),
		Title:  "Test Album 3",
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
		ID:       uuid.New(),
		Title:    "Test Song 4 - With Artist",
		ArtistID: &[]uuid.UUID{Artists[0].ID}[0],
		UserID:   Users[0].ID,
	},
	{
		ID:       uuid.New(),
		Title:    "Test Song 5 - With Different Artist",
		ArtistID: &[]uuid.UUID{Artists[1].ID}[0],
		UserID:   Users[0].ID,
	},
}
