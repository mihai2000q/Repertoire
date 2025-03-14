package album

import (
	"github.com/google/uuid"
	"repertoire/server/internal"
	"repertoire/server/internal/enums"
	"repertoire/server/model"
	"strings"
	"time"
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
		ID:        id,
		Title:     a.Title,
		ImageUrl:  a.ImageUrl,
		UpdatedAt: a.UpdatedAt,
	}
}

var ArtistSearches = []any{
	model.ArtistSearch{
		Name:      "Metal",
		ImageUrl:  &[]internal.FilePath{"artist-image.png"}[0],
		UpdatedAt: time.Now().UTC(),
		SearchBase: model.SearchBase{
			ID:     "artist-" + uuid.New().String(),
			Type:   enums.Artist,
			UserID: UserID,
		},
	},
}

var AlbumSearches = []any{
	model.AlbumSearch{
		Title:     "Justice For All",
		ImageUrl:  &[]internal.FilePath{"album-image.png"}[0],
		UpdatedAt: time.Now().UTC(),
		Artist:    fromArtistSearchToAlbumArtistSearch(ArtistSearches[0].(model.ArtistSearch)),
		SearchBase: model.SearchBase{
			ID:     "album-" + uuid.New().String(),
			Type:   enums.Album,
			UserID: UserID,
		},
	},
}

var SongSearches = []any{
	model.SongSearch{
		Title:     "Justice",
		ImageUrl:  &[]internal.FilePath{"song-image.png"}[0],
		UpdatedAt: time.Now().UTC(),
		Artist:    fromArtistSearchToSongArtistSearch(ArtistSearches[0].(model.ArtistSearch)),
		Album:     fromAlbumSearchToSongAlbumSearch(AlbumSearches[0].(model.AlbumSearch)),
		SearchBase: model.SearchBase{
			ID:     "song-" + uuid.New().String(),
			Type:   enums.Song,
			UserID: UserID,
		},
	},
	model.SongSearch{
		Title:     "Justce", // typo on purpose
		ImageUrl:  &[]internal.FilePath{"song-image.png"}[0],
		UpdatedAt: time.Now().UTC(),
		SearchBase: model.SearchBase{
			ID:     "song-" + uuid.New().String(),
			Type:   enums.Song,
			UserID: UserID,
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
		Title:     "Best of all time",
		ImageUrl:  &[]internal.FilePath{"playlist-image.png"}[0],
		UpdatedAt: time.Now().UTC(),
		SearchBase: model.SearchBase{
			ID:     "playlist-" + uuid.New().String(),
			Type:   enums.Playlist,
			UserID: UserID,
		},
	},
	model.PlaylistSearch{
		Title:     "Just ice",
		ImageUrl:  &[]internal.FilePath{"playlist-image.png"}[0],
		UpdatedAt: time.Now().UTC(),
		SearchBase: model.SearchBase{
			ID:     "playlist-" + uuid.New().String(),
			Type:   enums.Playlist,
			UserID: UserID,
		},
	},
}
