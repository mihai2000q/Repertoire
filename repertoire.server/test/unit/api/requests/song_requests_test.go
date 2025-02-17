package requests

import (
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/api/validation"
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
				OrderBy:     []string{"title asc", "created_at desc"},
				SearchBy:    []string{"title = something", "is_recorded <> false"},
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
				ReleaseDate:    &[]time.Time{time.Now()}[0],
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
			"All Null",
			requests.UpdateSongRequest{
				ID:    uuid.New(),
				Title: validSongTitle,
			},
		},
		{
			"Nothing Null",
			requests.UpdateSongRequest{
				ID:             uuid.New(),
				Title:          validSongTitle,
				Description:    "Something",
				IsRecorded:     true,
				Bpm:            &[]uint{120}[0],
				SongsterrLink:  &[]string{"http://songsterr.com/some-song"}[0],
				YoutubeLink:    &[]string{"https://www.youtube.com/watch?v=IHgFJEJgUrg&pp=ygUMeW91ciBiZXRyYXlh"}[0],
				ReleaseDate:    &[]time.Time{time.Now()}[0],
				Difficulty:     &[]enums.Difficulty{enums.Easy}[0],
				GuitarTuningID: &[]uuid.UUID{uuid.New()}[0],
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

// Sections

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
