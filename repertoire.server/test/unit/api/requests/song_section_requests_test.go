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

var validSectionName = "James Solo"

func TestValidateCreateSongSectionRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	tests := []struct {
		name    string
		request requests.CreateSongSectionRequest
	}{
		{
			"Minimal",
			requests.CreateSongSectionRequest{
				SongID: uuid.New(),
				Name:   validSectionName,
				TypeID: uuid.New(),
			},
		},
		{
			"Maximal",
			requests.CreateSongSectionRequest{
				SongID:       uuid.New(),
				Name:         validSectionName,
				TypeID:       uuid.New(),
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

func TestValidateCreateSongSectionRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.CreateSongSectionRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// Song ID Test Cases
		{
			"Song ID is invalid because it's required",
			requests.CreateSongSectionRequest{
				SongID: uuid.Nil,
				Name:   validSectionName,
				TypeID: uuid.New(),
			},
			"SongID",
			"required",
		},
		// Name Test Cases
		{
			"Name is invalid because it's required",
			requests.CreateSongSectionRequest{
				SongID: uuid.New(),
				Name:   "",
				TypeID: uuid.New(),
			},
			"Name",
			"required",
		},
		// Name Test Cases
		{
			"Name is invalid because it has too many characters",
			requests.CreateSongSectionRequest{
				SongID: uuid.New(),
				Name:   strings.Repeat("a", 31),
				TypeID: uuid.New(),
			},
			"Name",
			"max",
		},
		// Type ID Test Cases
		{
			"Type ID is invalid because it's required",
			requests.CreateSongSectionRequest{
				SongID: uuid.New(),
				Name:   validSectionName,
				TypeID: uuid.Nil,
			},
			"TypeID",
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
			assert.Contains(t, errCode.Error.Error(), "CreateSongSectionRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateAddPerfectSongRehearsalRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.AddPerfectSongRehearsalRequest{
		ID: uuid.New(),
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)

}

func TestValidateAddPerfectSongRehearsalRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.AddPerfectSongRehearsalRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// ID Test Cases
		{
			"ID is invalid because it's required",
			requests.AddPerfectSongRehearsalRequest{ID: uuid.Nil},
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
			assert.Contains(t, errCode.Error.Error(), "AddPerfectSongRehearsalRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateAddPerfectSongRehearsalsRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.AddPerfectSongRehearsalsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)

}

func TestValidateAddPerfectSongRehearsalsRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.AddPerfectSongRehearsalsRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// IDs Test Cases
		{
			"IDs is invalid because it requires at least one ID",
			requests.AddPerfectSongRehearsalsRequest{IDs: []uuid.UUID{}},
			"ID",
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
			assert.Contains(t, errCode.Error.Error(), "AddPerfectSongRehearsalsRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateAddPartialSongRehearsalRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.AddPartialSongRehearsalRequest{
		ID: uuid.New(),
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)

}

func TestValidateAddPartialSongRehearsalRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.AddPartialSongRehearsalRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// ID Test Cases
		{
			"ID is invalid because it's required",
			requests.AddPartialSongRehearsalRequest{ID: uuid.Nil},
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
			assert.Contains(t, errCode.Error.Error(), "AddPartialSongRehearsalRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateUpdateSongSectionRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	tests := []struct {
		name    string
		request requests.UpdateSongSectionRequest
	}{
		{
			"Minimal",
			requests.UpdateSongSectionRequest{
				ID:         uuid.New(),
				Name:       validSectionName,
				Confidence: 100,
				Rehearsals: 23,
				TypeID:     uuid.New(),
			},
		},
		{
			"Maximal",
			requests.UpdateSongSectionRequest{
				ID:           uuid.New(),
				Name:         validSectionName,
				Confidence:   100,
				Rehearsals:   23,
				TypeID:       uuid.New(),
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

func TestValidateUpdateSongSectionRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.UpdateSongSectionRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// ID Test Cases
		{
			"Song ID is invalid because it's required",
			requests.UpdateSongSectionRequest{
				ID:     uuid.Nil,
				Name:   validSectionName,
				TypeID: uuid.New(),
			},
			"ID",
			"required",
		},
		// Name Test Cases
		{
			"Name is invalid because it's required",
			requests.UpdateSongSectionRequest{
				ID:     uuid.New(),
				Name:   "",
				TypeID: uuid.New(),
			},
			"Name",
			"required",
		},
		{
			"Name is invalid because it has too many characters",
			requests.UpdateSongSectionRequest{
				ID:     uuid.New(),
				Name:   strings.Repeat("a", 31),
				TypeID: uuid.New(),
			},
			"Name",
			"max",
		},
		// Confidence Test Cases
		{
			"Confidence is invalid because it is greater than 100",
			requests.UpdateSongSectionRequest{
				ID:         uuid.New(),
				Name:       validSectionName,
				Confidence: 101,
				TypeID:     uuid.New(),
			},
			"Confidence",
			"max",
		},
		// Type ID Test Cases
		{
			"Type ID is invalid because it's required",
			requests.UpdateSongSectionRequest{
				ID:     uuid.New(),
				Name:   validSectionName,
				TypeID: uuid.Nil,
			},
			"TypeID",
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
			assert.Contains(t, errCode.Error.Error(), "UpdateSongSectionRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateUpdateSongSectionsOccurrencesRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.UpdateSongSectionsOccurrencesRequest{
		SongID: uuid.New(),
		Sections: []requests.UpdateSectionOccurrencesRequest{
			{ID: uuid.New()},
			{ID: uuid.New(), Occurrences: 1},
		},
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)

}

func TestValidateUpdateSongSectionsOccurrencesRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.UpdateSongSectionsOccurrencesRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// Song ID Test Cases
		{
			"Song ID is invalid because it's required",
			requests.UpdateSongSectionsOccurrencesRequest{
				SongID:   uuid.Nil,
				Sections: []requests.UpdateSectionOccurrencesRequest{{ID: uuid.New()}},
			},
			"SongID",
			"required",
		},
		// Sections Test Cases
		{
			"Sections is invalid because it requires at least one element",
			requests.UpdateSongSectionsOccurrencesRequest{
				SongID:   uuid.New(),
				Sections: []requests.UpdateSectionOccurrencesRequest{},
			},
			"Sections",
			"min",
		},
		{
			"Sections is invalid because it requires that all elements have an id",
			requests.UpdateSongSectionsOccurrencesRequest{
				SongID:   uuid.New(),
				Sections: []requests.UpdateSectionOccurrencesRequest{{ID: uuid.Nil}},
			},
			"Sections[0].ID",
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
			assert.Contains(t, errCode.Error.Error(), "UpdateSongSectionsOccurrencesRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateUpdateSongSectionsPartialOccurrencesRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.UpdateSongSectionsPartialOccurrencesRequest{
		SongID: uuid.New(),
		Sections: []requests.UpdateSectionPartialOccurrencesRequest{
			{ID: uuid.New()},
			{ID: uuid.New(), PartialOccurrences: 1},
		},
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)

}

func TestValidateUpdateSongSectionsPartialOccurrencesRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.UpdateSongSectionsPartialOccurrencesRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// Song ID Test Cases
		{
			"Song ID is invalid because it's required",
			requests.UpdateSongSectionsPartialOccurrencesRequest{
				SongID:   uuid.Nil,
				Sections: []requests.UpdateSectionPartialOccurrencesRequest{{ID: uuid.New()}},
			},
			"SongID",
			"required",
		},
		// Sections Test Cases
		{
			"Sections is invalid because it requires at least one element",
			requests.UpdateSongSectionsPartialOccurrencesRequest{
				SongID:   uuid.New(),
				Sections: []requests.UpdateSectionPartialOccurrencesRequest{},
			},
			"Sections",
			"min",
		},
		{
			"Sections is invalid because it requires that all elements have an id",
			requests.UpdateSongSectionsPartialOccurrencesRequest{
				SongID:   uuid.New(),
				Sections: []requests.UpdateSectionPartialOccurrencesRequest{{ID: uuid.Nil}},
			},
			"Sections[0].ID",
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
			assert.Contains(t, errCode.Error.Error(), "UpdateSongSectionsPartialOccurrencesRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateUpdateAllSongSectionsRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	tests := []struct {
		name    string
		request requests.UpdateAllSongSectionsRequest
	}{
		{
			"Non Optional",
			requests.UpdateAllSongSectionsRequest{
				SongID: uuid.New(),
			},
		},
		{
			"All Filled",
			requests.UpdateAllSongSectionsRequest{
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

func TestValidateUpdateAllSongSectionsRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.UpdateAllSongSectionsRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// Song ID Test Cases
		{
			"Song ID is invalid because it's required",
			requests.UpdateAllSongSectionsRequest{SongID: uuid.Nil},
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
			assert.Contains(t, errCode.Error.Error(), "UpdateAllSongSectionsRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateMoveSongSectionRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.MoveSongSectionRequest{
		ID:     uuid.New(),
		OverID: uuid.New(),
		SongID: uuid.New(),
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)
}

func TestValidateMoveSongSectionRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.MoveSongSectionRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// ID Test Cases
		{
			"ID is invalid because it's required",
			requests.MoveSongSectionRequest{ID: uuid.Nil, OverID: uuid.New(), SongID: uuid.New()},
			"ID",
			"required",
		},
		// Over ID Test Cases
		{
			"Over ID is invalid because it's required",
			requests.MoveSongSectionRequest{ID: uuid.New(), OverID: uuid.Nil, SongID: uuid.New()},
			"OverID",
			"required",
		},
		// Song ID Test Cases
		{
			"Song ID is invalid because it's required",
			requests.MoveSongSectionRequest{ID: uuid.New(), OverID: uuid.New(), SongID: uuid.Nil},
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
			assert.Contains(t, errCode.Error.Error(), "MoveSongSectionRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateBulkRehearsalsSongSectionsRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.BulkRehearsalsSongSectionsRequest{
		Sections: []requests.BulkRehearsalsSongSectionRequest{{ID: uuid.New(), Rehearsal: 0}},
		SongID:   uuid.New(),
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)
}

func TestValidateBulkRehearsalsSongSectionsRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.BulkRehearsalsSongSectionsRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// Sections Test Cases
		{
			"Sections is invalid because it requires at least 1 section",
			requests.BulkRehearsalsSongSectionsRequest{
				Sections: []requests.BulkRehearsalsSongSectionRequest{},
				SongID:   uuid.New(),
			},
			"Sections",
			"min",
		},
		// Sections Test Cases
		{
			"Sections is invalid because it requires at least 1 section",
			requests.BulkRehearsalsSongSectionsRequest{
				Sections: []requests.BulkRehearsalsSongSectionRequest{{ID: uuid.Nil, Rehearsal: 0}},
				SongID:   uuid.New(),
			},
			"Sections[0].ID",
			"required",
		},
		// Song ID Test Cases
		{
			"Song ID is invalid because it's required",
			requests.BulkRehearsalsSongSectionsRequest{
				Sections: []requests.BulkRehearsalsSongSectionRequest{{ID: uuid.New(), Rehearsal: 0}},
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
			assert.Contains(t, errCode.Error.Error(), "BulkRehearsalsSongSectionsRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateBulkDeleteSongSectionsRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.BulkDeleteSongSectionsRequest{
		IDs:    []uuid.UUID{uuid.New()},
		SongID: uuid.New(),
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)
}

func TestValidateBulkDeleteSongSectionsRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.BulkDeleteSongSectionsRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// IDs Test Cases
		{
			"IDs is invalid because it requirest at least 1 ID",
			requests.BulkDeleteSongSectionsRequest{IDs: []uuid.UUID{}, SongID: uuid.New()},
			"IDs",
			"min",
		},
		// Song ID Test Cases
		{
			"Song ID is invalid because it's required",
			requests.BulkDeleteSongSectionsRequest{IDs: []uuid.UUID{uuid.New()}, SongID: uuid.Nil},
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
			assert.Contains(t, errCode.Error.Error(), "BulkDeleteSongSectionsRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}
