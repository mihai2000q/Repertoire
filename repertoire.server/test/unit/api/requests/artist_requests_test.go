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

func TestValidateGetArtistsRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	tests := []struct {
		name    string
		request requests.GetArtistsRequest
	}{
		{
			"All Null",
			requests.GetArtistsRequest{},
		},
		{
			"Nothing Null",
			requests.GetArtistsRequest{
				CurrentPage: &[]int{1}[0],
				PageSize:    &[]int{1}[0],
				OrderBy:     []string{"name asc"},
				SearchBy:    []string{"name = Metallica"},
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

func TestValidateGetArtistsRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.GetArtistsRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// Current Page Test Cases
		{
			"Current Page is invalid because it should be greater than 0",
			requests.GetArtistsRequest{CurrentPage: &[]int{0}[0], PageSize: &[]int{1}[0]},
			"CurrentPage",
			"gt",
		},
		{
			"Current Page is invalid because page size is null",
			requests.GetArtistsRequest{PageSize: &[]int{1}[0]},
			"CurrentPage",
			"required_with",
		},
		// Page Size Test Cases
		{
			"Page Size is invalid because it should be greater than 0",
			requests.GetArtistsRequest{PageSize: &[]int{0}[0], CurrentPage: &[]int{1}[0]},
			"PageSize",
			"gt",
		},
		{
			"Page Size is invalid because current page is null",
			requests.GetArtistsRequest{CurrentPage: &[]int{1}[0]},
			"PageSize",
			"required_with",
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
			assert.Contains(t, errCode.Error.Error(), "GetArtistsRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

var validArtistName = "Metallica"

func TestValidateCreateArtistRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.CreateArtistRequest{
		Name: validArtistName,
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)
}

func TestValidateCreateArtistRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.CreateArtistRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// Name Test Cases
		{
			"Name is invalid because it's required",
			requests.CreateArtistRequest{Name: ""},
			"Name",
			"required",
		},
		{
			"Name is invalid because it has more than 100 characters",
			requests.CreateArtistRequest{Name: strings.Repeat("a", 101)},
			"Name",
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
			assert.Contains(t, errCode.Error.Error(), "CreateArtistRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateAddAlbumsToArtistRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.AddAlbumsToArtistRequest{
		ID:       uuid.New(),
		AlbumIDs: []uuid.UUID{uuid.New()},
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)
}

func TestValidateAddAlbumsToArtistRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.AddAlbumsToArtistRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// ID Test Cases
		{
			"ID is invalid because it's required",
			requests.AddAlbumsToArtistRequest{ID: uuid.Nil, AlbumIDs: []uuid.UUID{uuid.New()}},
			"ID",
			"required",
		},
		// Album IDs Test Cases
		{
			"Album IDs is invalid because it requires at least one ID",
			requests.AddAlbumsToArtistRequest{ID: uuid.New(), AlbumIDs: []uuid.UUID{}},
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
			assert.Contains(t, errCode.Error.Error(), "AddAlbumsToArtistRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateAddSongsToArtistRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.AddSongsToArtistRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New()},
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)
}

func TestValidateAddSongsToArtistRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.AddSongsToArtistRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// ID Test Cases
		{
			"ID is invalid because it's required",
			requests.AddSongsToArtistRequest{ID: uuid.Nil, SongIDs: []uuid.UUID{uuid.New()}},
			"ID",
			"required",
		},
		// Song IDs Test Cases
		{
			"Song IDs is invalid because it requires at least one ID",
			requests.AddSongsToArtistRequest{ID: uuid.New(), SongIDs: []uuid.UUID{}},
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
			assert.Contains(t, errCode.Error.Error(), "AddSongsToArtistRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateUpdateArtistRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.UpdateArtistRequest{
		ID:   uuid.New(),
		Name: validArtistName,
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)
}

func TestValidateUpdateArtistRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.UpdateArtistRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// ID Test Cases
		{
			"ID is invalid because it's required",
			requests.UpdateArtistRequest{ID: uuid.Nil, Name: validArtistName},
			"ID",
			"required",
		},
		// Name Test Cases
		{
			"Name is invalid because it's required",
			requests.UpdateArtistRequest{ID: uuid.New(), Name: ""},
			"Name",
			"required",
		},
		{
			"Name is invalid because it has more than 100 characters",
			requests.UpdateArtistRequest{ID: uuid.New(), Name: strings.Repeat("a", 101)},
			"Name",
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
			assert.Contains(t, errCode.Error.Error(), "UpdateArtistRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateRemoveAlbumsFromArtistRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.RemoveAlbumsFromArtistRequest{
		ID:       uuid.New(),
		AlbumIDs: []uuid.UUID{uuid.New()},
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)
}

func TestValidateRemoveAlbumsFromArtistRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.RemoveAlbumsFromArtistRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// ID Test Cases
		{
			"ID is invalid because it's required",
			requests.RemoveAlbumsFromArtistRequest{ID: uuid.Nil, AlbumIDs: []uuid.UUID{uuid.New()}},
			"ID",
			"required",
		},
		// Album IDs Test Cases
		{
			"Album IDs is invalid because it requires at least one ID",
			requests.RemoveAlbumsFromArtistRequest{ID: uuid.New(), AlbumIDs: []uuid.UUID{}},
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
			assert.Contains(t, errCode.Error.Error(), "RemoveAlbumsFromArtistRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateRemoveSongsFromArtistRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.RemoveSongsFromArtistRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New()},
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)
}

func TestValidateRemoveSongsFromArtistRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.RemoveSongsFromArtistRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// ID Test Cases
		{
			"ID is invalid because it's required",
			requests.RemoveSongsFromArtistRequest{ID: uuid.Nil, SongIDs: []uuid.UUID{uuid.New()}},
			"ID",
			"required",
		},
		// Song IDs Test Cases
		{
			"Song IDs is invalid because it requires at least one ID",
			requests.RemoveSongsFromArtistRequest{ID: uuid.New(), SongIDs: []uuid.UUID{}},
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
			assert.Contains(t, errCode.Error.Error(), "RemoveSongsFromArtistRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateDeleteArtistRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	tests := []struct {
		name    string
		request requests.DeleteArtistRequest
	}{
		{
			"Minimal",
			requests.DeleteArtistRequest{ID: uuid.New()},
		},
		{
			"Maximal",
			requests.DeleteArtistRequest{
				ID:         uuid.New(),
				WithSongs:  true,
				WithAlbums: true,
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

func TestValidateDeleteArtistRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.DeleteArtistRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// ID Cases
		{
			"ID is invalid because it is required",
			requests.DeleteArtistRequest{ID: uuid.Nil},
			"ID",
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
			assert.Contains(t, errCode.Error.Error(), "DeleteArtistRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}
