package requests

import "github.com/google/uuid"

type GetArtistsRequest struct {
	CurrentPage *int     `form:"currentPage" validate:"required_with=PageSize,omitempty,gt=0"`
	PageSize    *int     `form:"pageSize" validate:"required_with=CurrentPage,omitempty,gt=0"`
	OrderBy     []string `form:"orderBy"`
	SearchBy    []string `form:"orderBy"`
}

type CreateArtistRequest struct {
	Name string `validate:"required,max=100"`
}

type AddAlbumsToArtistRequest struct {
	ID       uuid.UUID   `validate:"required"`
	AlbumIDs []uuid.UUID `validate:"min=1"`
}

type AddSongsToArtistRequest struct {
	ID      uuid.UUID   `validate:"required"`
	SongIDs []uuid.UUID `validate:"min=1"`
}

type UpdateArtistRequest struct {
	ID   uuid.UUID `validate:"required"`
	Name string    `validate:"required,max=100"`
}

type RemoveAlbumsFromArtistRequest struct {
	ID       uuid.UUID   `validate:"required"`
	AlbumIDs []uuid.UUID `validate:"min=1"`
}

type RemoveSongsFromArtistRequest struct {
	ID      uuid.UUID   `validate:"required"`
	SongIDs []uuid.UUID `validate:"min=1"`
}
