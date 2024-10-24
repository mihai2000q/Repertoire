package request

import (
	"repertoire/api/validation"
	"strings"
	"testing"

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
				OrderBy:     "title asc, created_at desc",
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
			_uut := validation.NewValidator(nil)

			errCode := _uut.Validate(tt.request)

			err := errCode.Error.Error()

			assert.NotNil(t, errCode)
			assert.Len(t, errCode.Error, 1)
			assert.Contains(t, err, "GetAlbumsRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, 400, errCode.Code)
		})
	}
}

var validAlbumTitle = "Justice For All"

func TestValidateCreateAlbumRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	_uut := validation.NewValidator(nil)

	request := CreateAlbumRequest{
		Title: validAlbumTitle,
	}

	errCode := _uut.Validate(request)

	assert.Nil(t, errCode)
}

func TestValidateCreateAlbumRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              CreateAlbumRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// Title Test Cases
		{
			"Title is invalid because it's required",
			CreateAlbumRequest{Title: ""},
			"Title",
			"required",
		},
		{
			"Title is invalid because it has more than 100 characters",
			CreateAlbumRequest{Title: strings.Repeat("a", 101)},
			"Title",
			"max",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_uut := validation.NewValidator(nil)

			errCode := _uut.Validate(tt.request)

			err := errCode.Error.Error()

			assert.NotNil(t, errCode)
			assert.Len(t, errCode.Error, 1)
			assert.Contains(t, err, "CreateAlbumRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, 400, errCode.Code)
		})
	}
}

func TestValidateUpdateAlbumRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	_uut := validation.NewValidator(nil)

	request := UpdateAlbumRequest{
		ID:    uuid.New(),
		Title: validAlbumTitle,
	}

	errCode := _uut.Validate(request)

	assert.Nil(t, errCode)
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
			_uut := validation.NewValidator(nil)

			errCode := _uut.Validate(tt.request)

			err := errCode.Error.Error()

			assert.NotNil(t, errCode)
			assert.Len(t, errCode.Error, 1)
			assert.Contains(t, err, "UpdateAlbumRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, 400, errCode.Code)
		})
	}
}
