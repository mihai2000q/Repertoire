package requests

import (
	"repertoire/api/validation"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidateGetArtistsRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	tests := []struct {
		name    string
		request GetArtistsRequest
	}{
		{
			"All Null",
			GetArtistsRequest{
				UserID: uuid.New(),
			},
		},
		{
			"Nothing Null",
			GetArtistsRequest{
				UserID:      uuid.New(),
				CurrentPage: &[]int{1}[0],
				PageSize:    &[]int{1}[0],
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
		request              GetArtistsRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// User ID Test Cases
		{
			"User ID is invalid because it's required",
			GetArtistsRequest{UserID: uuid.Nil},
			"UserID",
			"required",
		},
		// Current Page Test Cases
		{
			"Current Page is invalid because it should be greater than 0",
			GetArtistsRequest{UserID: uuid.New(), CurrentPage: &[]int{0}[0], PageSize: &[]int{1}[0]},
			"CurrentPage",
			"gt",
		},
		{
			"Current Page is invalid because page size is null",
			GetArtistsRequest{UserID: uuid.New(), PageSize: &[]int{1}[0]},
			"CurrentPage",
			"required_with",
		},
		// Page Size Test Cases
		{
			"Page Size is invalid because it should be greater than 0",
			GetArtistsRequest{UserID: uuid.New(), PageSize: &[]int{0}[0], CurrentPage: &[]int{1}[0]},
			"PageSize",
			"gt",
		},
		{
			"Page Size is invalid because current page is null",
			GetArtistsRequest{UserID: uuid.New(), CurrentPage: &[]int{1}[0]},
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
			assert.Contains(t, err, "GetArtistsRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, 400, errCode.Code)
		})
	}
}

var validArtistName = "Metallica"

func TestValidateCreateArtistRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	_uut := validation.NewValidator(nil)

	request := CreateArtistRequest{
		Name: validArtistName,
	}

	errCode := _uut.Validate(request)

	assert.Nil(t, errCode)
}

func TestValidateCreateArtistRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              CreateArtistRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// Name Test Cases
		{
			"Name is invalid because it's required",
			CreateArtistRequest{Name: ""},
			"Name",
			"required",
		},
		{
			"Name is invalid because it has more than 100 characters",
			CreateArtistRequest{Name: strings.Repeat("a", 101)},
			"Name",
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
			assert.Contains(t, err, "CreateArtistRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, 400, errCode.Code)
		})
	}
}

func TestValidateUpdateArtistRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	_uut := validation.NewValidator(nil)

	request := UpdateArtistRequest{
		ID:   uuid.New(),
		Name: validArtistName,
	}

	errCode := _uut.Validate(request)

	assert.Nil(t, errCode)
}

func TestValidateUpdateArtistRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              UpdateArtistRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// ID Test Cases
		{
			"ID is invalid because it's required",
			UpdateArtistRequest{ID: uuid.Nil, Name: validArtistName},
			"ID",
			"required",
		},
		// Name Test Cases
		{
			"Name is invalid because it's required",
			UpdateArtistRequest{ID: uuid.New(), Name: ""},
			"Name",
			"required",
		},
		{
			"Name is invalid because it has more than 100 characters",
			UpdateArtistRequest{ID: uuid.New(), Name: strings.Repeat("a", 101)},
			"Name",
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
			assert.Contains(t, err, "UpdateArtistRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, 400, errCode.Code)
		})
	}
}
