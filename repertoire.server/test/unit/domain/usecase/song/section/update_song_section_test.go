package section

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song/section"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/domain/processor"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestUpdateSongSection_WhenRehearsalsIsDecreasing_ShouldReturnConflictError(t *testing.T) {
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
	assert.Equal(t, http.StatusConflict, errCode.Code)
	assert.Equal(t, "rehearsals can only be increased", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestUpdateSongSection_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	tests := []struct {
		name    string
		request requests.UpdateSongSectionRequest
	}{
		{
			"When updating the confidence level",
			requests.UpdateSongSectionRequest{
				ID:         uuid.New(),
				Name:       "Some Section",
				Confidence: 50,
				TypeID:     uuid.New(),
			},
		},
		{
			"When updating the number of rehearsals",
			requests.UpdateSongSectionRequest{
				ID:         uuid.New(),
				Name:       "Some Section",
				Rehearsals: 50,
				TypeID:     uuid.New(),
			},
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			songRepository := new(repository.SongRepositoryMock)
			_uut := section.NewUpdateSongSection(songRepository, nil)

			// given - mocking
			mockSection := &model.SongSection{
				ID:     tt.request.ID,
				Name:   "Old name",
				SongID: uuid.New(),
			}
			songRepository.On("GetSection", new(model.SongSection), tt.request.ID).
				Return(nil, mockSection).
				Once()

			internalError := errors.New("internal error")
			songRepository.On("Get", new(model.Song), mockSection.SongID).
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

func TestUpdateSongSection_WhenCountSectionsFails_ShouldReturnInternalServerError(t *testing.T) {
	tests := []struct {
		name    string
		request requests.UpdateSongSectionRequest
	}{
		{
			"When updating the confidence level",
			requests.UpdateSongSectionRequest{
				ID:         uuid.New(),
				Name:       "Some Section",
				Confidence: 50,
				TypeID:     uuid.New(),
			},
		},
		{
			"When updating the number of rehearsals",
			requests.UpdateSongSectionRequest{
				ID:         uuid.New(),
				Name:       "Some Section",
				Rehearsals: 50,
				TypeID:     uuid.New(),
			},
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			songRepository := new(repository.SongRepositoryMock)
			_uut := section.NewUpdateSongSection(songRepository, nil)

			// given - mocking
			mockSection := &model.SongSection{
				ID:     tt.request.ID,
				Name:   "Old name",
				SongID: uuid.New(),
			}
			songRepository.On("GetSection", new(model.SongSection), tt.request.ID).
				Return(nil, mockSection).
				Once()

			mockSong := &model.Song{}
			songRepository.On("Get", new(model.Song), mockSection.SongID).
				Return(nil, mockSong).
				Once()

			internalError := errors.New("internal error")
			songRepository.On("CountSectionsBySong", mock.Anything, mockSection.SongID).
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

func TestUpdateSongSection_WhenIsBandMemberAssociatedWithSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewUpdateSongSection(songRepository, nil)

	request := requests.UpdateSongSectionRequest{
		ID:           uuid.New(),
		Name:         "Some Section",
		TypeID:       uuid.New(),
		BandMemberID: &[]uuid.UUID{uuid.New()}[0],
	}

	// given - mocking
	mockSection := &model.SongSection{
		ID:     uuid.New(),
		Name:   "Old name",
		SongID: uuid.New(),
	}
	songRepository.On("GetSection", new(model.SongSection), request.ID).
		Return(nil, mockSection).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("IsBandMemberAssociatedWithSong", mockSection.SongID, *request.BandMemberID).
		Return(false, internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestUpdateSongSection_WhenBandMemberIsNotAssociatedWithSong_ShouldReturnConflictError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewUpdateSongSection(songRepository, nil)

	request := requests.UpdateSongSectionRequest{
		ID:           uuid.New(),
		Name:         "Some Section",
		TypeID:       uuid.New(),
		BandMemberID: &[]uuid.UUID{uuid.New()}[0],
	}

	// given - mocking
	mockSection := &model.SongSection{
		ID:     uuid.New(),
		Name:   "Old name",
		SongID: uuid.New(),
	}
	songRepository.On("GetSection", new(model.SongSection), request.ID).
		Return(nil, mockSection).
		Once()

	songRepository.On("IsBandMemberAssociatedWithSong", mockSection.SongID, *request.BandMemberID).
		Return(false, nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusConflict, errCode.Code)
	assert.Equal(t, "band member is not part of the artist associated with this song", errCode.Error.Error())

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
				ID:     tt.request.ID,
				Name:   "Old name",
				SongID: uuid.New(),
			}
			songRepository.On("GetSection", new(model.SongSection), tt.request.ID).
				Return(nil, mockSection).
				Once()

			mockSong := &model.Song{}
			songRepository.On("Get", new(model.Song), mockSection.SongID).
				Return(nil, mockSong).
				Once()

			sectionsCount := &[]int64{20}[0]
			songRepository.On("CountSectionsBySong", mock.IsType(sectionsCount), mockSection.SongID).
				Return(nil, sectionsCount).
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
				ID:     tt.request.ID,
				Name:   "Old name",
				SongID: uuid.New(),
			}
			songRepository.On("GetSection", new(model.SongSection), tt.request.ID).
				Return(nil, mockSection).
				Once()

			mockSong := &model.Song{}
			songRepository.On("Get", new(model.Song), mockSection.SongID).
				Return(nil, mockSong).
				Once()

			sectionsCount := &[]int64{20}[0]
			songRepository.On("CountSectionsBySong", mock.IsType(sectionsCount), mockSection.SongID).
				Return(nil, sectionsCount).
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

func TestUpdateSongSection_WhenUpdateSongFails_ShouldReturnInternalServerError(t *testing.T) {
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
			progressProcessor := new(processor.ProgressProcessorMock)
			_uut := section.NewUpdateSongSection(songRepository, progressProcessor)

			// given - mocking
			mockSection := &model.SongSection{
				ID:     tt.request.ID,
				Name:   "Old name",
				SongID: uuid.New(),
			}
			songRepository.On("GetSection", new(model.SongSection), tt.request.ID).
				Return(nil, mockSection).
				Once()

			mockSong := &model.Song{
				ID: mockSection.SongID,
			}
			songRepository.On("Get", new(model.Song), mockSection.SongID).
				Return(nil, mockSong).
				Once()

			sectionsCount := &[]int64{20}[0]
			songRepository.On("CountSectionsBySong", mock.IsType(sectionsCount), mockSection.SongID).
				Return(nil, sectionsCount).
				Once()

			var history []model.SongSectionHistory
			songSectionHistoryTimes := 0

			if mockSection.Rehearsals != tt.request.Rehearsals {
				songSectionHistoryTimes++
				progressProcessor.On("ComputeRehearsalsScore", history).Return(uint64(125)).Once()
			}

			if mockSection.Confidence != tt.request.Confidence {
				songSectionHistoryTimes++
				progressProcessor.On("ComputeConfidenceScore", history).Return(uint(12)).Once()
			}

			if songSectionHistoryTimes > 0 {
				songRepository.On("CreateSongSectionHistory", mock.IsType(new(model.SongSectionHistory))).
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

				progressProcessor.On("ComputeProgress", mock.IsType(*mockSection)).
					Return(uint64(780)).
					Times(songSectionHistoryTimes)
			}

			internalError := errors.New("internal error")
			songRepository.On("Update", mock.IsType(mockSong)).Return(internalError).Once()

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

func TestUpdateSongSection_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	id := uuid.New()
	songID := uuid.New()

	tests := []struct {
		name                   string
		songSection            *model.SongSection
		request                requests.UpdateSongSectionRequest
		progress               uint64
		song                   *model.Song
		sectionsCount          *int64
		expectedSongConfidence float64
		expectedSongRehearsals float64
		expectedSongProgress   float64
	}{
		{
			"without confidence or rehearsals",
			&model.SongSection{
				ID:     id,
				Name:   "Old name",
				SongID: uuid.New(),
			},
			requests.UpdateSongSectionRequest{
				ID:     id,
				Name:   "Some New Name",
				TypeID: uuid.New(),
			},
			0,
			nil,
			nil,
			0,
			0,
			0,
		},
		{
			"with Confidence",
			&model.SongSection{
				ID:         id,
				Name:       "Old name",
				Confidence: 35,
				Progress:   7,
				SongID:     songID,
			},
			requests.UpdateSongSectionRequest{
				ID:         id,
				Name:       "Some New Name",
				Confidence: 50,
				TypeID:     uuid.New(),
			},
			8,
			&model.Song{
				ID:         songID,
				Confidence: 45.5,
				Progress:   13.5,
			},
			&[]int64{4}[0],
			49.25,
			0,
			13.75,
		},
		{
			"with Rehearsals",
			&model.SongSection{
				ID:         id,
				Name:       "Old name",
				Rehearsals: 35,
				Progress:   7,
				SongID:     uuid.New(),
			},
			requests.UpdateSongSectionRequest{
				ID:         id,
				Name:       "Some New Name",
				Rehearsals: 40,
				TypeID:     uuid.New(),
			},
			8,
			&model.Song{
				ID:         songID,
				Rehearsals: 27.5,
				Progress:   13.5,
			},
			&[]int64{4}[0],
			0,
			28.75,
			13.75,
		},
		{
			"with Rehearsals",
			&model.SongSection{
				ID:         id,
				Name:       "Old name",
				Confidence: 35,
				Rehearsals: 35,
				Progress:   7,
				SongID:     uuid.New(),
			},
			requests.UpdateSongSectionRequest{
				ID:         id,
				Name:       "Some New Name",
				Rehearsals: 40,
				Confidence: 50,
				TypeID:     uuid.New(),
			},
			8,
			&model.Song{
				ID:         songID,
				Rehearsals: 27.5,
				Confidence: 45.5,
				Progress:   13.5,
			},
			&[]int64{4}[0],
			49.25,
			28.75,
			13.75,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			songRepository := new(repository.SongRepositoryMock)
			progressProcessor := new(processor.ProgressProcessorMock)
			_uut := section.NewUpdateSongSection(songRepository, progressProcessor)

			// given - mocking
			songRepository.On("GetSection", new(model.SongSection), tt.request.ID).
				Return(nil, tt.songSection).
				Once()

			var history []model.SongSectionHistory
			songSectionHistoryTimes := 0

			var rehearsalScore uint64
			var confidenceScore uint

			if tt.songSection.Rehearsals != tt.request.Rehearsals {
				songSectionHistoryTimes++
				rehearsalScore = 125
				progressProcessor.On("ComputeRehearsalsScore", history).Return(rehearsalScore).Once()
			}

			if tt.songSection.Confidence != tt.request.Confidence {
				songSectionHistoryTimes++
				confidenceScore = 88
				progressProcessor.On("ComputeConfidenceScore", history).Return(confidenceScore).Once()
			}

			if songSectionHistoryTimes > 0 {
				songRepository.On("CreateSongSectionHistory", mock.IsType(new(model.SongSectionHistory))).
					Run(func(args mock.Arguments) {
						newHistory := args.Get(0).(*model.SongSectionHistory)

						assert.NotEmpty(t, newHistory.ID)
						assert.Equal(t, tt.songSection.ID, newHistory.SongSectionID)
						assert.NotEmpty(t, newHistory.Property)
						if newHistory.Property == model.ConfidenceProperty {
							assert.Equal(t, tt.songSection.Confidence, newHistory.From)
							assert.Equal(t, tt.request.Confidence, newHistory.To)
						} else if newHistory.Property == model.RehearsalsProperty {
							assert.Equal(t, tt.songSection.Rehearsals, newHistory.From)
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
						tt.songSection.ID,
						mock.IsType(model.ConfidenceProperty),
					).
					Run(func(args mock.Arguments) {
						property := args.Get(2).(model.SongSectionProperty)
						assert.True(t, property == model.ConfidenceProperty || property == model.RehearsalsProperty)
					}).
					Return(nil, &history).
					Times(songSectionHistoryTimes)

				progressProcessor.On("ComputeProgress", mock.IsType(*tt.songSection)).
					Run(func(args mock.Arguments) {
						newSection := args.Get(0).(model.SongSection)
						assert.Equal(t, tt.songSection.ID, newSection.ID)
					}).
					Return(tt.progress).
					Times(songSectionHistoryTimes)

				songRepository.On("Get", new(model.Song), tt.songSection.SongID).
					Return(nil, tt.song).
					Once()

				if tt.request.BandMemberID != tt.songSection.BandMemberID && tt.request.BandMemberID != nil {
					songRepository.On("IsBandMemberAssociatedWithSong", tt.songSection.SongID, *tt.request.BandMemberID).
						Return(true, nil).
						Once()
				}

				songRepository.On("CountSectionsBySong", mock.IsType(tt.sectionsCount), tt.songSection.SongID).
					Return(nil, tt.sectionsCount).
					Once()

				songRepository.On("Update", mock.IsType(tt.song)).
					Run(func(args mock.Arguments) {
						newSong := args.Get(0).(*model.Song)

						assert.Equal(t, tt.expectedSongConfidence, newSong.Confidence)
						assert.Equal(t, tt.expectedSongRehearsals, newSong.Rehearsals)
						assert.Equal(t, tt.expectedSongProgress, newSong.Progress)

						if tt.songSection.Rehearsals != tt.request.Rehearsals {
							assert.WithinDuration(t, time.Now(), *newSong.LastTimePlayed, 1*time.Minute)
						}
					}).
					Return(nil).
					Once()
			}

			songRepository.On("UpdateSection", mock.IsType(new(model.SongSection))).
				Run(func(args mock.Arguments) {
					newSection := args.Get(0).(*model.SongSection)
					assertUpdatedSongSection(t, tt.request, *newSection, rehearsalScore, confidenceScore, tt.progress)
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
	progress uint64,
) {
	assert.Equal(t, request.Name, section.Name)
	assert.Equal(t, request.Confidence, section.Confidence)
	assert.Equal(t, request.Rehearsals, section.Rehearsals)
	assert.Equal(t, rehearsalsScore, section.RehearsalsScore)
	assert.Equal(t, confidenceScore, section.ConfidenceScore)
	assert.Equal(t, progress, section.Progress)
	assert.Equal(t, request.TypeID, section.SongSectionTypeID)
	assert.Equal(t, request.BandMemberID, section.BandMemberID)
	assert.Equal(t, request.InstrumentID, section.InstrumentID)
}
