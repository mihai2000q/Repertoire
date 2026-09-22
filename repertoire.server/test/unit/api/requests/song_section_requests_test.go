package requests

import (
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/api/validation"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateGetSongSectionsRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.GetSongSectionsRequest{
		SongID: uuid.New(),
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)
}

func TestValidateGetSongSectionsRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.GetSongSectionsRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// Song ID Test Cases
		{
			"Song ID is invalid because it's required",
			requests.GetSongSectionsRequest{
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
			require.NotNil(t, errCode)
			assert.Len(t, errCode.Error, 1)
			assert.Contains(t, errCode.Error.Error(), "GetSongSectionsRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

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
				SongID: uuid.New(),
				Name:   validSectionName,
				TypeID: uuid.New(),
				Parts: []requests.CreateSongSectionPartRequest{
					{PartID: &[]uuid.UUID{uuid.New()}[0]},
					{PartID: &[]uuid.UUID{uuid.New()}[0], BandMemberIDs: []uuid.UUID{uuid.New()}},
					{NewPart: &requests.CreateNewSongPartRequest{Name: "Chorus-1"}},
					{
						NewPart: &requests.CreateNewSongPartRequest{
							Name:         "Chorus-1",
							InstrumentID: &[]uuid.UUID{uuid.New()}[0],
						},
					},
					{
						NewPart:       &requests.CreateNewSongPartRequest{Name: "Chorus-1"},
						BandMemberIDs: []uuid.UUID{uuid.New()},
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

func TestValidateCreateSongSectionRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	partID := &[]uuid.UUID{uuid.New()}[0]

	bandMemberID := uuid.New()

	tests := []struct {
		name                  string
		request               requests.CreateSongSectionRequest
		expectedInvalidFields []string
		expectedFailedTags    []string
	}{
		// Song ID Test Cases
		{
			"Song ID is invalid because it's required",
			requests.CreateSongSectionRequest{
				SongID: uuid.Nil,
				Name:   validSectionName,
				TypeID: uuid.New(),
			},
			[]string{"SongID"},
			[]string{"required"},
		},
		// Name Test Cases
		{
			"Name is invalid because it's required",
			requests.CreateSongSectionRequest{
				SongID: uuid.New(),
				Name:   "",
				TypeID: uuid.New(),
			},
			[]string{"Name"},
			[]string{"required"},
		},
		{
			"Name is invalid because it has too many characters",
			requests.CreateSongSectionRequest{
				SongID: uuid.New(),
				Name:   strings.Repeat("a", 31),
				TypeID: uuid.New(),
			},
			[]string{"Name"},
			[]string{"max"},
		},
		// Type ID Test Cases
		{
			"Type ID is invalid because it's required",
			requests.CreateSongSectionRequest{
				SongID: uuid.New(),
				Name:   validSectionName,
				TypeID: uuid.Nil,
			},
			[]string{"TypeID"},
			[]string{"required"},
		},
		// Parts Test Cases
		{
			"Parts is invalid because the ids on Part ID must be unique",
			requests.CreateSongSectionRequest{
				SongID: uuid.New(),
				Name:   validSectionName,
				TypeID: uuid.New(),
				Parts: []requests.CreateSongSectionPartRequest{
					{PartID: partID},
					{PartID: partID},
				},
			},
			[]string{"Parts"},
			[]string{"unique_ids"},
		},
		// Parts Test Cases - Part ID and New Part
		{
			"Parts is invalid because Part ID and New Part cannot be both set at the same time",
			requests.CreateSongSectionRequest{
				SongID: uuid.New(),
				Name:   validSectionName,
				TypeID: uuid.New(),
				Parts: []requests.CreateSongSectionPartRequest{
					{
						PartID:  &[]uuid.UUID{uuid.New()}[0],
						NewPart: &requests.CreateNewSongPartRequest{Name: "Chorus-1"},
					},
				},
			},
			[]string{"Parts[0].PartID", "Parts[0].NewPart"},
			[]string{"excluded_with", "excluded_with"},
		},
		// Parts Test Cases - Part ID and New Part
		{
			"Parts is invalid because Part ID and New Part cannot be both set at the same time",
			requests.CreateSongSectionRequest{
				SongID: uuid.New(),
				Name:   validSectionName,
				TypeID: uuid.New(),
				Parts: []requests.CreateSongSectionPartRequest{
					{
						PartID:  &[]uuid.UUID{uuid.New()}[0],
						NewPart: &requests.CreateNewSongPartRequest{Name: "Chorus-1"},
					},
				},
			},
			[]string{"Parts[0].PartID", "Parts[0].NewPart"},
			[]string{"excluded_with", "excluded_with"},
		},
		// Parts Test Cases - Band Member IDs
		{
			"Parts is invalid because Band Member ID requires either Part ID or New Part to be set",
			requests.CreateSongSectionRequest{
				SongID: uuid.New(),
				Name:   validSectionName,
				TypeID: uuid.New(),
				Parts: []requests.CreateSongSectionPartRequest{
					{BandMemberIDs: []uuid.UUID{uuid.New()}},
				},
			},
			[]string{"Parts[0].BandMemberIDs"},
			[]string{"excluded_without_all"},
		},
		// Parts Test Cases - Band Member IDs
		{
			"Parts is invalid because Band Member IDs requires the ids to be unique",
			requests.CreateSongSectionRequest{
				SongID: uuid.New(),
				Name:   validSectionName,
				TypeID: uuid.New(),
				Parts: []requests.CreateSongSectionPartRequest{
					{PartID: &[]uuid.UUID{uuid.New()}[0], BandMemberIDs: []uuid.UUID{bandMemberID, bandMemberID}},
				},
			},
			[]string{"Parts[0].BandMemberIDs"},
			[]string{"unique"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			_uut := validation.NewValidator(nil)

			// when
			errCode := _uut.Validate(tt.request)

			// then
			require.NotNil(t, errCode)
			assert.Len(t, tt.expectedFailedTags, len(tt.expectedInvalidFields))
			assert.Len(t, errCode.Error, len(tt.expectedFailedTags))
			for _, expectedInvalidField := range tt.expectedInvalidFields {
				assert.Contains(t, errCode.Error.Error(), "CreateSongSectionRequest."+expectedInvalidField)
			}
			for _, expectedFailedTag := range tt.expectedFailedTags {
				assert.Contains(t, errCode.Error.Error(), "'"+expectedFailedTag+"' tag")
			}
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
				ID:     uuid.New(),
				Name:   validSectionName,
				TypeID: uuid.New(),
			},
		},
		{
			"Maximal",
			requests.UpdateSongSectionRequest{
				ID:      uuid.New(),
				Name:    validSectionName,
				TypeID:  uuid.New(),
				PartIDs: []uuid.UUID{uuid.New()},
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
	partID := uuid.New()

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
		// Part IDs Test Cases
		{
			"Part IDs is invalid because it requires the ids to be unique",
			requests.UpdateSongSectionRequest{
				ID:      uuid.New(),
				Name:    validSectionName,
				TypeID:  uuid.New(),
				PartIDs: []uuid.UUID{partID, partID},
			},
			"PartIDs",
			"unique",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			_uut := validation.NewValidator(nil)

			// when
			errCode := _uut.Validate(tt.request)

			// then
			require.NotNil(t, errCode)
			assert.Len(t, errCode.Error, 1)
			assert.Contains(t, errCode.Error.Error(), "UpdateSongSectionRequest."+tt.expectedInvalidField)
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
			require.NotNil(t, errCode)
			assert.Len(t, errCode.Error, 1)
			assert.Contains(t, errCode.Error.Error(), "MoveSongSectionRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateBulkDeleteSongSectionsRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	tests := []struct {
		name    string
		request requests.BulkDeleteSongSectionsRequest
	}{
		{
			"Minimal",
			requests.BulkDeleteSongSectionsRequest{
				IDs:    []uuid.UUID{uuid.New()},
				SongID: uuid.New(),
			},
		},
		{
			"Maximal",
			requests.BulkDeleteSongSectionsRequest{
				IDs:     []uuid.UUID{uuid.New()},
				SongID:  uuid.New(),
				PartIDs: []uuid.UUID{uuid.New()},
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

func TestValidateBulkDeleteSongSectionsRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	id := uuid.New()

	tests := []struct {
		name                 string
		request              requests.BulkDeleteSongSectionsRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// IDs Test Cases
		{
			"IDs is invalid because it requires at least 1 ID",
			requests.BulkDeleteSongSectionsRequest{IDs: []uuid.UUID{}, SongID: uuid.New()},
			"IDs",
			"min",
		},
		{
			"IDs is invalid because it requires the ids to be unique",
			requests.BulkDeleteSongSectionsRequest{IDs: []uuid.UUID{id, id}, SongID: uuid.New()},
			"IDs",
			"unique",
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
			require.NotNil(t, errCode)
			assert.Len(t, errCode.Error, 1)
			assert.Contains(t, errCode.Error.Error(), "BulkDeleteSongSectionsRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}
