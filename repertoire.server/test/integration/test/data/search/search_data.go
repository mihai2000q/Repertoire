package album

import (
	"repertoire/server/internal"
	"repertoire/server/internal/enums"
	"repertoire/server/model"
	"strings"
	"time"

	"github.com/google/uuid"
)

func GetSearchDocuments() []any {
	var documents []any
	documents = append(documents, ArtistSearches...)
	documents = append(documents, AlbumSearches...)
	documents = append(documents, SongSearches...)
	documents = append(documents, PlaylistSearches...)
	return documents
}

var UserID = uuid.New()

func fromArtistSearchToAlbumArtistSearch(a model.ArtistSearch) *model.AlbumArtistSearch {
	id, _ := uuid.Parse(strings.Replace(a.ID, "artist-", "", 1))
	return &model.AlbumArtistSearch{
		ID:        id,
		Name:      a.Name,
		ImageUrl:  a.ImageUrl,
		UpdatedAt: a.UpdatedAt,
	}
}

func fromArtistSearchToSongArtistSearch(a model.ArtistSearch) *model.SongArtistSearch {
	id, _ := uuid.Parse(strings.Replace(a.ID, "artist-", "", 1))
	return &model.SongArtistSearch{
		ID:        id,
		Name:      a.Name,
		ImageUrl:  a.ImageUrl,
		UpdatedAt: a.UpdatedAt,
	}
}

func fromAlbumSearchToSongAlbumSearch(a model.AlbumSearch) *model.SongAlbumSearch {
	id, _ := uuid.Parse(strings.Replace(a.ID, "album-", "", 1))
	return &model.SongAlbumSearch{
		ID:          id,
		Title:       a.Title,
		ImageUrl:    a.ImageUrl,
		ReleaseDate: a.ReleaseDate,
		UpdatedAt:   a.UpdatedAt,
	}
}

var ArtistSearches = []any{
	model.ArtistSearch{
		Name:     "Metal",
		ImageUrl: &[]internal.FilePath{"artist-image.png"}[0],
		SearchBase: model.SearchBase{
			ID:        "artist-" + uuid.New().String(),
			UpdatedAt: time.Now().UTC(),
			CreatedAt: time.Now().UTC(),
			Type:      enums.Artist,
			UserID:    UserID,
		},
	},
}

var AlbumSearches = []any{
	model.AlbumSearch{
		Title:       "Justice For All",
		ImageUrl:    &[]internal.FilePath{"album-image.png"}[0],
		ReleaseDate: &[]string{time.Now().Format("YYYY-MM-DD")}[0],
		Artist:      fromArtistSearchToAlbumArtistSearch(ArtistSearches[0].(model.ArtistSearch)),
		SearchBase: model.SearchBase{
			ID:        "album-" + uuid.New().String(),
			UpdatedAt: time.Now().UTC(),
			CreatedAt: time.Now().UTC(),
			Type:      enums.Album,
			UserID:    UserID,
		},
	},
}

var SongSearches = []any{
	model.SongSearch{
		Title:       "Justice",
		ReleaseDate: &[]string{time.Now().Format("YYYY-MM-DD")}[0],
		ImageUrl:    &[]internal.FilePath{"song-image.png"}[0],
		Artist:      fromArtistSearchToSongArtistSearch(ArtistSearches[0].(model.ArtistSearch)),
		Album:       fromAlbumSearchToSongAlbumSearch(AlbumSearches[0].(model.AlbumSearch)),
		SearchBase: model.SearchBase{
			ID:        "song-" + uuid.New().String(),
			UpdatedAt: time.Now().UTC(),
			CreatedAt: time.Now().UTC(),
			Type:      enums.Song,
			UserID:    UserID,
		},
	},
	model.SongSearch{
		Title:    "Justce", // typo on purpose
		ImageUrl: &[]internal.FilePath{"song-image.png"}[0],
		SearchBase: model.SearchBase{
			ID:        "song-" + uuid.New().String(),
			UpdatedAt: time.Now().UTC(),
			CreatedAt: time.Now().UTC(),
			Type:      enums.Song,
			UserID:    UserID,
		},
	},
	model.SongSearch{
		Title:  "Master",
		Artist: fromArtistSearchToSongArtistSearch(ArtistSearches[0].(model.ArtistSearch)),
		SearchBase: model.SearchBase{
			ID:     "song-" + uuid.New().String(),
			Type:   enums.Song,
			UserID: UserID,
		},
	},
}

var PlaylistSearches = []any{
	model.PlaylistSearch{
		Title:    "Best of all time",
		ImageUrl: &[]internal.FilePath{"playlist-image.png"}[0],
		SearchBase: model.SearchBase{
			ID:        "playlist-" + uuid.New().String(),
			UpdatedAt: time.Now().UTC(),
			CreatedAt: time.Now().UTC(),
			Type:      enums.Playlist,
			UserID:    UserID,
		},
	},
	model.PlaylistSearch{
		Title:    "Just ice",
		ImageUrl: &[]internal.FilePath{"playlist-image.png"}[0],
		SearchBase: model.SearchBase{
			ID:        "playlist-" + uuid.New().String(),
			UpdatedAt: time.Now().UTC(),
			CreatedAt: time.Now().UTC(),
			Type:      enums.Playlist,
			UserID:    UserID,
		},
	},
}
