package album

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
	return documents
}

var UserSearchID = uuid.New()
var AlbumSearchID = uuid.New()
var ArtistSearchID = uuid.New()

var SongSearches = []any{
	model.SongSearch{
		Title:       "Song 1",
		ReleaseDate: &[]string{time.Now().Format("2006-01-02")}[0],
		ImageUrl:    &[]internal.FilePath{"song-image.png"}[0],
		Artist: &model.SongArtistSearch{
			ID:        ArtistSearchID,
			Name:      "Artist 1",
			ImageUrl:  &[]internal.FilePath{"song/artist-image.png"}[0],
			UpdatedAt: time.Now().UTC(),
		},
		Album: &model.SongAlbumSearch{
			ID:          AlbumSearchID,
			Title:       "Album 1",
			ReleaseDate: &[]string{time.Now().Format("2006-01-02")}[0],
			ImageUrl:    &[]internal.FilePath{"song/artist-image.png"}[0],
			UpdatedAt:   time.Now().UTC(),
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
		Album: &model.SongAlbumSearch{
			ID:          AlbumSearchID,
			Title:       "Album 1",
			ReleaseDate: &[]string{time.Now().Format("2006-01-02")}[0],
			ImageUrl:    &[]internal.FilePath{"song/artist-image.png"}[0],
			UpdatedAt:   time.Now().UTC(),
		},
		SearchBase: model.SearchBase{
			ID:        "song-" + uuid.New().String(),
			UpdatedAt: time.Now().UTC(),
			CreatedAt: time.Now().UTC(),
			Type:      enums.Song,
			UserID:    UserSearchID,
		},
	},
}
