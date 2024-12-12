package section

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song/section"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/domain/processor"
	"testing"
)

func TestUpdateSongSection_WhenGetSectionsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewUpdateSongSection(songRepository, nil)

	request := requests.UpdateSongSectionRequest{
		ID:     uuid.New(),
		Name:   "Some Artist",
		TypeID: uuid.New(),
	}

	internalError := errors.New("internal error")
	songRepository.On("GetSection", new(model.SongSection), request.ID).
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

func TestUpdateSongSection_WhenSectionsIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewUpdateSongSection(songRepository, nil)

	request := requests.UpdateSongSectionRequest{
		ID:     uuid.New(),
		Name:   "Some Artist",
		TypeID: uuid.New(),
	}

	songRepository.On("GetSection", new(model.SongSection), request.ID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song section not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestUpdateSongSection_WhenRehearsalsIsDecreasing_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewUpdateSongSection(songRepository, nil)

	request := requests.UpdateSongSectionRequest{
		ID:         uuid.New(),
		Name:       "Some Artist",
		Rehearsals: 23,
		TypeID:     uuid.New(),
	}

	mockSection := &model.SongSection{
		ID:         request.ID,
		Rehearsals: request.Rehearsals + 1,
	}
	songRepository.On("GetSection", new(model.SongSection), request.ID).
		Return(nil, mockSection).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusBadRequest, errCode.Code)
	assert.Equal(t, "rehearsals can only be increased", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestUpdateSongSection_WhenCreateHistoryFails_ShouldReturnInternalServerError(t *testing.T) {
	tests := []struct {
		name     string
		request  requests.UpdateSongSectionRequest
		property model.SongSectionProperty
	}{
		{
			"When updating the confidence level",
			requests.UpdateSongSectionRequest{
				ID:         uuid.New(),
				Name:       "Some Section",
				Confidence: 50,
				TypeID:     uuid.New(),
			},
			model.ConfidenceProperty,
		},
		{
			"When updating the number of rehearsals",
			requests.UpdateSongSectionRequest{
				ID:         uuid.New(),
				Name:       "Some Section",
				Rehearsals: 50,
				TypeID:     uuid.New(),
			},
			model.RehearsalsProperty,
		},
		{
			"When updating the number of rehearsals and confidence, the rehearsals will fail",
			requests.UpdateSongSectionRequest{
				ID:         uuid.New(),
				Name:       "Some Section",
				Rehearsals: 50,
				Confidence: 64,
				TypeID:     uuid.New(),
			},
			model.RehearsalsProperty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			songRepository := new(repository.SongRepositoryMock)
			_uut := section.NewUpdateSongSection(songRepository, nil)

			// given - mocking
			mockSection := &model.SongSection{
				ID:   tt.request.ID,
				Name: "Old name",
			}
			songRepository.On("GetSection", new(model.SongSection), tt.request.ID).
				Return(nil, mockSection).
				Once()

			internalError := errors.New("internal error")
			songRepository.On("CreateSongSectionHistory", mock.IsType(new(model.SongSectionHistory))).
				Run(func(args mock.Arguments) {
					newHistory := args.Get(0).(*model.SongSectionHistory)
					assert.Equal(t, tt.property, newHistory.Property)
				}).
				Return(internalError).
				Once()

			// when
			errCode := _uut.Handle(tt.request)

			// then
			assert.NotNil(t, errCode)
			assert.Equal(t, http.StatusInternalServerError, errCode.Code)
			assert.Equal(t, internalError, errCode.Error)

			songRepository.AssertExpectations(t)
		})
	}
}

func TestUpdateSongSection_WhenGetHistoryFails_ShouldReturnInternalServerError(t *testing.T) {
	tests := []struct {
		name     string
		request  requests.UpdateSongSectionRequest
		property model.SongSectionProperty
	}{
		{
			"When updating the confidence level",
			requests.UpdateSongSectionRequest{
				ID:         uuid.New(),
				Name:       "Some Section",
				Confidence: 50,
				TypeID:     uuid.New(),
			},
			model.ConfidenceProperty,
		},
		{
			"When updating the number of rehearsals",
			requests.UpdateSongSectionRequest{
				ID:         uuid.New(),
				Name:       "Some Section",
				Rehearsals: 50,
				TypeID:     uuid.New(),
			},
			model.RehearsalsProperty,
		},
		{
			"When updating the number of rehearsals and confidence, the rehearsals will fail",
			requests.UpdateSongSectionRequest{
				ID:         uuid.New(),
				Name:       "Some Section",
				Rehearsals: 50,
				Confidence: 64,
				TypeID:     uuid.New(),
			},
			model.RehearsalsProperty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			songRepository := new(repository.SongRepositoryMock)
			_uut := section.NewUpdateSongSection(songRepository, nil)

			// given - mocking
			mockSection := &model.SongSection{
				ID:   tt.request.ID,
				Name: "Old name",
			}
			songRepository.On("GetSection", new(model.SongSection), tt.request.ID).
				Return(nil, mockSection).
				Once()

			songRepository.On("CreateSongSectionHistory", mock.IsType(new(model.SongSectionHistory))).
				Return(nil).
				Once()

			internalError := errors.New("internal error")
			songRepository.
				On(
					"GetSongSectionHistory",
					mock.IsType(new([]model.SongSectionHistory)),
					mockSection.ID,
					tt.property,
				).
				Return(internalError).
				Once()

			// when
			errCode := _uut.Handle(tt.request)

			// then
			assert.NotNil(t, errCode)
			assert.Equal(t, http.StatusInternalServerError, errCode.Code)
			assert.Equal(t, internalError, errCode.Error)

			songRepository.AssertExpectations(t)
		})
	}
}

func TestUpdateSongSection_WhenUpdateSectionFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewUpdateSongSection(songRepository, nil)

	request := requests.UpdateSongSectionRequest{
		ID:     uuid.New(),
		Name:   "Some Artist",
		TypeID: uuid.New(),
	}

	// given - mocking
	mockSection := &model.SongSection{
		ID:   request.ID,
		Name: "Old name",
	}
	songRepository.On("GetSection", new(model.SongSection), request.ID).
		Return(nil, mockSection).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("UpdateSection", mock.IsType(new(model.SongSection))).
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

func TestUpdateSongSection_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	tests := []struct {
		name    string
		request requests.UpdateSongSectionRequest
	}{
		{
			"without confidence or rehearsals",
			requests.UpdateSongSectionRequest{
				ID:     uuid.New(),
				Name:   "Some New Name",
				TypeID: uuid.New(),
			},
		},
		{
			"with Confidence",
			requests.UpdateSongSectionRequest{
				ID:         uuid.New(),
				Name:       "Some New Name",
				Confidence: 50,
				TypeID:     uuid.New(),
			},
		},
		{
			"with rehearsals",
			requests.UpdateSongSectionRequest{
				ID:         uuid.New(),
				Name:       "Some New Name",
				Rehearsals: 12,
				TypeID:     uuid.New(),
			},
		},
		{
			"with confidence and rehearsals",
			requests.UpdateSongSectionRequest{
				ID:         uuid.New(),
				Name:       "Some New Name",
				Confidence: 50,
				Rehearsals: 12,
				TypeID:     uuid.New(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			songRepository := new(repository.SongRepositoryMock)
			progressProcessor := new(processor.ProgressProcessorMock)
			_uut := section.NewUpdateSongSection(songRepository, progressProcessor)

			// given - mocking
			mockSection := &model.SongSection{
				ID:   tt.request.ID,
				Name: "Old name",
			}
			songRepository.On("GetSection", new(model.SongSection), tt.request.ID).
				Return(nil, mockSection).
				Once()

			var history []model.SongSectionHistory
			songSectionHistoryTimes := 0

			var rehearsalScore uint64
			var confidenceScore uint

			if mockSection.Rehearsals != tt.request.Rehearsals {
				songSectionHistoryTimes++
				rehearsalScore = 125
				progressProcessor.On("ComputeRehearsalsScore", history).Return(rehearsalScore).Once()
			}

			if mockSection.Confidence != tt.request.Confidence {
				songSectionHistoryTimes++
				confidenceScore = 88
				progressProcessor.On("ComputeConfidenceScore", history).Return(confidenceScore).Once()
			}

			if songSectionHistoryTimes > 0 {
				songRepository.On("CreateSongSectionHistory", mock.IsType(new(model.SongSectionHistory))).
					Run(func(args mock.Arguments) {
						newHistory := args.Get(0).(*model.SongSectionHistory)

						assert.NotEmpty(t, newHistory.ID)
						assert.Equal(t, mockSection.ID, newHistory.SongSectionID)
						assert.NotEmpty(t, newHistory.Property)
						if newHistory.Property == model.ConfidenceProperty {
							assert.Equal(t, mockSection.Confidence, newHistory.From)
							assert.Equal(t, tt.request.Confidence, newHistory.To)
						} else if newHistory.Property == model.RehearsalsProperty {
							assert.Equal(t, mockSection.Rehearsals, newHistory.From)
							assert.Equal(t, tt.request.Rehearsals, newHistory.To)
						} else {
							assert.Fail(t, "invalid property")
						}
					}).
					Return(nil).
					Times(songSectionHistoryTimes)

				songRepository.
					On(
						"GetSongSectionHistory",
						mock.IsType(new([]model.SongSectionHistory)),
						mockSection.ID,
						mock.IsType(model.ConfidenceProperty),
					).
					Run(func(args mock.Arguments) {
						property := args.Get(2).(model.SongSectionProperty)
						assert.True(t, property == model.ConfidenceProperty || property == model.RehearsalsProperty)
					}).
					Return(nil, &history).
					Times(songSectionHistoryTimes)
			}

			songRepository.On("UpdateSection", mock.IsType(new(model.SongSection))).
				Run(func(args mock.Arguments) {
					newSection := args.Get(0).(*model.SongSection)
					assertUpdatedSongSection(t, tt.request, *newSection, rehearsalScore, confidenceScore)
				}).
				Return(nil).
				Once()

			// when
			errCode := _uut.Handle(tt.request)

			// then
			assert.Nil(t, errCode)

			songRepository.AssertExpectations(t)
			progressProcessor.AssertExpectations(t)
		})
	}
}

func assertUpdatedSongSection(
	t *testing.T,
	request requests.UpdateSongSectionRequest,
	section model.SongSection,
	rehearsalsScore uint64,
	confidenceScore uint,
) {
	assert.Equal(t, request.Name, section.Name)
	assert.Equal(t, request.Confidence, section.Confidence)
	assert.Equal(t, request.Rehearsals, section.Rehearsals)
	assert.Equal(t, rehearsalsScore, section.RehearsalsScore)
	assert.Equal(t, confidenceScore, section.ConfidenceScore)
	assert.Equal(t, section.RehearsalsScore*uint64(section.ConfidenceScore), section.Progress)
	assert.Equal(t, request.TypeID, section.SongSectionTypeID)
}
