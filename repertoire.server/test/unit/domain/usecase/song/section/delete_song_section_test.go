package section

import (
	"errors"
	"math"
	"net/http"
	"repertoire/server/domain/usecase/song/section"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteSongSection_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewDeleteSongSection(songRepository)

	id := uuid.New()
	songID := uuid.New()

	internalError := errors.New("internal error")
	songRepository.On("GetWithSections", new(model.Song), songID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestDeleteSongSection_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewDeleteSongSection(songRepository)

	id := uuid.New()
	songID := uuid.New()

	songRepository.On("GetWithSections", new(model.Song), songID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestDeleteSongSection_WhenSectionIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewDeleteSongSection(songRepository)

	id := uuid.New()
	songID := uuid.New()

	// given - mocking
	song := &model.Song{
		ID: songID,
		Sections: []model.SongSection{
			{ID: uuid.New(), Order: 0},
		},
	}
	songRepository.On("GetWithSections", new(model.Song), songID).
		Return(nil, song).
		Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song section not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestDeleteSongSection_WhenUpdateSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewDeleteSongSection(songRepository)

	id := uuid.New()
	songID := uuid.New()

	// given - mocking
	song := &model.Song{
		ID: songID,
		Sections: []model.SongSection{
			{ID: id, Order: 0},
		},
	}
	songRepository.On("GetWithSections", new(model.Song), songID).
		Return(nil, song).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("UpdateWithAssociations", mock.IsType(song)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestDeleteSongSection_WhenDeleteSectionFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewDeleteSongSection(songRepository)

	id := uuid.New()
	songID := uuid.New()

	// given - mocking
	song := &model.Song{
		ID: songID,
		Sections: []model.SongSection{
			{ID: id, Order: 0},
		},
	}
	songRepository.On("GetWithSections", new(model.Song), songID).
		Return(nil, song).
		Once()

	songRepository.On("UpdateWithAssociations", mock.IsType(song)).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("DeleteSection", id).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestDeleteSongSection_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	tests := []struct {
		name                   string
		song                   model.Song
		sectionsIndex          uint
		expectedSongConfidence float64
		expectedSongRehearsals float64
		expectedSongProgress   float64
	}{
		{
			"1 - When it was the only section",
			model.Song{
				ID: uuid.New(),
				Sections: []model.SongSection{
					{ID: uuid.New(), Order: 0},
				},
			},
			0,
			0,
			0,
			0,
		},
		{
			"2 - When it was the only section with stats",
			model.Song{
				ID: uuid.New(),
				Sections: []model.SongSection{
					{ID: uuid.New(), Order: 0, Confidence: 55, Rehearsals: 25, Progress: 39},
				},
				Confidence: 55,
				Rehearsals: 25,
				Progress:   39,
			},
			0,
			0,
			0,
			0,
		},
		{
			"3 - When there are more sections, but not stats",
			model.Song{
				ID: uuid.New(),
				Sections: []model.SongSection{
					{ID: uuid.New(), Order: 0},
					{ID: uuid.New(), Order: 1},
					{ID: uuid.New(), Order: 2},
					{ID: uuid.New(), Order: 3},
					{ID: uuid.New(), Order: 4},
				},
			},
			2,
			0,
			0,
			0,
		},
		{
			"4 - When there are more sections with stats",
			model.Song{
				ID: uuid.New(),
				Sections: []model.SongSection{
					{ID: uuid.New(), Order: 0, Confidence: 55, Rehearsals: 12, Progress: 45},
					{ID: uuid.New(), Order: 1, Confidence: 23, Rehearsals: 5, Progress: 15},
					{ID: uuid.New(), Order: 2, Confidence: 78, Rehearsals: 25, Progress: 100},
					{ID: uuid.New(), Order: 3, Confidence: 40, Rehearsals: 6, Progress: 63},
					{ID: uuid.New(), Order: 4, Confidence: 80, Rehearsals: 19, Progress: 170},
				},
				Confidence: 55.2,
				Rehearsals: 13.4,
				Progress:   78.6,
			},
			2,
			50,
			11,
			73,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			songRepository := new(repository.SongRepositoryMock)
			_uut := section.NewDeleteSongSection(songRepository)

			id := tt.song.Sections[tt.sectionsIndex].ID
			songID := tt.song.ID

			// given - mocking
			songRepository.On("GetWithSections", new(model.Song), songID).
				Return(nil, &tt.song).
				Once()

			songRepository.On("UpdateWithAssociations", mock.IsType(&tt.song)).
				Run(func(args mock.Arguments) {
					newSong := args.Get(0).(*model.Song)

					// stats updated
					assert.Equal(t, tt.expectedSongConfidence, math.Round(newSong.Confidence))
					assert.Equal(t, tt.expectedSongRehearsals, math.Round(newSong.Rehearsals))
					assert.Equal(t, tt.expectedSongProgress, math.Round(newSong.Progress))

					// sections ordered
					sections := slices.Clone(newSong.Sections)
					sections = slices.DeleteFunc(sections, func(a model.SongSection) bool {
						return a.ID == id
					})
					for i, s := range sections {
						assert.Equal(t, uint(i), s.Order)
					}
				}).
				Return(nil).
				Once()

			songRepository.On("DeleteSection", id).Return(nil).Once()

			// when
			errCode := _uut.Handle(id, songID)

			// then
			assert.Nil(t, errCode)

			songRepository.AssertExpectations(t)
		})
	}
}
