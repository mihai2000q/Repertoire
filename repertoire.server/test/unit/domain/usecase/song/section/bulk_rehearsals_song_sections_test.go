package section

import (
	"errors"
	"math"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song/section"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/database/transaction"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/domain/processor"
	"slices"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBulkRehearsalsSongSections_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewBulkRehearsalsSongSections(songRepository, nil, nil)

	request := requests.BulkRehearsalsSongSectionsRequest{
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

func TestBulkRehearsalsSongSections_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewBulkRehearsalsSongSections(songRepository, nil, nil)

	request := requests.BulkRehearsalsSongSectionsRequest{
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

func TestBulkRehearsalsSongSections_WhenSectionsAreNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewBulkRehearsalsSongSections(songRepository, nil, nil)

	request := requests.BulkRehearsalsSongSectionsRequest{
		Sections: []requests.BulkRehearsalsSongSectionRequest{{ID: uuid.New(), Rehearsals: 12}},
		SongID:   uuid.New(),
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

func TestBulkRehearsalsSongSections_WhenNotAllSectionsAreFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewBulkRehearsalsSongSections(songRepository, nil, nil)

	request := requests.BulkRehearsalsSongSectionsRequest{
		Sections: []requests.BulkRehearsalsSongSectionRequest{
			{ID: uuid.New(), Rehearsals: 12},
			{ID: uuid.New(), Rehearsals: 14},
		},
		SongID: uuid.New(),
	}

	song := &model.Song{
		ID: request.SongID,
		Sections: []model.SongSection{
			{ID: request.Sections[0].ID, Order: 0},
			{ID: uuid.New(), Order: 1},
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

func TestBulkRehearsalsSongSections_WhenTransactionExecuteFails_ShouldReturnError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	progressProcessor := new(processor.ProgressProcessorMock)
	_uut := section.NewBulkRehearsalsSongSections(songRepository, transactionManager, progressProcessor)

	repositoryFactory := new(transaction.RepositoryFactoryMock)

	request := requests.BulkRehearsalsSongSectionsRequest{
		Sections: []requests.BulkRehearsalsSongSectionRequest{
			{ID: uuid.New(), Rehearsals: 1},
		},
		SongID: uuid.New(),
	}

	song := &model.Song{
		ID: request.SongID,
		Sections: []model.SongSection{
			{ID: request.Sections[0].ID, Order: 0},
		},
	}

	// given - mocking
	songRepository.On("GetWithSections", new(model.Song), request.SongID).
		Return(nil, song).
		Once()

	internalError := errors.New("internal error")
	transactionManager.On("Execute", mock.Anything).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, errCode.Error, internalError)
	assert.Equal(t, errCode.Code, http.StatusInternalServerError)

	songRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	progressProcessor.AssertExpectations(t)

	repositoryFactory.AssertExpectations(t)
}

func TestBulkRehearsalsSongSections_WhenCreateHistoryFails_ShouldReturnInternalError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	progressProcessor := new(processor.ProgressProcessorMock)
	_uut := section.NewBulkRehearsalsSongSections(songRepository, transactionManager, progressProcessor)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionSongSectionRepository := new(repository.SongSectionRepositoryMock)
	transactionSongRepository := new(repository.SongRepositoryMock)

	request := requests.BulkRehearsalsSongSectionsRequest{
		Sections: []requests.BulkRehearsalsSongSectionRequest{
			{ID: uuid.New(), Rehearsals: 1},
		},
		SongID: uuid.New(),
	}

	song := &model.Song{
		ID: request.SongID,
		Sections: []model.SongSection{
			{ID: request.Sections[0].ID, Order: 0},
		},
	}

	// given - mocking
	songRepository.On("GetWithSections", new(model.Song), request.SongID).
		Return(nil, song).
		Once()

	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(transactionSongSectionRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(transactionSongRepository).Once()

	internalError := errors.New("internal error")
	transactionSongSectionRepository.On("CreateHistory", mock.IsType(new(model.SongSectionHistory))).
		Return(internalError).
		Times(len(request.Sections))

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, errCode.Error, internalError)
	assert.Equal(t, errCode.Code, http.StatusInternalServerError)

	songRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	progressProcessor.AssertExpectations(t)

	repositoryFactory.AssertExpectations(t)
	transactionSongRepository.AssertExpectations(t)
	transactionSongSectionRepository.AssertExpectations(t)
}

func TestBulkRehearsalsSongSections_WhenGetHistoryFails_ShouldReturnInternalError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	progressProcessor := new(processor.ProgressProcessorMock)
	_uut := section.NewBulkRehearsalsSongSections(songRepository, transactionManager, progressProcessor)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionSongSectionRepository := new(repository.SongSectionRepositoryMock)
	transactionSongRepository := new(repository.SongRepositoryMock)

	request := requests.BulkRehearsalsSongSectionsRequest{
		Sections: []requests.BulkRehearsalsSongSectionRequest{
			{ID: uuid.New(), Rehearsals: 1},
		},
		SongID: uuid.New(),
	}

	song := &model.Song{
		ID: request.SongID,
		Sections: []model.SongSection{
			{ID: request.Sections[0].ID, Order: 0},
		},
	}

	// given - mocking
	songRepository.On("GetWithSections", new(model.Song), request.SongID).
		Return(nil, song).
		Once()

	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(transactionSongSectionRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(transactionSongRepository).Once()

	transactionSongSectionRepository.On("CreateHistory", mock.IsType(new(model.SongSectionHistory))).
		Return(nil).
		Times(len(request.Sections))

	internalError := errors.New("internal error")
	transactionSongSectionRepository.
		On(
			"GetHistory",
			new([]model.SongSectionHistory),
			mock.IsType(uuid.UUID{}),
			model.RehearsalsProperty,
		).
		Return(internalError).
		Times(len(request.Sections))

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, errCode.Error, internalError)
	assert.Equal(t, errCode.Code, http.StatusInternalServerError)

	songRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	progressProcessor.AssertExpectations(t)

	repositoryFactory.AssertExpectations(t)
	transactionSongRepository.AssertExpectations(t)
	transactionSongSectionRepository.AssertExpectations(t)
}

func TestBulkRehearsalsSongSections_WhenUpdateFails_ShouldReturnInternalError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	progressProcessor := new(processor.ProgressProcessorMock)
	_uut := section.NewBulkRehearsalsSongSections(songRepository, transactionManager, progressProcessor)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionSongSectionRepository := new(repository.SongSectionRepositoryMock)
	transactionSongRepository := new(repository.SongRepositoryMock)

	request := requests.BulkRehearsalsSongSectionsRequest{
		Sections: []requests.BulkRehearsalsSongSectionRequest{
			{ID: uuid.New(), Rehearsals: 1},
		},
		SongID: uuid.New(),
	}

	song := &model.Song{
		ID: request.SongID,
		Sections: []model.SongSection{
			{ID: request.Sections[0].ID, Order: 0},
		},
	}

	// given - mocking
	songRepository.On("GetWithSections", new(model.Song), request.SongID).
		Return(nil, song).
		Once()

	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(transactionSongSectionRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(transactionSongRepository).Once()

	transactionSongSectionRepository.On("CreateHistory", mock.IsType(new(model.SongSectionHistory))).
		Return(nil).
		Times(len(request.Sections))

	transactionSongSectionRepository.
		On(
			"GetHistory",
			new([]model.SongSectionHistory),
			mock.IsType(uuid.UUID{}),
			model.RehearsalsProperty,
		).
		Return(nil).
		Times(len(request.Sections))

	progressProcessor.On("ComputeRehearsalsScore", mock.IsType([]model.SongSectionHistory{})).
		Return(uint64(0)).
		Times(len(request.Sections))

	progressProcessor.On("ComputeProgress", mock.IsType(model.SongSection{})).
		Return(uint64(0)).
		Times(len(request.Sections))

	internalError := errors.New("internal error")
	transactionSongRepository.On("UpdateWithAssociations", mock.IsType(new(model.Song))).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, errCode.Error, internalError)
	assert.Equal(t, errCode.Code, http.StatusInternalServerError)

	songRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	progressProcessor.AssertExpectations(t)

	repositoryFactory.AssertExpectations(t)
	transactionSongRepository.AssertExpectations(t)
	transactionSongSectionRepository.AssertExpectations(t)
}

func TestBulkRehearsalsSongSections_WhenSongIsNotUpdated_ShouldNotUpdateSong(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	progressProcessor := new(processor.ProgressProcessorMock)
	_uut := section.NewBulkRehearsalsSongSections(songRepository, transactionManager, progressProcessor)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionSongSectionRepository := new(repository.SongSectionRepositoryMock)
	transactionSongRepository := new(repository.SongRepositoryMock)

	request := requests.BulkRehearsalsSongSectionsRequest{
		Sections: []requests.BulkRehearsalsSongSectionRequest{
			{ID: uuid.New(), Rehearsals: 1},
		},
		SongID: uuid.New(),
	}

	song := &model.Song{
		ID: request.SongID,
		Sections: []model.SongSection{
			{ID: request.Sections[0].ID, Order: 0, Rehearsals: 1},
		},
	}

	// given - mocking
	songRepository.On("GetWithSections", new(model.Song), request.SongID).
		Return(nil, song).
		Once()

	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(transactionSongSectionRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(transactionSongRepository).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	progressProcessor.AssertExpectations(t)

	repositoryFactory.AssertExpectations(t)
	transactionSongRepository.AssertExpectations(t)
	transactionSongSectionRepository.AssertExpectations(t)
}

func TestBulkRehearsalsSongSections_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	type RequestSectionsTestData struct {
		Index         uint
		NewRehearsals uint
	}

	tests := []struct {
		name             string
		song             model.Song
		sectionsTestData []RequestSectionsTestData
		expectedSong     model.Song
	}{
		{
			"1 - When there are more sections, but without stats",
			model.Song{
				ID: uuid.New(),
				Sections: []model.SongSection{
					{ID: uuid.New(), Order: 0},
					{ID: uuid.New(), Order: 1},
					{ID: uuid.New(), Order: 2},
					{ID: uuid.New(), Order: 3},
					{ID: uuid.New(), Order: 4},
				},
				Rehearsals: 0,
				Progress:   0,
			},
			[]RequestSectionsTestData{
				{Index: 0, NewRehearsals: 1},
				{Index: 2, NewRehearsals: 1},
			},
			model.Song{
				ID: uuid.New(),
				Sections: []model.SongSection{
					{ID: uuid.New(), Order: 0, Rehearsals: 1, RehearsalsScore: 5, Progress: 1},
					{ID: uuid.New(), Order: 1},
					{ID: uuid.New(), Order: 2, Rehearsals: 1, RehearsalsScore: 2, Progress: 1},
					{ID: uuid.New(), Order: 3},
					{ID: uuid.New(), Order: 4},
				},
				Rehearsals: 0,
				Progress:   0,
			},
		},
		{
			"2 - When there are more sections, but without stats and skip one",
			model.Song{
				ID: uuid.New(),
				Sections: []model.SongSection{
					{ID: uuid.New(), Order: 0, Rehearsals: 1, RehearsalsScore: 1, Progress: 3},
					{ID: uuid.New(), Order: 1},
					{ID: uuid.New(), Order: 2},
					{ID: uuid.New(), Order: 3},
					{ID: uuid.New(), Order: 4},
				},
				Rehearsals: 1,
				Progress:   1,
			},
			[]RequestSectionsTestData{
				{Index: 0, NewRehearsals: 1},
				{Index: 2, NewRehearsals: 1},
			},
			model.Song{
				ID: uuid.New(),
				Sections: []model.SongSection{
					{ID: uuid.New(), Order: 0, Rehearsals: 1, RehearsalsScore: 1, Progress: 3},
					{ID: uuid.New(), Order: 1},
					{ID: uuid.New(), Order: 2, Rehearsals: 1, RehearsalsScore: 1, Progress: 1},
					{ID: uuid.New(), Order: 3},
					{ID: uuid.New(), Order: 4},
				},
				Rehearsals: 1,
				Progress:   1,
			},
		},
		{
			"3 - When there are more sections with stats",
			model.Song{
				ID: uuid.New(),
				Sections: []model.SongSection{
					{ID: uuid.New(), Order: 0, Rehearsals: 12, RehearsalsScore: 5, Progress: 45},
					{ID: uuid.New(), Order: 1, Rehearsals: 5, RehearsalsScore: 2, Progress: 15},
					{ID: uuid.New(), Order: 2, Rehearsals: 25, RehearsalsScore: 45, Progress: 100},
					{ID: uuid.New(), Order: 3, Rehearsals: 6, RehearsalsScore: 1, Progress: 63},
					{ID: uuid.New(), Order: 4, Rehearsals: 19, RehearsalsScore: 25, Progress: 170},
				},
				Rehearsals: 13.4,
				Progress:   78.6,
			},
			[]RequestSectionsTestData{
				{Index: 0, NewRehearsals: 22},
				{Index: 2, NewRehearsals: 40},
			},
			model.Song{
				Sections: []model.SongSection{
					{ID: uuid.New(), Order: 0, Rehearsals: 22, RehearsalsScore: 56, Progress: 102},
					{ID: uuid.New(), Order: 1, Rehearsals: 5, RehearsalsScore: 2, Progress: 15},
					{ID: uuid.New(), Order: 2, Rehearsals: 40, RehearsalsScore: 120, Progress: 201},
					{ID: uuid.New(), Order: 3, Rehearsals: 6, RehearsalsScore: 1, Progress: 63},
					{ID: uuid.New(), Order: 4, Rehearsals: 19, RehearsalsScore: 25, Progress: 170},
				},
				Rehearsals: 18,
				Progress:   110,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			songRepository := new(repository.SongRepositoryMock)
			transactionManager := new(transaction.ManagerMock)
			progressProcessor := new(processor.ProgressProcessorMock)
			_uut := section.NewBulkRehearsalsSongSections(songRepository, transactionManager, progressProcessor)

			repositoryFactory := new(transaction.RepositoryFactoryMock)
			transactionSongSectionRepository := new(repository.SongSectionRepositoryMock)
			transactionSongRepository := new(repository.SongRepositoryMock)

			var requestSections []requests.BulkRehearsalsSongSectionRequest
			for _, s := range tt.sectionsTestData {
				requestSections = append(requestSections, requests.BulkRehearsalsSongSectionRequest{
					ID:         tt.song.Sections[s.Index].ID,
					Rehearsals: s.NewRehearsals,
				})
			}
			request := requests.BulkRehearsalsSongSectionsRequest{
				SongID:   tt.song.ID,
				Sections: requestSections,
			}

			// given - mocking
			songRepository.On("GetWithSections", new(model.Song), request.SongID).
				Return(nil, &tt.song).
				Once()

			transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()
			repositoryFactory.On("NewSongSectionRepository").Return(transactionSongSectionRepository).Once()
			repositoryFactory.On("NewSongRepository").Return(transactionSongRepository).Once()

			for _, s := range request.Sections {
				sectionIndex := slices.IndexFunc(tt.song.Sections, func(sec model.SongSection) bool {
					return sec.ID == s.ID
				})
				if sectionIndex == -1 || s.Rehearsals == tt.song.Sections[sectionIndex].Rehearsals {
					continue
				}

				transactionSongSectionRepository.On("CreateHistory", mock.IsType(new(model.SongSectionHistory))).
					Run(func(args mock.Arguments) {
						newHistory := args.Get(0).(*model.SongSectionHistory)

						assert.NotEmpty(t, newHistory.ID)
						assert.Equal(t, model.RehearsalsProperty, newHistory.Property)
						assert.Equal(t, s.Rehearsals, newHistory.To)
						songSectionIndex := slices.IndexFunc(tt.song.Sections, func(sec model.SongSection) bool {
							return sec.ID == newHistory.SongSectionID
						})
						assert.Equal(t, tt.song.Sections[songSectionIndex].Rehearsals, newHistory.From)
					}).
					Return(nil).
					Once()

				var history []model.SongSectionHistory
				transactionSongSectionRepository.
					On(
						"GetHistory",
						new([]model.SongSectionHistory),
						s.ID,
						model.RehearsalsProperty,
					).
					Return(nil, &history).
					Once()

				progressProcessor.On("ComputeRehearsalsScore", history).
					Return(tt.expectedSong.Sections[sectionIndex].RehearsalsScore).
					Once()

				progressProcessor.On("ComputeProgress", mock.IsType(model.SongSection{})).
					Run(func(args mock.Arguments) {
						ss := args.Get(0).(model.SongSection)
						assert.Equal(t, s.ID, ss.ID)
					}).
					Return(tt.expectedSong.Sections[sectionIndex].Progress).
					Once()
			}

			transactionSongRepository.On("UpdateWithAssociations", mock.IsType(&tt.song)).
				Run(func(args mock.Arguments) {
					newSong := args.Get(0).(*model.Song)

					// song stats updated
					assert.Equal(t, tt.expectedSong.Rehearsals, math.Round(newSong.Rehearsals))
					assert.Equal(t, tt.expectedSong.Progress, math.Round(newSong.Progress))
					assert.WithinDuration(t, time.Now(), *newSong.LastTimePlayed, 1*time.Minute)

					// sections stats updated
					for i, newSection := range newSong.Sections {
						expectedSection := tt.expectedSong.Sections[i]
						assert.Equal(t, expectedSection.Rehearsals, newSection.Rehearsals)
						assert.Equal(t, expectedSection.RehearsalsScore, newSection.RehearsalsScore)
						assert.Equal(t, expectedSection.Progress, newSection.Progress)
					}
				}).
				Return(nil).
				Once()

			// when
			errCode := _uut.Handle(request)

			// then
			assert.Nil(t, errCode)

			songRepository.AssertExpectations(t)
			transactionManager.AssertExpectations(t)
			progressProcessor.AssertExpectations(t)

			repositoryFactory.AssertExpectations(t)
			transactionSongRepository.AssertExpectations(t)
			transactionSongSectionRepository.AssertExpectations(t)
		})
	}
}
