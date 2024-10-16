package requests

import (
	"repertoire/api/validation"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidateGetPlaylistsRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	_uut := validation.NewValidator(nil)

	request := GetPlaylistsRequest{
		UserID: uuid.New(),
	}

	errCode := _uut.Validate(request)

	assert.Nil(t, errCode)
}

func TestValidateGetPlaylistsRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              GetPlaylistsRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// User ID Test Cases
		{
			"User ID is invalid because it's required",
			GetPlaylistsRequest{UserID: uuid.Nil},
			"UserID",
			"required",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_uut := validation.NewValidator(nil)

			errCode := _uut.Validate(tt.request)

			err := errCode.Error.Error()

			assert.NotNil(t, errCode)
			assert.Len(t, errCode.Error, 1)
			assert.Contains(t, err, "GetPlaylistsRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, 400, errCode.Code)
		})
	}
}

var validPlaylistTitle = "New Playlist"

func TestValidateCreatePlaylistRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	_uut := validation.NewValidator(nil)

	request := CreatePlaylistRequest{
		Title: validPlaylistTitle,
	}

	errCode := _uut.Validate(request)

	assert.Nil(t, errCode)
}

func TestValidateCreatePlaylistRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              CreatePlaylistRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// Title Test Cases
		{
			"Title is invalid because it's required",
			CreatePlaylistRequest{Title: ""},
			"Title",
			"required",
		},
		{
			"Title is invalid because it has more than 100 characters",
			CreatePlaylistRequest{Title: strings.Repeat("a", 101)},
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
			assert.Contains(t, err, "CreatePlaylistRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, 400, errCode.Code)
		})
	}
}

func TestValidateUpdatePlaylistRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	_uut := validation.NewValidator(nil)

	request := UpdatePlaylistRequest{
		ID:    uuid.New(),
		Title: validPlaylistTitle,
	}

	errCode := _uut.Validate(request)

	assert.Nil(t, errCode)
}

func TestValidateUpdatePlaylistRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              UpdatePlaylistRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// ID Test Cases
		{
			"ID is invalid because it's required",
			UpdatePlaylistRequest{ID: uuid.Nil, Title: validPlaylistTitle},
			"ID",
			"required",
		},
		// Title Test Cases
		{
			"Title is invalid because it's required",
			UpdatePlaylistRequest{ID: uuid.New(), Title: ""},
			"Title",
			"required",
		},
		{
			"Title is invalid because it has more than 100 characters",
			UpdatePlaylistRequest{ID: uuid.New(), Title: strings.Repeat("a", 101)},
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
			assert.Contains(t, err, "UpdatePlaylistRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, 400, errCode.Code)
		})
	}
}
