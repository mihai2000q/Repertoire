package requests

import "github.com/google/uuid"

type GetPlaylistRequest struct {
	ID uuid.UUID `validate:"required"`
}

type GetPlaylistsRequest struct {
	CurrentPage *int     `form:"currentPage" validate:"required_with=PageSize,omitempty,gt=0"`
	PageSize    *int     `form:"pageSize" validate:"required_with=CurrentPage,omitempty,gt=0"`
	OrderBy     []string `form:"orderBy" validate:"order_by"`
	SearchBy    []string `form:"searchBy" validate:"search_by"`
}

type GetPlaylistSongsRequest struct {
	ID          uuid.UUID `validate:"required"`
	CurrentPage *int      `form:"currentPage" validate:"required_with=PageSize,omitempty,gt=0"`
	PageSize    *int      `form:"pageSize" validate:"required_with=CurrentPage,omitempty,gt=0"`
	OrderBy     []string  `form:"orderBy" validate:"order_by"`
}

type GetPlaylistFiltersMetadataRequest struct {
	SearchBy []string `form:"searchBy" validate:"search_by"`
}

type CreatePlaylistRequest struct {
	Title       string `validate:"required,max=100"`
	Description string
}

type AddAlbumsToPlaylistRequest struct {
	ID       uuid.UUID   `validate:"required"`
	AlbumIDs []uuid.UUID `validate:"min=1"`
	ForceAdd *bool
}

type AddArtistsToPlaylistRequest struct {
	ID        uuid.UUID   `validate:"required"`
	ArtistIDs []uuid.UUID `validate:"min=1"`
	ForceAdd  *bool
}

type AddSongsToPlaylistRequest struct {
	ID       uuid.UUID   `validate:"required"`
	SongIDs  []uuid.UUID `validate:"min=1"`
	ForceAdd *bool
}

type ShufflePlaylistSongsRequest struct {
	ID uuid.UUID `validate:"required"`
}

type UpdatePlaylistRequest struct {
	ID          uuid.UUID `validate:"required"`
	Title       string    `validate:"required,max=100"`
	Description string
}

type MoveSongFromPlaylistRequest struct {
	ID                 uuid.UUID `validate:"required"`
	PlaylistSongID     uuid.UUID `validate:"required"`
	OverPlaylistSongID uuid.UUID `validate:"required"`
}

type RemoveSongsFromPlaylistRequest struct {
	ID              uuid.UUID   `validate:"required"`
	PlaylistSongIDs []uuid.UUID `validate:"min=1"`
}
