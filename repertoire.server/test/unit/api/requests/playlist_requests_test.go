package requests

import (
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/api/validation"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidateGetPlaylistsRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	tests := []struct {
		name    string
		request requests.GetPlaylistsRequest
	}{
		{
			"All Null",
			requests.GetPlaylistsRequest{},
		},
		{
			"Nothing Null",
			requests.GetPlaylistsRequest{
				CurrentPage: &[]int{1}[0],
				PageSize:    &[]int{1}[0],
				OrderBy:     []string{"title asc"},
				SearchBy:    []string{"title = something", "title is not null"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			_uut := validation.NewValidator(nil)

			// when
			errCode := _uut.Validate(tt.request)

			// then
			assert.Nil(t, errCode)
		})
	}
}

func TestValidateGetPlaylistsRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.GetPlaylistsRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// Current Page Test Cases
		{
			"Current Page is invalid because it should be greater than 0",
			requests.GetPlaylistsRequest{CurrentPage: &[]int{0}[0], PageSize: &[]int{1}[0]},
			"CurrentPage",
			"gt",
		},
		{
			"Current Page is invalid because page size is null",
			requests.GetPlaylistsRequest{PageSize: &[]int{1}[0]},
			"CurrentPage",
			"required_with",
		},
		// Page Size Test Cases
		{
			"Page Size is invalid because it should be greater than 0",
			requests.GetPlaylistsRequest{PageSize: &[]int{0}[0], CurrentPage: &[]int{1}[0]},
			"PageSize",
			"gt",
		},
		{
			"Page Size is invalid because current page is null",
			requests.GetPlaylistsRequest{CurrentPage: &[]int{1}[0]},
			"PageSize",
			"required_with",
		},
		// Search By Test Cases
		{
			"Search By is invalid because the operator is not supported",
			requests.GetPlaylistsRequest{SearchBy: []string{"title = okay", "songs is nullish"}},
			"SearchBy",
			"search_by",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			_uut := validation.NewValidator(nil)

			// when
			errCode := _uut.Validate(tt.request)

			// then
			assert.NotNil(t, errCode)
			assert.Len(t, errCode.Error, 1)
			assert.Contains(t, errCode.Error.Error(), "GetPlaylistsRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateGetPlaylistFiltersMetadataRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	tests := []struct {
		name    string
		request requests.GetPlaylistFiltersMetadataRequest
	}{
		{
			"All Null",
			requests.GetPlaylistFiltersMetadataRequest{},
		},
		{
			"Nothing Null",
			requests.GetPlaylistFiltersMetadataRequest{
				SearchBy: []string{"title = something", "title is not null"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			_uut := validation.NewValidator(nil)

			// when
			errCode := _uut.Validate(tt.request)

			// then
			assert.Nil(t, errCode)
		})
	}
}

func TestValidateGetPlaylistFiltersMetadataRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.GetPlaylistFiltersMetadataRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// Search By Test Cases
		{
			"Search By is invalid because the operator is not supported",
			requests.GetPlaylistFiltersMetadataRequest{SearchBy: []string{"title = okay", "songs is nullish"}},
			"SearchBy",
			"search_by",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			_uut := validation.NewValidator(nil)

			// when
			errCode := _uut.Validate(tt.request)

			// then
			assert.NotNil(t, errCode)
			assert.Len(t, errCode.Error, 1)
			assert.Contains(t, errCode.Error.Error(), "GetPlaylistFiltersMetadataRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

var validPlaylistTitle = "New Playlist"

func TestValidateCreatePlaylistRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.CreatePlaylistRequest{
		Title: validPlaylistTitle,
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)
}

func TestValidateCreatePlaylistRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.CreatePlaylistRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// Title Test Cases
		{
			"Title is invalid because it's required",
			requests.CreatePlaylistRequest{Title: ""},
			"Title",
			"required",
		},
		{
			"Title is invalid because it has more than 100 characters",
			requests.CreatePlaylistRequest{Title: strings.Repeat("a", 101)},
			"Title",
			"max",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			_uut := validation.NewValidator(nil)

			// when
			errCode := _uut.Validate(tt.request)

			// then
			assert.NotNil(t, errCode)
			assert.Len(t, errCode.Error, 1)
			assert.Contains(t, errCode.Error.Error(), "CreatePlaylistRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateAddAlbumsToPlaylistRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.AddAlbumsToPlaylistRequest{
		ID:       uuid.New(),
		AlbumIDs: []uuid.UUID{uuid.New()},
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)
}

func TestValidateAddAlbumsToPlaylistRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.AddAlbumsToPlaylistRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// ID Test Cases
		{
			"ID is invalid because it's required",
			requests.AddAlbumsToPlaylistRequest{ID: uuid.Nil, AlbumIDs: []uuid.UUID{uuid.New()}},
			"ID",
			"required",
		},
		// Album IDs Test Cases
		{
			"Album IDs is invalid because it requires at least 1 ID",
			requests.AddAlbumsToPlaylistRequest{ID: uuid.New(), AlbumIDs: []uuid.UUID{}},
			"AlbumIDs",
			"min",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			_uut := validation.NewValidator(nil)

			// when
			errCode := _uut.Validate(tt.request)

			// then
			assert.NotNil(t, errCode)
			assert.Len(t, errCode.Error, 1)
			assert.Contains(t, errCode.Error.Error(), "AddAlbumsToPlaylistRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateAddArtistsToPlaylistRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.AddArtistsToPlaylistRequest{
		ID:        uuid.New(),
		ArtistIDs: []uuid.UUID{uuid.New()},
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)
}

func TestValidateAddArtistsToPlaylistRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.AddArtistsToPlaylistRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// ID Test Cases
		{
			"ID is invalid because it's required",
			requests.AddArtistsToPlaylistRequest{ID: uuid.Nil, ArtistIDs: []uuid.UUID{uuid.New()}},
			"ID",
			"required",
		},
		// Artist IDs Test Cases
		{
			"Album IDs is invalid because it requires at least 1 ID",
			requests.AddArtistsToPlaylistRequest{ID: uuid.New(), ArtistIDs: []uuid.UUID{}},
			"ArtistIDs",
			"min",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			_uut := validation.NewValidator(nil)

			// when
			errCode := _uut.Validate(tt.request)

			// then
			assert.NotNil(t, errCode)
			assert.Len(t, errCode.Error, 1)
			assert.Contains(t, errCode.Error.Error(), "AddArtistsToPlaylistRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateAddSongsToPlaylistRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.AddSongsToPlaylistRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New()},
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)
}

func TestValidateAddSongsToPlaylistRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.AddSongsToPlaylistRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// ID Test Cases
		{
			"ID is invalid because it's required",
			requests.AddSongsToPlaylistRequest{ID: uuid.Nil, SongIDs: []uuid.UUID{uuid.New()}},
			"ID",
			"required",
		},
		// Song IDs Test Cases
		{
			"Song IDs is invalid because it requires at least 1 ID",
			requests.AddSongsToPlaylistRequest{ID: uuid.New(), SongIDs: []uuid.UUID{}},
			"SongIDs",
			"min",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			_uut := validation.NewValidator(nil)

			// when
			errCode := _uut.Validate(tt.request)

			// then
			assert.NotNil(t, errCode)
			assert.Len(t, errCode.Error, 1)
			assert.Contains(t, errCode.Error.Error(), "AddSongsToPlaylistRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateUpdatePlaylistRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.UpdatePlaylistRequest{
		ID:    uuid.New(),
		Title: validPlaylistTitle,
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)
}

func TestValidateUpdatePlaylistRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.UpdatePlaylistRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// ID Test Cases
		{
			"ID is invalid because it's required",
			requests.UpdatePlaylistRequest{ID: uuid.Nil, Title: validPlaylistTitle},
			"ID",
			"required",
		},
		// Title Test Cases
		{
			"Title is invalid because it's required",
			requests.UpdatePlaylistRequest{ID: uuid.New(), Title: ""},
			"Title",
			"required",
		},
		{
			"Title is invalid because it has more than 100 characters",
			requests.UpdatePlaylistRequest{ID: uuid.New(), Title: strings.Repeat("a", 101)},
			"Title",
			"max",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			_uut := validation.NewValidator(nil)

			// when
			errCode := _uut.Validate(tt.request)

			// then
			assert.NotNil(t, errCode)
			assert.Len(t, errCode.Error, 1)
			assert.Contains(t, errCode.Error.Error(), "UpdatePlaylistRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateMoveSongFromPlaylistRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.MoveSongFromPlaylistRequest{
		ID:         uuid.New(),
		SongID:     uuid.New(),
		OverSongID: uuid.New(),
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)
}

func TestValidateMoveSongFromPlaylistRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.MoveSongFromPlaylistRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// ID Test Cases
		{
			"ID is invalid because it's required",
			requests.MoveSongFromPlaylistRequest{
				ID:         uuid.Nil,
				SongID:     uuid.New(),
				OverSongID: uuid.New(),
			},
			"ID",
			"required",
		},
		// Song ID Test Cases
		{
			"Song ID is invalid because it's required",
			requests.MoveSongFromPlaylistRequest{
				ID:         uuid.New(),
				SongID:     uuid.Nil,
				OverSongID: uuid.New(),
			},
			"SongID",
			"required",
		},
		// Over Song ID Test Cases
		{
			"Over Song ID is invalid because it's required",
			requests.MoveSongFromPlaylistRequest{
				ID:         uuid.New(),
				SongID:     uuid.New(),
				OverSongID: uuid.Nil,
			},
			"OverSongID",
			"required",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			_uut := validation.NewValidator(nil)

			// when
			errCode := _uut.Validate(tt.request)

			// then
			assert.NotNil(t, errCode)
			assert.Len(t, errCode.Error, 1)
			assert.Contains(t, errCode.Error.Error(), "MoveSongFromPlaylistRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateRemoveSongsFromPlaylistRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.RemoveSongsFromPlaylistRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New()},
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)
}

func TestValidateRemoveSongsFromPlaylistRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.RemoveSongsFromPlaylistRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// ID Test Cases
		{
			"ID is invalid because it's required",
			requests.RemoveSongsFromPlaylistRequest{ID: uuid.Nil, SongIDs: []uuid.UUID{uuid.New()}},
			"ID",
			"required",
		},
		// Song IDs Test Cases
		{
			"Song IDs is invalid because it requires at least 1 ID",
			requests.RemoveSongsFromPlaylistRequest{ID: uuid.New(), SongIDs: []uuid.UUID{}},
			"SongIDs",
			"min",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			_uut := validation.NewValidator(nil)

			// when
			errCode := _uut.Validate(tt.request)

			// then
			assert.NotNil(t, errCode)
			assert.Len(t, errCode.Error, 1)
			assert.Contains(t, errCode.Error.Error(), "RemoveSongsFromPlaylistRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}
