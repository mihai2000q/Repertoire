package requests

import "github.com/google/uuid"

type GetArtistsRequest struct {
	CurrentPage *int     `form:"currentPage" validate:"required_with=PageSize,omitempty,gt=0"`
	PageSize    *int     `form:"pageSize" validate:"required_with=CurrentPage,omitempty,gt=0"`
	OrderBy     []string `form:"orderBy"`
	SearchBy    []string `form:"searchBy"`
}

type GetArtistFiltersMetadataRequest struct {
	SearchBy []string `form:"searchBy"`
}

type CreateArtistRequest struct {
	Name   string `validate:"required,max=100"`
	IsBand bool
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
	ID     uuid.UUID `validate:"required"`
	Name   string    `validate:"required,max=100"`
	IsBand bool
}

type RemoveAlbumsFromArtistRequest struct {
	ID       uuid.UUID   `validate:"required"`
	AlbumIDs []uuid.UUID `validate:"min=1"`
}

type RemoveSongsFromArtistRequest struct {
	ID      uuid.UUID   `validate:"required"`
	SongIDs []uuid.UUID `validate:"min=1"`
}

type DeleteArtistRequest struct {
	ID         uuid.UUID `validate:"required"`
	WithAlbums bool      `form:"withAlbums"`
	WithSongs  bool      `form:"withSongs"`
}

// Band Members

type CreateBandMemberRequest struct {
	Name     string      `validate:"required,max=100"`
	Color    *string     `validate:"omitempty,color"`
	RoleIDs  []uuid.UUID `validate:"min=1"`
	ArtistID uuid.UUID   `validate:"required"`
}

type UpdateBandMemberRequest struct {
	ID      uuid.UUID   `validate:"required"`
	Name    string      `validate:"required,max=100"`
	Color   *string     `validate:"omitempty,color"`
	RoleIDs []uuid.UUID `validate:"min=1"`
}

type MoveBandMemberRequest struct {
	ID       uuid.UUID `validate:"required"`
	OverID   uuid.UUID `validate:"required"`
	ArtistID uuid.UUID `validate:"required"`
}
