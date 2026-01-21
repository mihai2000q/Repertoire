package requests

import (
	"repertoire/server/internal"

	"github.com/google/uuid"
)

type GetAlbumRequest struct {
	ID           uuid.UUID `validate:"required"`
	SongsOrderBy []string  `form:"songsOrderBy" validate:"order_by"`
}

type GetAlbumsRequest struct {
	CurrentPage *int     `form:"currentPage" validate:"required_with=PageSize,omitempty,gt=0"`
	PageSize    *int     `form:"pageSize" validate:"required_with=CurrentPage,omitempty,gt=0"`
	OrderBy     []string `form:"orderBy" validate:"order_by"`
	SearchBy    []string `form:"searchBy" validate:"search_by"`
}

type GetAlbumFiltersMetadataRequest struct {
	SearchBy []string `form:"searchBy" validate:"search_by"`
}

type CreateAlbumRequest struct {
	Title       string `validate:"required,max=100"`
	ReleaseDate *internal.Date
	ArtistID    *uuid.UUID `validate:"omitempty,excluded_with=ArtistName"`
	ArtistName  *string    `validate:"omitempty,excluded_with=ArtistID,max=100"`
}

type AddSongsToAlbumRequest struct {
	ID      uuid.UUID   `validate:"required"`
	SongIDs []uuid.UUID `validate:"min=1"`
}

type AddPerfectRehearsalsToAlbumsRequest struct {
	IDs []uuid.UUID `validate:"min=1"`
}

type UpdateAlbumRequest struct {
	ID          uuid.UUID `validate:"required"`
	Title       string    `validate:"required,max=100"`
	ReleaseDate *internal.Date
	ArtistID    *uuid.UUID
}

type MoveSongFromAlbumRequest struct {
	ID         uuid.UUID `validate:"required"`
	SongID     uuid.UUID `validate:"required"`
	OverSongID uuid.UUID `validate:"required"`
}

type RemoveSongsFromAlbumRequest struct {
	ID      uuid.UUID   `validate:"required"`
	SongIDs []uuid.UUID `validate:"min=1"`
}

type BulkDeleteAlbumsRequest struct {
	IDs       []uuid.UUID `validate:"min=1"`
	WithSongs bool
}

type DeleteAlbumRequest struct {
	ID        uuid.UUID `validate:"required"`
	WithSongs bool      `form:"withSongs"`
}
