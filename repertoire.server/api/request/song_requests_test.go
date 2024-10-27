package request

import (
	"repertoire/api/validation"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidateGetSongsRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	tests := []struct {
		name    string
		request GetSongsRequest
	}{
		{
			"All Null",
			GetSongsRequest{},
		},
		{
			"Nothing Null",
			GetSongsRequest{
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

func TestValidateGetSongsRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              GetSongsRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// Current Page Test Cases
		{
			"Current Page is invalid because it should be greater than 0",
			GetSongsRequest{CurrentPage: &[]int{0}[0], PageSize: &[]int{1}[0]},
			"CurrentPage",
			"gt",
		},
		{
			"Current Page is invalid because page size is null",
			GetSongsRequest{PageSize: &[]int{1}[0]},
			"CurrentPage",
			"required_with",
		},
		// Page Size Test Cases
		{
			"Page Size is invalid because it should be greater than 0",
			GetSongsRequest{PageSize: &[]int{0}[0], CurrentPage: &[]int{1}[0]},
			"PageSize",
			"gt",
		},
		{
			"Page Size is invalid because current page is null",
			GetSongsRequest{CurrentPage: &[]int{1}[0]},
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
			assert.Contains(t, err, "GetSongsRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, 400, errCode.Code)
		})
	}
}

var validSongTitle = "Justice For All"

func TestValidateCreateSongRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	tests := []struct {
		name    string
		request CreateSongRequest
	}{
		{
			"All Null",
			CreateSongRequest{Title: validSongTitle},
		},
		{
			"Nothing Null",
			CreateSongRequest{
				Title:         validSongTitle,
				Bpm:           &[]uint{12}[0],
				SongsterrLink: &[]string{"http://songsterr.com/some-song"}[0],
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_uut := validation.NewValidator(nil)

			errCode := _uut.Validate(tt.request)

			assert.Nil(t, errCode)
		})
	}
}

func TestValidateCreateSongRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              CreateSongRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// Title Test Cases
		{
			"Title is invalid because it's required",
			CreateSongRequest{Title: ""},
			"Title",
			"required",
		},
		{
			"Title is invalid because it has more than 100 characters",
			CreateSongRequest{Title: strings.Repeat("a", 101)},
			"Title",
			"max",
		},
		// SongsterrLink Test Cases
		{
			"Songsterr Link is invalid because it is not an url",
			CreateSongRequest{
				Title:         validSongTitle,
				SongsterrLink: &[]string{"somethingcom"}[0],
			},
			"SongsterrLink",
			"url",
		},
		{
			"Songsterr Link is invalid because it is not a songsterr link",
			CreateSongRequest{
				Title:         validSongTitle,
				SongsterrLink: &[]string{"http://google.com"}[0],
			},
			"SongsterrLink",
			"contains",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_uut := validation.NewValidator(nil)

			errCode := _uut.Validate(tt.request)

			err := errCode.Error.Error()

			assert.NotNil(t, errCode)
			assert.Len(t, errCode.Error, 1)
			assert.Contains(t, err, "CreateSongRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, 400, errCode.Code)
		})
	}
}

func TestValidateUpdateSongRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	tests := []struct {
		name    string
		request UpdateSongRequest
	}{
		{
			"All Null",
			UpdateSongRequest{
				ID:    uuid.New(),
				Title: validSongTitle,
			},
		},
		{
			"Nothing Null",
			UpdateSongRequest{
				ID:            uuid.New(),
				Title:         validSongTitle,
				IsRecorded:    true,
				Reharsals:     12,
				Bpm:           &[]uint{120}[0],
				SongsterrLink: &[]string{"http://songsterr.com/some-song"}[0],
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_uut := validation.NewValidator(nil)

			errCode := _uut.Validate(tt.request)

			assert.Nil(t, errCode)
		})
	}
}

func TestValidateUpdateSongRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              UpdateSongRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// ID Test Cases
		{
			"ID is invalid because it's required",
			UpdateSongRequest{ID: uuid.Nil, Title: validSongTitle},
			"ID",
			"required",
		},
		// Title Test Cases
		{
			"Title is invalid because it's required",
			UpdateSongRequest{ID: uuid.New(), Title: ""},
			"Title",
			"required",
		},
		{
			"Title is invalid because it has more than 100 characters",
			UpdateSongRequest{ID: uuid.New(), Title: strings.Repeat("a", 101)},
			"Title",
			"max",
		},
		// SongsterrLink Test Cases
		{
			"Songsterr Link is invalid because it is not an url",
			UpdateSongRequest{
				ID:            uuid.New(),
				Title:         validSongTitle,
				SongsterrLink: &[]string{"somethingcom"}[0],
			},
			"SongsterrLink",
			"url",
		},
		{
			"Songsterr Link is invalid because it is not a songsterr link",
			UpdateSongRequest{
				ID:            uuid.New(),
				Title:         validSongTitle,
				SongsterrLink: &[]string{"http://google.com"}[0],
			},
			"SongsterrLink",
			"contains",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			_uut := validation.NewValidator(nil)

			errCode := _uut.Validate(tt.request)

			// when
			err := errCode.Error.Error()

			// then
			assert.NotNil(t, errCode)
			assert.Len(t, errCode.Error, 1)
			assert.Contains(t, err, "UpdateSongRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, 400, errCode.Code)
		})
	}
}
