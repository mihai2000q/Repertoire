package album

import (
	"repertoire/server/internal"
	"repertoire/server/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedData(db *gorm.DB) {
	db.Create(&Users)
	db.Create(&Artists)
	db.Create(&Songs)
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
			{
				ID:    uuid.New(),
				Name:  "Voice",
				Order: 2,
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

var Songs = []model.Song{
	{
		ID:             uuid.New(),
		Title:          "Test Song 1",
		GuitarTuningID: &[]uuid.UUID{Users[0].GuitarTunings[0].ID}[0],
		ArtistID:       &[]uuid.UUID{Artists[0].ID}[0],
		Sections: []model.SongSection{
			{
				ID:                uuid.New(),
				Name:              "Verse 1 - used on update",
				SongSectionTypeID: Users[0].SongSectionTypes[2].ID,
				BandMemberID:      &Artists[0].BandMembers[0].ID,
				InstrumentID:      &Users[0].Instruments[1].ID,
				Order:             0,
			},
			{
				ID:                uuid.New(),
				Name:              "Chorus 1 - used on delete",
				SongSectionTypeID: Users[0].SongSectionTypes[0].ID,
				BandMemberID:      &Artists[0].BandMembers[1].ID,
				InstrumentID:      &Users[0].Instruments[0].ID,
				Order:             1,
			},
		},
		UserID: Users[0].ID,
	},
	{
		ID:       uuid.New(),
		Title:    "Test Song 2",
		ArtistID: &[]uuid.UUID{Artists[0].ID}[0],
		UserID:   Users[0].ID,
	},
	{
		ID:     uuid.New(),
		Title:  "Test Song 3",
		UserID: Users[0].ID,
	},
}
