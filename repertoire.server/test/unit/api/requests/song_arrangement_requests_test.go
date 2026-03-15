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

var validArrangementName = "Perfect Rehearsal"

func TestValidateGetSongArrangementsRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.GetSongArrangementsRequest{
		SongID: uuid.New(),
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)
}

func TestValidateGetSongArrangementsRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.GetSongArrangementsRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// Song ID Test Cases
		{
			"Song ID is invalid because it's required",
			requests.GetSongArrangementsRequest{
				SongID: uuid.Nil,
			},
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
			assert.Contains(t, errCode.Error.Error(), "GetSongArrangementsRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateCreateSongArrangementRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.CreateSongArrangementRequest{
		SongID: uuid.New(),
		Name:   validArrangementName,
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)
}

func TestValidateCreateSongArrangementRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.CreateSongArrangementRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// Song ID Test Cases
		{
			"Song ID is invalid because it's required",
			requests.CreateSongArrangementRequest{
				SongID: uuid.Nil,
				Name:   validArrangementName,
			},
			"SongID",
			"required",
		},
		// Name Test Cases
		{
			"Name is invalid because it's required",
			requests.CreateSongArrangementRequest{
				SongID: uuid.New(),
				Name:   "",
			},
			"Name",
			"required",
		},
		{
			"Name is invalid because it has too many characters",
			requests.CreateSongArrangementRequest{
				SongID: uuid.New(),
				Name:   strings.Repeat("a", 31),
			},
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
			assert.Contains(t, errCode.Error.Error(), "CreateSongArrangementRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateBulkUpdateSongArrangementsRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	tests := []struct {
		name    string
		request requests.BulkUpdateSongArrangementsRequest
	}{
		{
			"Minimal",
			requests.BulkUpdateSongArrangementsRequest{
				SongID: uuid.New(),
				Requests: []requests.UpdateSongArrangementRequest{
					{
						ID:   uuid.New(),
						Name: validArrangementName,
					},
				},
			},
		},
		{
			"Maximal",
			requests.BulkUpdateSongArrangementsRequest{
				SongID: uuid.New(),
				Requests: []requests.UpdateSongArrangementRequest{
					{
						ID:   uuid.New(),
						Name: validArrangementName,
						Occurrences: []requests.UpdateSongSectionOccurrencesRequest{
							{SectionID: uuid.New(), Occurrences: 1},
							{SectionID: uuid.New(), Occurrences: 0},
							{SectionID: uuid.New(), Occurrences: 3},
						},
					},
				},
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

func TestValidateBulkUpdateSongArrangementsRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.BulkUpdateSongArrangementsRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// Song ID Test Cases
		{
			"Song ID invalid because it's required",
			requests.BulkUpdateSongArrangementsRequest{
				SongID: uuid.Nil,
				Requests: []requests.UpdateSongArrangementRequest{
					{
						ID:   uuid.New(),
						Name: validArrangementName,
					},
				},
			},
			"SongID",
			"required",
		},
		// Requests Test Cases
		{
			"Requests are invalid because it requires at least 1 element",
			requests.BulkUpdateSongArrangementsRequest{
				SongID:   uuid.New(),
				Requests: []requests.UpdateSongArrangementRequest{},
			},
			"Requests",
			"min",
		},
		// Requests - ID Test Cases
		{
			"ID is invalid because it's required",
			requests.BulkUpdateSongArrangementsRequest{
				SongID: uuid.New(),
				Requests: []requests.UpdateSongArrangementRequest{
					{
						ID:   uuid.Nil,
						Name: validArrangementName,
					},
				},
			},
			"Requests[0].ID",
			"required",
		},
		// Requests - Name Test Cases
		{
			"Name is invalid because it's required",
			requests.BulkUpdateSongArrangementsRequest{
				SongID: uuid.New(),
				Requests: []requests.UpdateSongArrangementRequest{
					{
						ID:   uuid.New(),
						Name: "",
					},
				},
			},
			"Requests[0].Name",
			"required",
		},
		{
			"Name is invalid because it has too many characters",
			requests.BulkUpdateSongArrangementsRequest{
				SongID: uuid.New(),
				Requests: []requests.UpdateSongArrangementRequest{
					{
						ID:   uuid.New(),
						Name: strings.Repeat("a", 31),
					},
				},
			},
			"Requests[0].Name",
			"max",
		},
		// Requests - Occurrences - ID Test Cases
		{
			"Occurrences are invalid because the first element has an empty SectionID",
			requests.BulkUpdateSongArrangementsRequest{
				SongID: uuid.New(),
				Requests: []requests.UpdateSongArrangementRequest{
					{
						ID:          uuid.New(),
						Name:        validArrangementName,
						Occurrences: []requests.UpdateSongSectionOccurrencesRequest{{SectionID: uuid.Nil, Occurrences: 1}},
					},
				},
			},
			"Requests[0].Occurrences[0].SectionID",
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
			assert.Contains(t, errCode.Error.Error(), "BulkUpdateSongArrangementsRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateUpdateDefaultSongArrangementRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.UpdateDefaultSongArrangementRequest{
		ID:     &[]uuid.UUID{uuid.New()}[0],
		SongID: uuid.New(),
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)
}

func TestValidateUpdateDefaultSongArrangementRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.UpdateDefaultSongArrangementRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// Song ID Test Cases
		{
			"Song ID is invalid because it's required",
			requests.UpdateDefaultSongArrangementRequest{
				ID:     &[]uuid.UUID{uuid.New()}[0],
				SongID: uuid.Nil,
			},
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
			assert.Contains(t, errCode.Error.Error(), "UpdateDefaultSongArrangementRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateMoveSongArrangementRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.MoveSongArrangementRequest{
		ID:     uuid.New(),
		OverID: uuid.New(),
		SongID: uuid.New(),
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)
}

func TestValidateMoveSongArrangementRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.MoveSongArrangementRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// ID Test Cases
		{
			"ID is invalid because it's required",
			requests.MoveSongArrangementRequest{ID: uuid.Nil, OverID: uuid.New(), SongID: uuid.New()},
			"ID",
			"required",
		},
		// Over ID Test Cases
		{
			"Over ID is invalid because it's required",
			requests.MoveSongArrangementRequest{ID: uuid.New(), OverID: uuid.Nil, SongID: uuid.New()},
			"OverID",
			"required",
		},
		// Song ID Test Cases
		{
			"Song ID is invalid because it's required",
			requests.MoveSongArrangementRequest{ID: uuid.New(), OverID: uuid.New(), SongID: uuid.Nil},
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
			assert.Contains(t, errCode.Error.Error(), "MoveSongArrangementRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}
