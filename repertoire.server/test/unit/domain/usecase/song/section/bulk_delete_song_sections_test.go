package section

import (
	"errors"
	"math"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song/section"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBulkDeleteSongSections_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewBulkDeleteSongSections(songRepository)

	request := requests.BulkDeleteSongSectionsRequest{
		SongID: uuid.New(),
	}

	internalError := errors.New("internal error")
	songRepository.On("GetWithSections", new(model.Song), request.SongID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestBulkDeleteSongSections_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewBulkDeleteSongSections(songRepository)

	request := requests.BulkDeleteSongSectionsRequest{
		SongID: uuid.New(),
	}

	songRepository.On("GetWithSections", new(model.Song), request.SongID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestBulkDeleteSongSections_WhenSectionsAreNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewBulkDeleteSongSections(songRepository)

	request := requests.BulkDeleteSongSectionsRequest{
		IDs:    []uuid.UUID{uuid.New()},
		SongID: uuid.New(),
	}

	song := &model.Song{
		ID: request.SongID,
		Sections: []model.SongSection{
			{ID: uuid.New(), Order: 0},
		},
	}
	songRepository.On("GetWithSections", new(model.Song), request.SongID).
		Return(nil, song).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song sections not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestBulkDeleteSongSections_WhenUpdateSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewBulkDeleteSongSections(songRepository)

	request := requests.BulkDeleteSongSectionsRequest{
		IDs:    []uuid.UUID{uuid.New()},
		SongID: uuid.New(),
	}

	// given - mocking
	song := &model.Song{
		ID: request.SongID,
		Sections: []model.SongSection{
			{ID: request.IDs[0], Order: 0},
		},
	}
	songRepository.On("GetWithSections", new(model.Song), request.SongID).
		Return(nil, song).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("UpdateWithAssociations", mock.IsType(song)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestBulkDeleteSongSections_WhenDeleteSectionFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewBulkDeleteSongSections(songRepository)

	request := requests.BulkDeleteSongSectionsRequest{
		IDs:    []uuid.UUID{uuid.New()},
		SongID: uuid.New(),
	}

	// given - mocking
	song := &model.Song{
		ID: request.SongID,
		Sections: []model.SongSection{
			{ID: request.IDs[0], Order: 0},
		},
	}
	songRepository.On("GetWithSections", new(model.Song), request.SongID).
		Return(nil, song).
		Once()

	songRepository.On("UpdateWithAssociations", mock.IsType(song)).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("DeleteSections", request.IDs).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestBulkDeleteSongSections_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	tests := []struct {
		name                   string
		song                   model.Song
		sectionIndexes         []uint
		expectedSongConfidence float64
		expectedSongRehearsals float64
		expectedSongProgress   float64
	}{
		{
			"1 - When it was the only sections",
			model.Song{
				ID: uuid.New(),
				Sections: []model.SongSection{
					{ID: uuid.New(), Order: 0},
					{ID: uuid.New(), Order: 1},
				},
			},
			[]uint{0},
			0,
			0,
			0,
		},
		{
			"2 - When it was the only sections with stats",
			model.Song{
				ID: uuid.New(),
				Sections: []model.SongSection{
					{ID: uuid.New(), Order: 0, Confidence: 55, Rehearsals: 30, Progress: 100},
					{ID: uuid.New(), Order: 0, Confidence: 50, Rehearsals: 15, Progress: 50},
				},
				Confidence: 53,
				Rehearsals: 20,
				Progress:   75,
			},
			[]uint{0, 1},
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
			[]uint{0, 2},
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
			[]uint{0, 2},
			48,
			10,
			83,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			songRepository := new(repository.SongRepositoryMock)
			_uut := section.NewBulkDeleteSongSections(songRepository)

			request := requests.BulkDeleteSongSectionsRequest{
				SongID: tt.song.ID,
				IDs:    []uuid.UUID{},
			}
			for _, songSectionIndex := range tt.sectionIndexes {
				request.IDs = append(request.IDs, tt.song.Sections[songSectionIndex].ID)
			}

			// given - mocking
			songRepository.On("GetWithSections", new(model.Song), request.SongID).
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
					order := uint(0)
					for _, s := range tt.song.Sections {
						if slices.ContainsFunc(request.IDs, func(id uuid.UUID) bool {
							return id == s.ID
						}) {
							continue
						}
						assert.Equal(t, order, s.Order)
						order++
					}
				}).
				Return(nil).
				Once()

			songRepository.On("DeleteSections", request.IDs).Return(nil).Once()

			// when
			errCode := _uut.Handle(request)

			// then
			assert.Nil(t, errCode)

			songRepository.AssertExpectations(t)
		})
	}
}
