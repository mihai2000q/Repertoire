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

var validPartName = "James Solo"

func TestValidateCreateSongPartRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	tests := []struct {
		name    string
		request requests.CreateSongPartRequest
	}{
		{
			"Minimal",
			requests.CreateSongPartRequest{
				SongID: uuid.New(),
				Name:   validPartName,
			},
		},
		{
			"Maximal",
			requests.CreateSongPartRequest{
				SongID:       uuid.New(),
				SectionID:    &[]uuid.UUID{uuid.New()}[0],
				Name:         validPartName,
				BandMemberID: &[]uuid.UUID{uuid.New()}[0],
				InstrumentID: &[]uuid.UUID{uuid.New()}[0],
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

func TestValidateCreateSongPartRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.CreateSongPartRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// Song ID Test Cases
		{
			"Song ID is invalid because it's required",
			requests.CreateSongPartRequest{
				SongID: uuid.Nil,
				Name:   validPartName,
			},
			"SongID",
			"required",
		},
		// Name Test Cases
		{
			"Name is invalid because it's required",
			requests.CreateSongPartRequest{
				SongID: uuid.New(),
				Name:   "",
			},
			"Name",
			"required",
		},
		{
			"Name is invalid because it has too many characters",
			requests.CreateSongPartRequest{
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
			assert.Contains(t, errCode.Error.Error(), "CreateSongPartRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateUpdateSongPartRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	tests := []struct {
		name    string
		request requests.UpdateSongPartRequest
	}{
		{
			"Minimal",
			requests.UpdateSongPartRequest{
				ID:         uuid.New(),
				Name:       validPartName,
				Confidence: 100,
				Rehearsals: 23,
			},
		},
		{
			"Maximal",
			requests.UpdateSongPartRequest{
				ID:           uuid.New(),
				Name:         validPartName,
				Confidence:   100,
				Rehearsals:   23,
				BandMemberID: &[]uuid.UUID{uuid.New()}[0],
				InstrumentID: &[]uuid.UUID{uuid.New()}[0],
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

func TestValidateUpdateSongPartRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.UpdateSongPartRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// ID Test Cases
		{
			"Song ID is invalid because it's required",
			requests.UpdateSongPartRequest{
				ID:   uuid.Nil,
				Name: validPartName,
			},
			"ID",
			"required",
		},
		// Name Test Cases
		{
			"Name is invalid because it's required",
			requests.UpdateSongPartRequest{
				ID:   uuid.New(),
				Name: "",
			},
			"Name",
			"required",
		},
		{
			"Name is invalid because it has too many characters",
			requests.UpdateSongPartRequest{
				ID:   uuid.New(),
				Name: strings.Repeat("a", 31),
			},
			"Name",
			"max",
		},
		// Confidence Test Cases
		{
			"Confidence is invalid because it is greater than 100",
			requests.UpdateSongPartRequest{
				ID:         uuid.New(),
				Name:       validPartName,
				Confidence: 101,
			},
			"Confidence",
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
			assert.Contains(t, errCode.Error.Error(), "UpdateSongPartRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateUpdateAllSongPartsRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	tests := []struct {
		name    string
		request requests.UpdateAllSongPartsRequest
	}{
		{
			"Non Optional",
			requests.UpdateAllSongPartsRequest{
				SongID: uuid.New(),
			},
		},
		{
			"All Filled",
			requests.UpdateAllSongPartsRequest{
				SongID:       uuid.New(),
				InstrumentID: &[]uuid.UUID{uuid.New()}[0],
				BandMemberID: &[]uuid.UUID{uuid.New()}[0],
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

func TestValidateUpdateAllSongPartsRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.UpdateAllSongPartsRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// Song ID Test Cases
		{
			"Song ID is invalid because it's required",
			requests.UpdateAllSongPartsRequest{SongID: uuid.Nil},
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
			assert.Contains(t, errCode.Error.Error(), "UpdateAllSongPartsRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateMoveSongPartInSongRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.MoveSongPartInSongRequest{
		ID:     uuid.New(),
		OverID: uuid.New(),
		SongID: uuid.New(),
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)
}

func TestValidateMoveSongPartInSongRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.MoveSongPartInSongRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// ID Test Cases
		{
			"ID is invalid because it's required",
			requests.MoveSongPartInSongRequest{ID: uuid.Nil, OverID: uuid.New(), SongID: uuid.New()},
			"ID",
			"required",
		},
		// Over ID Test Cases
		{
			"Over ID is invalid because it's required",
			requests.MoveSongPartInSongRequest{ID: uuid.New(), OverID: uuid.Nil, SongID: uuid.New()},
			"OverID",
			"required",
		},
		// Song ID Test Cases
		{
			"Song ID is invalid because it's required",
			requests.MoveSongPartInSongRequest{ID: uuid.New(), OverID: uuid.New(), SongID: uuid.Nil},
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
			assert.Contains(t, errCode.Error.Error(), "MoveSongPartInSongRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateMoveSongPartInSectionRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	tests := []struct {
		name    string
		request requests.MoveSongPartInSectionRequest
	}{
		{
			"Minimal",
			requests.MoveSongPartInSectionRequest{
				ID:        uuid.New(),
				OverID:    uuid.New(),
				SectionID: uuid.New(),
			},
		},
		{
			"Maximal",
			requests.MoveSongPartInSectionRequest{
				ID:            uuid.New(),
				OverID:        uuid.New(),
				SectionID:     uuid.New(),
				OverSectionID: &[]uuid.UUID{uuid.New()}[0],
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

func TestValidateMoveSongPartInSectionRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.MoveSongPartInSectionRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// ID Test Cases
		{
			"ID is invalid because it's required",
			requests.MoveSongPartInSectionRequest{ID: uuid.Nil, OverID: uuid.New(), SectionID: uuid.New()},
			"ID",
			"required",
		},
		// Over ID Test Cases
		{
			"Over ID is invalid because it's required",
			requests.MoveSongPartInSectionRequest{ID: uuid.New(), OverID: uuid.Nil, SectionID: uuid.New()},
			"OverID",
			"required",
		},
		// Section ID Test Cases
		{
			"Section ID is invalid because it's required",
			requests.MoveSongPartInSectionRequest{ID: uuid.New(), OverID: uuid.New(), SectionID: uuid.Nil},
			"SectionID",
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
			assert.Contains(t, errCode.Error.Error(), "MoveSongPartInSectionRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateBulkRehearsalsSongPartsRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.BulkRehearsalsSongPartsRequest{
		Requests: []requests.BulkRehearsalsSongPartRequest{{ID: uuid.New(), Rehearsals: 0}},
		SongID:   uuid.New(),
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)
}

func TestValidateBulkRehearsalsSongPartsRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.BulkRehearsalsSongPartsRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// Requests Test Cases
		{
			"Requests is invalid because it requires at least 1 part",
			requests.BulkRehearsalsSongPartsRequest{
				Requests: []requests.BulkRehearsalsSongPartRequest{},
				SongID:   uuid.New(),
			},
			"Requests",
			"min",
		},
		{
			"Requests is invalid because it requires at least 1 part",
			requests.BulkRehearsalsSongPartsRequest{
				Requests: []requests.BulkRehearsalsSongPartRequest{{ID: uuid.Nil, Rehearsals: 0}},
				SongID:   uuid.New(),
			},
			"Requests[0].ID",
			"required",
		},
		// Song ID Test Cases
		{
			"Song ID is invalid because it's required",
			requests.BulkRehearsalsSongPartsRequest{
				Requests: []requests.BulkRehearsalsSongPartRequest{{ID: uuid.New(), Rehearsals: 0}},
				SongID:   uuid.Nil,
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
			assert.Contains(t, errCode.Error.Error(), "BulkRehearsalsSongPartsRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateBulkDeleteSongPartsRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.BulkDeleteSongPartsRequest{
		IDs:    []uuid.UUID{uuid.New()},
		SongID: uuid.New(),
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)
}

func TestValidateBulkDeleteSongPartsRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.BulkDeleteSongPartsRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// IDs Test Cases
		{
			"IDs is invalid because it requirest at least 1 ID",
			requests.BulkDeleteSongPartsRequest{IDs: []uuid.UUID{}, SongID: uuid.New()},
			"IDs",
			"min",
		},
		// Song ID Test Cases
		{
			"Song ID is invalid because it's required",
			requests.BulkDeleteSongPartsRequest{IDs: []uuid.UUID{uuid.New()}, SongID: uuid.Nil},
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
			assert.Contains(t, errCode.Error.Error(), "BulkDeleteSongPartsRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}
