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

func TestValidateUpdateSongArrangementRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.UpdateSongArrangementRequest{
		ID:   uuid.New(),
		Name: validArrangementName,
		Occurrences: []requests.UpdateSectionOccurrencesRequest{
			{SectionID: uuid.New(), Occurrences: 1},
			{SectionID: uuid.New(), Occurrences: 0},
			{SectionID: uuid.New(), Occurrences: 3},
		},
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)
}

func TestValidateUpdateSongArrangementRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.UpdateSongArrangementRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// ID Test Cases
		{
			"ID is invalid because it's required",
			requests.UpdateSongArrangementRequest{
				ID:          uuid.Nil,
				Name:        validArrangementName,
				Occurrences: []requests.UpdateSectionOccurrencesRequest{{SectionID: uuid.New(), Occurrences: 1}},
			},
			"ID",
			"required",
		},
		// Name Test Cases
		{
			"Name is invalid because it's required",
			requests.UpdateSongArrangementRequest{
				ID:          uuid.New(),
				Name:        "",
				Occurrences: []requests.UpdateSectionOccurrencesRequest{{SectionID: uuid.New(), Occurrences: 1}},
			},
			"Name",
			"required",
		},
		{
			"Name is invalid because it has too many characters",
			requests.UpdateSongArrangementRequest{
				ID:          uuid.New(),
				Name:        strings.Repeat("a", 31),
				Occurrences: []requests.UpdateSectionOccurrencesRequest{{SectionID: uuid.New(), Occurrences: 1}},
			},
			"Name",
			"max",
		},
		// Occurrences - ID Test Cases
		{
			"Occurrences are invalid because the first element has an empty SectionID",
			requests.UpdateSongArrangementRequest{
				ID:          uuid.New(),
				Name:        validArrangementName,
				Occurrences: []requests.UpdateSectionOccurrencesRequest{{SectionID: uuid.Nil, Occurrences: 1}},
			},
			"Occurrences[0].SectionID",
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
			assert.Contains(t, errCode.Error.Error(), "UpdateSongArrangementRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateUpdateDefaultSongArrangementRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.UpdateDefaultSongArrangementRequest{
		ID:     uuid.New(),
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
		// ID Test Cases
		{
			"ID is invalid because it's required",
			requests.UpdateDefaultSongArrangementRequest{
				ID:     uuid.Nil,
				SongID: uuid.New(),
			},
			"ID",
			"required",
		},
		// Song ID Test Cases
		{
			"Song ID is invalid because it's required",
			requests.UpdateDefaultSongArrangementRequest{
				ID:     uuid.New(),
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
