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
		Instruments: []model.Instrument{
			{
				ID:    uuid.New(),
				Name:  "Guitar",
				Order: 0,
			},
			{
				ID:    uuid.New(),
				Name:  "Piano",
				Order: 1,
			},
		},
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
		ImageURL: &[]internal.FilePath{"userId/Some artist image path/somewhere.jpeg"}[0],
		UserID:   Users[0].ID,
		BandMembers: []model.BandMember{
			{
				ID:    uuid.New(),
				Name:  "Member 1",
				Order: 0,
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
		ImageURL:    &[]internal.FilePath{"userId/Some album image path/somewhere.jpeg"}[0],
		UserID:      Users[0].ID,
		ArtistID:    &[]uuid.UUID{Artists[0].ID}[0],
	},
	{
		ID:     uuid.New(),
		Title:  "Test Album 2",
		UserID: Users[0].ID,
	},
	{
		ID:       uuid.New(),
		Title:    "Test Album 3",
		UserID:   Users[0].ID,
		ArtistID: &[]uuid.UUID{Artists[1].ID}[0],
	},
}

var songSections = []model.SongSection{
	{
		ID:                uuid.New(),
		Name:              "Verse 1 - used on update",
		SongSectionTypeID: Users[0].SongSectionTypes[2].ID,
		BandMemberID:      &Artists[0].BandMembers[0].ID,
		InstrumentID:      &Users[0].Instruments[1].ID,
		Order:             0,
		Confidence:        10,
		Rehearsals:        10,
		ConfidenceScore:   12,
		RehearsalsScore:   45,
		Progress:          5,
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
		Name:              "Chorus 1 - used on delete",
		SongSectionTypeID: Users[0].SongSectionTypes[0].ID,
		BandMemberID:      &Artists[0].BandMembers[1].ID,
		Order:             1,
		Confidence:        25,
		Rehearsals:        50,
		ConfidenceScore:   30,
		RehearsalsScore:   550,
		Progress:          99,
		InstrumentID:      &Users[0].Instruments[0].ID,
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
}

var Songs = []model.Song{
	{
		ID:            uuid.New(),
		Title:         "Test Song 1",
		Description:   "Some description",
		ReleaseDate:   &[]internal.Date{internal.Date(time.Now())}[0],
		ImageURL:      &[]internal.FilePath{"userId/Some image path/somewhere.jpeg"}[0],
		IsRecorded:    true,
		Bpm:           &[]uint{123}[0],
		Difficulty:    &[]enums.Difficulty{enums.Easy}[0],
		SongsterrLink: &[]string{"https://songster.com/some-song"}[0],
		YoutubeLink:   &[]string{"https://youtube.com/some-song"}[0],

		Confidence: 8.75,
		Rehearsals: 15,
		Progress:   26,

		GuitarTuningID: &[]uuid.UUID{Users[0].GuitarTunings[0].ID}[0],
		ArtistID:       &[]uuid.UUID{Artists[0].ID}[0],
		AlbumID:        &[]uuid.UUID{Albums[0].ID}[0],
		AlbumTrackNo:   &[]uint{1}[0],

		Sections: songSections,
		UserID:   Users[0].ID,
	},
	{
		ID:           uuid.New(),
		Title:        "Test Song 2",
		ArtistID:     &[]uuid.UUID{Artists[0].ID}[0],
		AlbumID:      &[]uuid.UUID{Albums[0].ID}[0],
		AlbumTrackNo: &[]uint{2}[0],
		Settings:     model.SongSettings{ID: uuid.New()},
		UserID:       Users[0].ID,
	},
	{
		ID:           uuid.New(),
		Title:        "Test Song 3",
		ArtistID:     &[]uuid.UUID{Artists[0].ID}[0],
		AlbumID:      &[]uuid.UUID{Albums[0].ID}[0],
		AlbumTrackNo: &[]uint{3}[0],
		Bpm:          &[]uint{81}[0],
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
		Title:  "Test Song 5 - No Album - But Has Section Occurrences",
		UserID: Users[0].ID,
		Sections: []model.SongSection{
			{
				ID:                 uuid.New(),
				Name:               "Test Song Section 1",
				Order:              0,
				Rehearsals:         15,
				Occurrences:        2,
				PartialOccurrences: 1,
				SongSectionTypeID:  Users[0].SongSectionTypes[1].ID,
				History: []model.SongSectionHistory{
					{
						ID:       uuid.New(),
						From:     0,
						To:       15,
						Property: model.RehearsalsProperty,
					},
				},
			},
			{
				ID:                uuid.New(),
				Name:              "Test Song Section 2",
				Order:             1,
				SongSectionTypeID: Users[0].SongSectionTypes[0].ID,
			},
			{
				ID:                 uuid.New(),
				Name:               "Test Song Section 3",
				Order:              2,
				Occurrences:        10,
				PartialOccurrences: 7,
				SongSectionTypeID:  Users[0].SongSectionTypes[0].ID,
			},
		},
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

	{
		ID:           uuid.New(),
		Title:        "Test Song 10 - In Another Album",
		UserID:       Users[0].ID,
		AlbumID:      &[]uuid.UUID{Albums[2].ID}[0],
		AlbumTrackNo: &[]uint{1}[0],
	},
}

var PlaylistSongs = []model.PlaylistSong{
	{
		ID:          uuid.New(),
		PlaylistID:  Playlists[0].ID,
		SongID:      Songs[0].ID,
		SongTrackNo: 1,
	},

	{
		ID:          uuid.New(),
		PlaylistID:  Playlists[1].ID,
		SongID:      Songs[5].ID,
		SongTrackNo: 1,
	},
	{
		ID:          uuid.New(),
		PlaylistID:  Playlists[1].ID,
		SongID:      Songs[6].ID,
		SongTrackNo: 2,
	},
	{
		ID:          uuid.New(),
		PlaylistID:  Playlists[1].ID,
		SongID:      Songs[7].ID,
		SongTrackNo: 3,
	},
	{
		ID:          uuid.New(),
		PlaylistID:  Playlists[1].ID,
		SongID:      Songs[8].ID,
		SongTrackNo: 4,
	},
}
