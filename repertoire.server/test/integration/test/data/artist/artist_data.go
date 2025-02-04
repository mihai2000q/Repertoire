package artist

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
				Order:    1,
				ImageURL: &[]internal.FilePath{"userId/Some image path/somewhere.jpeg"}[0],
				Roles: []model.BandMemberRole{
					Users[0].BandMemberRoles[0],
					Users[0].BandMemberRoles[1],
				},
			},
			{
				ID:    uuid.New(),
				Name:  "Member 2",
				Order: 2,
				Roles: []model.BandMemberRole{Users[0].BandMemberRoles[1]},
			},
			{
				ID:    uuid.New(),
				Name:  "Member 3",
				Order: 3,
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
		},
	},
}

var Albums = []model.Album{
	{
		ID:          uuid.New(),
		Title:       "Test Album 1",
		ReleaseDate: &[]time.Time{time.Now()}[0],
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
