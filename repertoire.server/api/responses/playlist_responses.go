package responses

import "github.com/google/uuid"

type AddSongsToPlaylistResponse struct {
	Success    bool        `json:"success"`
	Duplicates []uuid.UUID `json:"duplicates"`
	Added      []uuid.UUID `json:"added"`
}

type AddAlbumsToPlaylistResponse struct {
	Success           bool        `json:"success"`
	DuplicateSongIDs  []uuid.UUID `json:"duplicateSongIDs"`
	DuplicateAlbumIDs []uuid.UUID `json:"duplicateAlbumIDs"`
	AddedSongIDs      []uuid.UUID `json:"addedSongIDs"`
}

type AddArtistsToPlaylistResponse struct {
	Success            bool        `json:"success"`
	DuplicateSongIDs   []uuid.UUID `json:"duplicateSongIDs"`
	DuplicateArtistIDs []uuid.UUID `json:"duplicateArtistIDs"`
	AddedSongIDs       []uuid.UUID `json:"addedSongIDs"`
}
