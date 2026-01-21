package requests

import (
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/api/validation"
	"repertoire/server/internal"
	"repertoire/server/internal/enums"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidateGetSongsRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	tests := []struct {
		name    string
		request requests.GetSongsRequest
	}{
		{
			"All Null",
			requests.GetSongsRequest{},
		},
		{
			"Nothing Null",
			requests.GetSongsRequest{
				CurrentPage: &[]int{1}[0],
				PageSize:    &[]int{1}[0],
				OrderBy:     []string{"title asc nulls first", "created_at desc"},
				SearchBy:    []string{"title ~* something entirely different", "is_recorded <> false"},
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
		request              requests.GetSongsRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// Current Page Test Cases
		{
			"Current Page is invalid because it should be greater than 0",
			requests.GetSongsRequest{CurrentPage: &[]int{0}[0], PageSize: &[]int{1}[0]},
			"CurrentPage",
			"gt",
		},
		{
			"Current Page is invalid because page size is null",
			requests.GetSongsRequest{PageSize: &[]int{1}[0]},
			"CurrentPage",
			"required_with",
		},
		// Page Size Test Cases
		{
			"Page Size is invalid because it should be greater than 0",
			requests.GetSongsRequest{PageSize: &[]int{0}[0], CurrentPage: &[]int{1}[0]},
			"PageSize",
			"gt",
		},
		{
			"Page Size is invalid because current page is null",
			requests.GetSongsRequest{CurrentPage: &[]int{1}[0]},
			"PageSize",
			"required_with",
		},
		// Order By Test Cases
		{
			"Order By is invalid because of the invalid 'first'",
			requests.GetSongsRequest{OrderBy: []string{"title asc", "songs asc nulls firsts"}},
			"OrderBy",
			"order_by",
		},
		// Search By Test Cases
		{
			"Search By is invalid because the operator is not supported",
			requests.GetSongsRequest{SearchBy: []string{"title != okay", "songs is not nullish"}},
			"SearchBy",
			"search_by",
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
			assert.Contains(t, errCode.Error.Error(), "GetSongsRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateGetSongFiltersMetadataRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	tests := []struct {
		name    string
		request requests.GetSongFiltersMetadataRequest
	}{
		{
			"All Null",
			requests.GetSongFiltersMetadataRequest{},
		},
		{
			"Nothing Null",
			requests.GetSongFiltersMetadataRequest{
				SearchBy: []string{"title = something", "is_recorded <> false"},
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

func TestValidateGetSongFiltersMetadataRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.GetSongFiltersMetadataRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// Search By Test Cases
		{
			"Search By is invalid because the value is missing",
			requests.GetSongFiltersMetadataRequest{SearchBy: []string{"title != okay", "songs ="}},
			"SearchBy",
			"search_by",
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
			assert.Contains(t, errCode.Error.Error(), "GetSongFiltersMetadataRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

var validSongTitle = "Justice For All"

func TestValidateCreateSongRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	tests := []struct {
		name    string
		request requests.CreateSongRequest
	}{
		{
			"All Null",
			requests.CreateSongRequest{Title: validSongTitle},
		},
		{
			"Nothing Null 1",
			requests.CreateSongRequest{
				Title:          validSongTitle,
				Description:    "Something",
				Bpm:            &[]uint{12}[0],
				SongsterrLink:  &[]string{"http://songsterr.com/some-song"}[0],
				YoutubeLink:    &[]string{"https://www.youtube.com/watch?v=9DyxtUCW84o&t=1m3s"}[0],
				GuitarTuningID: &[]uuid.UUID{uuid.New()}[0],
				AlbumID:        &[]uuid.UUID{uuid.New()}[0],
				Sections: []requests.CreateSectionRequest{
					{Name: "A section", TypeID: uuid.New()},
					{Name: "A Second Section", TypeID: uuid.New()},
				},
			},
		},
		{
			"Nothing Null 2",
			requests.CreateSongRequest{
				Title:          validSongTitle,
				Description:    "Something",
				Bpm:            &[]uint{12}[0],
				SongsterrLink:  &[]string{"https://songsterr.com/some-other"}[0],
				YoutubeLink:    &[]string{"https://youtu.be/9DyxtUCW84o?si=2pNX8eaV4KwKfOaF"}[0],
				ReleaseDate:    &[]internal.Date{internal.Date(time.Now())}[0],
				Difficulty:     &[]enums.Difficulty{enums.Easy}[0],
				GuitarTuningID: &[]uuid.UUID{uuid.New()}[0],
				AlbumTitle:     &[]string{"New Album Title"}[0],
				ArtistName:     &[]string{"New Artist Name"}[0],
				Sections: []requests.CreateSectionRequest{
					{Name: "A section", TypeID: uuid.New()},
					{Name: "A Second Section", TypeID: uuid.New()},
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

func TestValidateCreateSongRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                  string
		request               requests.CreateSongRequest
		expectedInvalidFields []string
		expectedFailedTags    []string
	}{
		// Title Test Cases
		{
			"Title is invalid because it's required",
			requests.CreateSongRequest{Title: ""},
			[]string{"Title"},
			[]string{"required"},
		},
		{
			"Title is invalid because it has more than 100 characters",
			requests.CreateSongRequest{Title: strings.Repeat("a", 101)},
			[]string{"Title"},
			[]string{"max"},
		},
		// SongsterrLink Test Cases
		{
			"Songsterr Link is invalid because it is not an url",
			requests.CreateSongRequest{
				Title:         validSongTitle,
				SongsterrLink: &[]string{"scom"}[0],
			},
			[]string{"SongsterrLink"},
			[]string{"url"},
		},
		{
			"Songsterr Link is invalid because it is not a songsterr link",
			requests.CreateSongRequest{
				Title:         validSongTitle,
				SongsterrLink: &[]string{"http://google.com"}[0],
			},
			[]string{"SongsterrLink"},
			[]string{"contains"},
		},
		// YoutubeLink Test Cases
		{
			"Youtube Link is invalid because it is not youtube link",
			requests.CreateSongRequest{
				Title:       validSongTitle,
				YoutubeLink: &[]string{"https://google.com"}[0],
			},
			[]string{"YoutubeLink"},
			[]string{"youtube_link"},
		},
		// Difficulty Test Cases
		{
			"Difficulty is invalid because it is not a Difficulty Enum",
			requests.CreateSongRequest{
				Title:      validSongTitle,
				Difficulty: &[]enums.Difficulty{"Something else"}[0],
			},
			[]string{"Difficulty"},
			[]string{"difficulty_enum"},
		},
		// Album ID
		{
			"Album ID is invalid because the Artist Id is also set",
			requests.CreateSongRequest{
				Title:    validSongTitle,
				AlbumID:  &[]uuid.UUID{uuid.New()}[0],
				ArtistID: &[]uuid.UUID{uuid.New()}[0],
			},
			[]string{"AlbumID"},
			[]string{"excluded_with"},
		},
		{
			"Album ID is invalid because the Artist Name is also set",
			requests.CreateSongRequest{
				Title:      validSongTitle,
				AlbumID:    &[]uuid.UUID{uuid.New()}[0],
				ArtistName: &[]string{"New Artist Name"}[0],
			},
			[]string{"AlbumID"},
			[]string{"excluded_with"},
		},
		// Album Title Test Case
		{
			"Album Title is invalid because it has too many characters",
			requests.CreateSongRequest{
				Title:      validSongTitle,
				AlbumTitle: &[]string{strings.Repeat("a", 101)}[0],
			},
			[]string{"AlbumTitle"},
			[]string{"max"},
		},
		// Album ID and Album Title Test Case
		{
			"Album Title and ID are invalid because only one can be set at a time",
			requests.CreateSongRequest{
				Title:      validSongTitle,
				AlbumID:    &[]uuid.UUID{uuid.New()}[0],
				AlbumTitle: &[]string{"New Album Title"}[0],
			},
			[]string{"AlbumID", "AlbumTitle"},
			[]string{"excluded_with", "excluded_with"},
		},
		// Artist Name Test Case
		{
			"Artist Name is invalid because it has too many characters",
			requests.CreateSongRequest{
				Title:      validSongTitle,
				ArtistName: &[]string{strings.Repeat("a", 101)}[0],
			},
			[]string{"ArtistName"},
			[]string{"max"},
		},
		// Artist ID and Artist Name Test Case
		{
			"Artist Name and ID are invalid because only one can be set at a time",
			requests.CreateSongRequest{
				Title:      validSongTitle,
				ArtistID:   &[]uuid.UUID{uuid.New()}[0],
				ArtistName: &[]string{"New Artist Name"}[0],
			},
			[]string{"ArtistID", "ArtistName"},
			[]string{"excluded_with", "excluded_with"},
		},
		// Sections - Name Test Cases
		{
			"Sections are invalid because the first element has an empty Name",
			requests.CreateSongRequest{
				Title: validSongTitle,
				Sections: []requests.CreateSectionRequest{
					{Name: "", TypeID: uuid.New()},
				},
			},
			[]string{"Sections[0].Name"},
			[]string{"required"},
		},
		// Sections - Name Test Cases
		{
			"Sections are invalid because the first element has a Name with too many characters",
			requests.CreateSongRequest{
				Title: validSongTitle,
				Sections: []requests.CreateSectionRequest{
					{Name: strings.Repeat("a", 31), TypeID: uuid.New()},
				},
			},
			[]string{"Sections[0].Name"},
			[]string{"max"},
		},
		// Sections - Type ID Test Cases
		{
			"Sections are invalid because the first element has an empty Type ID",
			requests.CreateSongRequest{
				Title: validSongTitle,
				Sections: []requests.CreateSectionRequest{
					{Name: "some Name", TypeID: uuid.Nil},
				},
			},
			[]string{"Sections[0].TypeID"},
			[]string{"required"},
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
			assert.Len(t, tt.expectedFailedTags, len(tt.expectedInvalidFields))
			assert.Len(t, errCode.Error, len(tt.expectedFailedTags))
			for _, expectedInvalidField := range tt.expectedInvalidFields {
				assert.Contains(t, errCode.Error.Error(), "CreateSongRequest."+expectedInvalidField)
			}
			for _, expectedFailedTag := range tt.expectedFailedTags {
				assert.Contains(t, errCode.Error.Error(), "'"+expectedFailedTag+"' tag")
			}
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateUpdateSongRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	tests := []struct {
		name    string
		request requests.UpdateSongRequest
	}{
		{
			"Non Optional",
			requests.UpdateSongRequest{
				ID:    uuid.New(),
				Title: validSongTitle,
			},
		},
		{
			"All Filled",
			requests.UpdateSongRequest{
				ID:             uuid.New(),
				Title:          validSongTitle,
				Description:    "Something",
				IsRecorded:     true,
				Bpm:            &[]uint{120}[0],
				SongsterrLink:  &[]string{"http://songsterr.com/some-song"}[0],
				YoutubeLink:    &[]string{"https://www.youtube.com/watch?v=IHgFJEJgUrg&t=120s"}[0],
				ReleaseDate:    &[]internal.Date{internal.Date(time.Now())}[0],
				Difficulty:     &[]enums.Difficulty{enums.Easy}[0],
				GuitarTuningID: &[]uuid.UUID{uuid.New()}[0],
				ArtistID:       &[]uuid.UUID{uuid.New()}[0],
				AlbumID:        &[]uuid.UUID{uuid.New()}[0],
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

func TestValidateUpdateSongRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.UpdateSongRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// ID Test Cases
		{
			"ID is invalid because it's required",
			requests.UpdateSongRequest{ID: uuid.Nil, Title: validSongTitle},
			"ID",
			"required",
		},
		// Title Test Cases
		{
			"Title is invalid because it's required",
			requests.UpdateSongRequest{ID: uuid.New(), Title: ""},
			"Title",
			"required",
		},
		{
			"Title is invalid because it has more than 100 characters",
			requests.UpdateSongRequest{ID: uuid.New(), Title: strings.Repeat("a", 101)},
			"Title",
			"max",
		},
		// SongsterrLink Test Cases
		{
			"Songsterr Link is invalid because it is not an url",
			requests.UpdateSongRequest{
				ID:            uuid.New(),
				Title:         validSongTitle,
				SongsterrLink: &[]string{"scom"}[0],
			},
			"SongsterrLink",
			"url",
		},
		{
			"Songsterr Link is invalid because it is not a songsterr link",
			requests.UpdateSongRequest{
				ID:            uuid.New(),
				Title:         validSongTitle,
				SongsterrLink: &[]string{"http://google.com"}[0],
			},
			"SongsterrLink",
			"contains",
		},
		// YoutubeLink Test Cases
		{
			"Youtube Link is invalid because it is a youtube link",
			requests.UpdateSongRequest{
				ID:          uuid.New(),
				Title:       validSongTitle,
				YoutubeLink: &[]string{"https://google.com"}[0],
			},
			"YoutubeLink",
			"youtube_link",
		},
		// Difficulty Test Cases
		{
			"Difficulty is invalid because it is not a Difficulty Enum",
			requests.UpdateSongRequest{
				ID:         uuid.New(),
				Title:      validSongTitle,
				Difficulty: &[]enums.Difficulty{"Something else"}[0],
			},
			"Difficulty",
			"difficulty_enum",
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
			assert.Contains(t, errCode.Error.Error(), "UpdateSongRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateUpdateSongSettingsRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	tests := []struct {
		name    string
		request requests.UpdateSongSettingsRequest
	}{
		{
			"Non Optional",
			requests.UpdateSongSettingsRequest{
				SettingsID: uuid.New(),
			},
		},
		{
			"All Filled",
			requests.UpdateSongSettingsRequest{
				SettingsID:          uuid.New(),
				DefaultInstrumentID: &[]uuid.UUID{uuid.New()}[0],
				DefaultBandMemberID: &[]uuid.UUID{uuid.New()}[0],
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

func TestValidateUpdateSongSettingsRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.UpdateSongSettingsRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// Settings ID Test Cases
		{
			"Settings ID is invalid because it's required",
			requests.UpdateSongSettingsRequest{SettingsID: uuid.Nil},
			"SettingsID",
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
			assert.Contains(t, errCode.Error.Error(), "UpdateSongSettingsRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateBulkDeleteSongsRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.BulkDeleteSongsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)
}

func TestValidateBulkDeleteSongsRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.BulkDeleteSongsRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// IDs Test Cases
		{
			"IDs is invalid because it requires at least 1 id",
			requests.BulkDeleteSongsRequest{IDs: []uuid.UUID{}},
			"IDs",
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
			assert.Contains(t, errCode.Error.Error(), "BulkDeleteSongsRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}
