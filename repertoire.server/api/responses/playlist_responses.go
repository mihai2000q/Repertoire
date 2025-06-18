package responses

import "github.com/google/uuid"

type AddSongsToPlaylistResponse struct {
	Success    bool        `json:"success"`
	Duplicates []uuid.UUID `json:"duplicates"`
	Added      []uuid.UUID `json:"added"`
}

type AddAlbumsToPlaylistResponse struct {
	Success           bool        `json:"success"`
	DuplicateSongIDs  []uuid.UUID `json:"duplicateSongIds"`
	DuplicateAlbumIDs []uuid.UUID `json:"duplicateAlbumIds"`
	AddedSongIDs      []uuid.UUID `json:"addedSongIds"`
}

type AddArtistsToPlaylistResponse struct {
	Success            bool        `json:"success"`
	DuplicateSongIDs   []uuid.UUID `json:"duplicateSongIds"`
	DuplicateArtistIDs []uuid.UUID `json:"duplicateArtistIds"`
	AddedSongIDs       []uuid.UUID `json:"addedSongIds"`
}
