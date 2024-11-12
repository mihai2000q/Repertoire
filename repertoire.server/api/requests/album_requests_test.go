package requests

import (
	"repertoire/server/api/validation"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidateGetAlbumsRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	tests := []struct {
		name    string
		request GetAlbumsRequest
	}{
		{
			"All Null",
			GetAlbumsRequest{},
		},
		{
			"Nothing Null",
			GetAlbumsRequest{
				CurrentPage: &[]int{1}[0],
				PageSize:    &[]int{1}[0],
				OrderBy:     []string{"title asc", "created_at desc"},
				SearchBy:    []string{"title = something", "is_recorded <> false"},
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

func TestValidateGetAlbumsRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              GetAlbumsRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// Current Page Test Cases
		{
			"Current Page is invalid because it should be greater than 0",
			GetAlbumsRequest{CurrentPage: &[]int{0}[0], PageSize: &[]int{1}[0]},
			"CurrentPage",
			"gt",
		},
		{
			"Current Page is invalid because page size is null",
			GetAlbumsRequest{PageSize: &[]int{1}[0]},
			"CurrentPage",
			"required_with",
		},
		// Page Size Test Cases
		{
			"Page Size is invalid because it should be greater than 0",
			GetAlbumsRequest{PageSize: &[]int{0}[0], CurrentPage: &[]int{1}[0]},
			"PageSize",
			"gt",
		},
		{
			"Page Size is invalid because current page is null",
			GetAlbumsRequest{CurrentPage: &[]int{1}[0]},
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
			assert.Contains(t, errCode.Error.Error(), "GetAlbumsRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, 400, errCode.Code)
		})
	}
}

var validAlbumTitle = "Justice For All"

func TestValidateCreateAlbumRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	tests := []struct {
		name    string
		request CreateAlbumRequest
	}{
		{
			"All Null",
			CreateAlbumRequest{
				Title: validAlbumTitle,
			},
		},
		{
			"All Filled With Existing Arist",
			CreateAlbumRequest{
				Title:       validAlbumTitle,
				ReleaseDate: &[]time.Time{time.Now()}[0],
				ArtistID:    &[]uuid.UUID{uuid.New()}[0],
			},
		},
		{
			"All Filled With New Arist",
			CreateAlbumRequest{
				Title:       validAlbumTitle,
				ReleaseDate: &[]time.Time{time.Now()}[0],
				ArtistName:  &[]string{"New Name"}[0],
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

func TestValidateCreateAlbumRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                  string
		request               CreateAlbumRequest
		expectedInvalidFields []string
		expectedFailedTags    []string
	}{
		// Title Test Cases
		{
			"Title is invalid because it's required",
			CreateAlbumRequest{Title: ""},
			[]string{"Title"},
			[]string{"required"},
		},
		{
			"Title is invalid because it has more than 100 characters",
			CreateAlbumRequest{Title: strings.Repeat("a", 101)},
			[]string{"Title"},
			[]string{"max"},
		},
		// Artist ID and Artist Title Test Case
		{
			"Artist Name and ID are invalid because only one can be set at a time",
			CreateAlbumRequest{
				Title:      validAlbumTitle,
				ArtistID:   &[]uuid.UUID{uuid.New()}[0],
				ArtistName: &[]string{"New Artist Name"}[0],
			},
			[]string{"ArtistID", "ArtistName"},
			[]string{"excluded_with", "excluded_with"},
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
			for _, expectedInvalidField := range tt.expectedInvalidFields {
				assert.Contains(t, errCode.Error.Error(), "CreateAlbumRequest."+expectedInvalidField)
			}
			for _, expectedFailedTag := range tt.expectedFailedTags {
				assert.Contains(t, errCode.Error.Error(), "'"+expectedFailedTag+"' tag")
			}
			assert.Equal(t, 400, errCode.Code)
		})
	}
}

func TestValidateAddSongToAlbumRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := AddSongToAlbumRequest{
		ID:     uuid.New(),
		SongID: uuid.New(),
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)
}

func TestValidateAddSongToAlbumRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              AddSongToAlbumRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// ID Test Cases
		{
			"ID is invalid because it's required",
			AddSongToAlbumRequest{ID: uuid.Nil, SongID: uuid.New()},
			"ID",
			"required",
		},
		// Song ID Test Cases
		{
			"Song ID is invalid because it's required",
			AddSongToAlbumRequest{ID: uuid.New(), SongID: uuid.Nil},
			"SongID",
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
			assert.Contains(t, errCode.Error.Error(), "AddSongToAlbumRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, 400, errCode.Code)
		})
	}
}

func TestValidateUpdateAlbumRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	tests := []struct {
		name    string
		request UpdateAlbumRequest
	}{
		{
			"All Null",
			UpdateAlbumRequest{
				ID:    uuid.New(),
				Title: validAlbumTitle,
			},
		},
		{
			"All Filled",
			UpdateAlbumRequest{
				ID:          uuid.New(),
				Title:       validAlbumTitle,
				ReleaseDate: &[]time.Time{time.Now()}[0],
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

func TestValidateUpdateAlbumRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              UpdateAlbumRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// ID Test Cases
		{
			"ID is invalid because it's required",
			UpdateAlbumRequest{ID: uuid.Nil, Title: validAlbumTitle},
			"ID",
			"required",
		},
		// Title Test Cases
		{
			"Title is invalid because it's required",
			UpdateAlbumRequest{ID: uuid.New(), Title: ""},
			"Title",
			"required",
		},
		{
			"Title is invalid because it has more than 100 characters",
			UpdateAlbumRequest{ID: uuid.New(), Title: strings.Repeat("a", 101)},
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
			assert.Contains(t, errCode.Error.Error(), "UpdateAlbumRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, 400, errCode.Code)
		})
	}
}

func TestValidateMoveSongFromAlbumRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := MoveSongFromAlbumRequest{
		ID:         uuid.New(),
		SongID:     uuid.New(),
		OverSongID: uuid.New(),
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)
}

func TestValidateMoveSongFromAlbumRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              MoveSongFromAlbumRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// ID Test Cases
		{
			"ID is invalid because it's required",
			MoveSongFromAlbumRequest{
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
			MoveSongFromAlbumRequest{
				ID:         uuid.New(),
				SongID:     uuid.Nil,
				OverSongID: uuid.New(),
			},
			"SongID",
			"required",
		},
		// Over Song ID Test Cases
		{
			"Song ID is invalid because it's required",
			MoveSongFromAlbumRequest{
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
			assert.Contains(t, errCode.Error.Error(), "MoveSongFromAlbumRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, 400, errCode.Code)
		})
	}
}
