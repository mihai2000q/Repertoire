package artist

import (
	"github.com/google/uuid"
	"repertoire/server/internal"
	"repertoire/server/internal/enums"
	"repertoire/server/model"
	"time"
)

func GetSearchDocuments() []any {
	var documents []any
	documents = append(documents, SongSearches...)
	documents = append(documents, AlbumSearches...)
	return documents
}

var UserSearchID = uuid.New()
var ArtistSearchID = uuid.New()

var AlbumSearches = []any{
	model.AlbumSearch{
		Title:    "Album 1",
		ImageUrl: &[]internal.FilePath{"song-image.png"}[0],
		Artist: &model.AlbumArtistSearch{
			ID:        ArtistSearchID,
			Name:      "Artist 1",
			ImageUrl:  &[]internal.FilePath{"song/artist-image.png"}[0],
			UpdatedAt: time.Now().UTC(),
		},
		SearchBase: model.SearchBase{
			ID:        "album-" + uuid.New().String(),
			UpdatedAt: time.Now().UTC(),
			CreatedAt: time.Now().UTC(),
			Type:      enums.Album,
			UserID:    UserSearchID,
		},
	},
	model.AlbumSearch{
		Title: "Album 2",
		Artist: &model.AlbumArtistSearch{
			ID:        ArtistSearchID,
			Name:      "Artist 1",
			ImageUrl:  &[]internal.FilePath{"song/artist-image.png"}[0],
			UpdatedAt: time.Now().UTC(),
		},
		SearchBase: model.SearchBase{
			ID:     "song-" + uuid.New().String(),
			Type:   enums.Song,
			UserID: UserSearchID,
		},
	},
}

var SongSearches = []any{
	model.SongSearch{
		Title:    "Song 1",
		ImageUrl: &[]internal.FilePath{"song-image.png"}[0],
		Artist: &model.SongArtistSearch{
			ID:        ArtistSearchID,
			Name:      "Artist 1",
			ImageUrl:  &[]internal.FilePath{"song/artist-image.png"}[0],
			UpdatedAt: time.Now().UTC(),
		},
		Album: &model.SongAlbumSearch{
			ID:        uuid.New(),
			Title:     "Album 1",
			ImageUrl:  &[]internal.FilePath{"song/artist-image.png"}[0],
			UpdatedAt: time.Now().UTC(),
		},
		SearchBase: model.SearchBase{
			ID:        "song-" + uuid.New().String(),
			UpdatedAt: time.Now().UTC(),
			CreatedAt: time.Now().UTC(),
			Type:      enums.Song,
			UserID:    UserSearchID,
		},
	},
	model.SongSearch{
		Title: "Song 2",
		Artist: &model.SongArtistSearch{
			ID:        ArtistSearchID,
			Name:      "Artist 1",
			ImageUrl:  &[]internal.FilePath{"song/artist-image.png"}[0],
			UpdatedAt: time.Now().UTC(),
		},
		SearchBase: model.SearchBase{
			ID:     "song-" + uuid.New().String(),
			Type:   enums.Song,
			UserID: UserSearchID,
		},
	},
}
