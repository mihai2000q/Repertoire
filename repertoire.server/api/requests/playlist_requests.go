package requests

import "github.com/google/uuid"

type GetPlaylistRequest struct {
	ID           uuid.UUID `validate:"required"`
	SongsOrderBy []string  `form:"songsOrderBy" validate:"order_by"`
}

type GetPlaylistsRequest struct {
	CurrentPage *int     `form:"currentPage" validate:"required_with=PageSize,omitempty,gt=0"`
	PageSize    *int     `form:"pageSize" validate:"required_with=CurrentPage,omitempty,gt=0"`
	OrderBy     []string `form:"orderBy" validate:"order_by"`
	SearchBy    []string `form:"searchBy" validate:"search_by"`
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
}

type AddArtistsToPlaylistRequest struct {
	ID        uuid.UUID   `validate:"required"`
	ArtistIDs []uuid.UUID `validate:"min=1"`
}

type AddSongsToPlaylistRequest struct {
	ID      uuid.UUID   `validate:"required"`
	SongIDs []uuid.UUID `validate:"min=1"`
}

type UpdatePlaylistRequest struct {
	ID          uuid.UUID `validate:"required"`
	Title       string    `validate:"required,max=100"`
	Description string
}

type MoveSongFromPlaylistRequest struct {
	ID         uuid.UUID `validate:"required"`
	SongID     uuid.UUID `validate:"required"`
	OverSongID uuid.UUID `validate:"required"`
}

type RemoveSongsFromPlaylistRequest struct {
	ID      uuid.UUID   `validate:"required"`
	SongIDs []uuid.UUID `validate:"min=1"`
}
