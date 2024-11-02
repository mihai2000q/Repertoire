package requests

import (
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
			// given
			_uut := validation.NewValidator(nil)

			// when
			errCode := _uut.Validate(tt.request)

			// then
			assert.NotNil(t, errCode)
			assert.Len(t, errCode.Error, 1)
			assert.Contains(t, errCode.Error.Error(), "GetSongsRequest."+tt.expectedInvalidField)
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
			"Nothing Null 1",
			CreateSongRequest{
				Title:          validSongTitle,
				Description:    "Something",
				Bpm:            &[]uint{12}[0],
				SongsterrLink:  &[]string{"http://songsterr.com/some-song"}[0],
				GuitarTuningID: &[]uuid.UUID{uuid.New()}[0],
				AlbumID:        &[]uuid.UUID{uuid.New()}[0],
				ArtistID:       &[]uuid.UUID{uuid.New()}[0],
				Sections: []CreateSectionRequest{
					{Name: "A section", TypeID: uuid.New()},
					{Name: "A Second Section", TypeID: uuid.New()},
				},
			},
		},
		{
			"Nothing Null 2",
			CreateSongRequest{
				Title:          validSongTitle,
				Description:    "Something",
				Bpm:            &[]uint{12}[0],
				SongsterrLink:  &[]string{"https://songsterr.com/some-other"}[0],
				ReleaseDate:    &[]time.Time{time.Now()}[0],
				Difficulty:     &[]enums.Difficulty{enums.Easy}[0],
				GuitarTuningID: &[]uuid.UUID{uuid.New()}[0],
				AlbumTitle:     &[]string{"New Album Title"}[0],
				ArtistName:     &[]string{"New Artist Name"}[0],
				Sections: []CreateSectionRequest{
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
		request               CreateSongRequest
		expectedInvalidFields []string
		expectedFailedTags    []string
	}{
		// Title Test Cases
		{
			"Title is invalid because it's required",
			CreateSongRequest{Title: ""},
			[]string{"Title"},
			[]string{"required"},
		},
		{
			"Title is invalid because it has more than 100 characters",
			CreateSongRequest{Title: strings.Repeat("a", 101)},
			[]string{"Title"},
			[]string{"max"},
		},
		// SongsterrLink Test Cases
		{
			"Songsterr Link is invalid because it is not an url",
			CreateSongRequest{
				Title:         validSongTitle,
				SongsterrLink: &[]string{"scom"}[0],
			},
			[]string{"SongsterrLink"},
			[]string{"url"},
		},
		{
			"Songsterr Link is invalid because it is not a songsterr link",
			CreateSongRequest{
				Title:         validSongTitle,
				SongsterrLink: &[]string{"http://google.com"}[0],
			},
			[]string{"SongsterrLink"},
			[]string{"contains"},
		},
		// Difficulty Test Cases
		{
			"Difficulty is invalid because it is not a Difficulty Enum",
			CreateSongRequest{
				Title:      validSongTitle,
				Difficulty: &[]enums.Difficulty{"Something else"}[0],
			},
			[]string{"Difficulty"},
			[]string{"isDifficultyEnum"},
		},
		// Album ID and Album Title Test Case
		{
			"Album Title and ID are invalid because only one can be set at a time",
			CreateSongRequest{
				Title:      validSongTitle,
				AlbumID:    &[]uuid.UUID{uuid.New()}[0],
				AlbumTitle: &[]string{"New Album Title"}[0],
			},
			[]string{"AlbumID", "AlbumTitle"},
			[]string{"excluded_with", "excluded_with"},
		},
		// Artist ID and Artist Title Test Case
		{
			"Artist Name and ID are invalid because only one can be set at a time",
			CreateSongRequest{
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
			CreateSongRequest{
				Title: validSongTitle,
				Sections: []CreateSectionRequest{
					{Name: "", TypeID: uuid.New()},
				},
			},
			[]string{"Sections[0].Name"},
			[]string{"required"},
		},
		// Sections - Name Test Cases
		{
			"Sections are invalid because the first element has a Name with too many characters",
			CreateSongRequest{
				Title: validSongTitle,
				Sections: []CreateSectionRequest{
					{Name: strings.Repeat("a", 31), TypeID: uuid.New()},
				},
			},
			[]string{"Sections[0].Name"},
			[]string{"max"},
		},
		// Sections - Type ID Test Cases
		{
			"Sections are invalid because the first element has an empty Type ID",
			CreateSongRequest{
				Title: validSongTitle,
				Sections: []CreateSectionRequest{
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

			err := errCode.Error.Error()

			// then
			assert.NotNil(t, errCode)
			assert.Len(t, tt.expectedFailedTags, len(tt.expectedInvalidFields))
			assert.Len(t, errCode.Error, len(tt.expectedFailedTags))
			for _, expectedInvalidField := range tt.expectedInvalidFields {
				assert.Contains(t, err, "CreateSongRequest."+expectedInvalidField)
			}
			for _, expectedFailedTag := range tt.expectedFailedTags {
				assert.Contains(t, errCode.Error.Error(), "'"+expectedFailedTag+"' tag")
			}
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
				ID:             uuid.New(),
				Title:          validSongTitle,
				Description:    "Something",
				IsRecorded:     true,
				Bpm:            &[]uint{120}[0],
				SongsterrLink:  &[]string{"http://songsterr.com/some-song"}[0],
				ReleaseDate:    &[]time.Time{time.Now()}[0],
				Difficulty:     &[]enums.Difficulty{enums.Easy}[0],
				GuitarTuningID: &[]uuid.UUID{uuid.New()}[0],
				AlbumID:        &[]uuid.UUID{uuid.New()}[0],
				ArtistID:       &[]uuid.UUID{uuid.New()}[0],
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
				SongsterrLink: &[]string{"scom"}[0],
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
		// Difficulty Test Cases
		{
			"Difficulty is invalid because it is not a Difficulty Enum",
			UpdateSongRequest{
				ID:         uuid.New(),
				Title:      validSongTitle,
				Difficulty: &[]enums.Difficulty{"Something else"}[0],
			},
			"Difficulty",
			"isDifficultyEnum",
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
			assert.Equal(t, 400, errCode.Code)
		})
	}
}

// Sections

var validSectionName = "James Solo"

func TestValidateCreateSongSectionRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := CreateSongSectionRequest{
		SongID: uuid.New(),
		Name:   validSectionName,
		TypeID: uuid.New(),
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)
}

func TestValidateCreateSongSectionRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              CreateSongSectionRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// Song ID Test Cases
		{
			"Song ID is invalid because it's required",
			CreateSongSectionRequest{
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
			CreateSongSectionRequest{
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
			CreateSongSectionRequest{
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
			CreateSongSectionRequest{
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
			assert.Equal(t, 400, errCode.Code)
		})
	}
}

func TestValidateUpdateSongSectionRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := UpdateSongSectionRequest{
		ID:         uuid.New(),
		Name:       validSectionName,
		Rehearsals: 23,
		TypeID:     uuid.New(),
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)
}

func TestValidateUpdateSongSectionRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              UpdateSongSectionRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// ID Test Cases
		{
			"Song ID is invalid because it's required",
			UpdateSongSectionRequest{
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
			UpdateSongSectionRequest{
				ID:     uuid.New(),
				Name:   "",
				TypeID: uuid.New(),
			},
			"Name",
			"required",
		},
		// Name Test Cases
		{
			"Name is invalid because it has too many characters",
			UpdateSongSectionRequest{
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
			UpdateSongSectionRequest{
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
			assert.Equal(t, 400, errCode.Code)
		})
	}
}

func TestValidateMoveSongSectionRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := MoveSongSectionRequest{
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
		request              MoveSongSectionRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// ID Test Cases
		{
			"ID is invalid because it's required",
			MoveSongSectionRequest{ID: uuid.Nil, OverID: uuid.New(), SongID: uuid.New()},
			"ID",
			"required",
		},
		// Over ID Test Cases
		{
			"Over ID is invalid because it's required",
			MoveSongSectionRequest{ID: uuid.New(), OverID: uuid.Nil, SongID: uuid.New()},
			"OverID",
			"required",
		},
		// Song ID Test Cases
		{
			"Song ID is invalid because it's required",
			MoveSongSectionRequest{ID: uuid.New(), OverID: uuid.New(), SongID: uuid.Nil},
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
			assert.Equal(t, 400, errCode.Code)
		})
	}
}
