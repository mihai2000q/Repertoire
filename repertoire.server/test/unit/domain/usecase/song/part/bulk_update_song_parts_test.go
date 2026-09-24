package part

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song/part"
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
	"github.com/stretchr/testify/require"
)

func TestBulkUpdateSongParts_WhenTransactionExecuteFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewBulkUpdateSongParts(transactionManager, nil)

	request := requests.BulkUpdateSongPartsRequest{
		SongID: uuid.New(),
		Requests: []requests.BulkUpdateSongPartRequest{
			{ID: uuid.New(), Rehearsals: 5},
		},
	}

	internalError := errors.New("transaction error")
	transactionManager.On("Execute", mock.Anything).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	transactionManager.AssertExpectations(t)
}

func TestBulkUpdateSongParts_WhenGetWithPartsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewBulkUpdateSongParts(transactionManager, nil)

	request := requests.BulkUpdateSongPartsRequest{
		SongID: uuid.New(),
		Requests: []requests.BulkUpdateSongPartRequest{
			{ID: uuid.New(), Rehearsals: 5},
		},
	}

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	internalError := errors.New("get error")
	txSongRepo.On("GetWithParts", new(model.Song), request.SongID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongPartRepo.AssertExpectations(t)
}

func TestBulkUpdateSongParts_WhenSongNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewBulkUpdateSongParts(transactionManager, nil)

	request := requests.BulkUpdateSongPartsRequest{
		SongID: uuid.New(),
		Requests: []requests.BulkUpdateSongPartRequest{
			{ID: uuid.New(), Rehearsals: 5},
		},
	}

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txSongRepo.On("GetWithParts", new(model.Song), request.SongID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())

	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongPartRepo.AssertExpectations(t)
}

func TestBulkUpdateSongParts_WhenPartsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewBulkUpdateSongParts(transactionManager, nil)

	partID := uuid.New()
	request := requests.BulkUpdateSongPartsRequest{
		SongID: uuid.New(),
		Requests: []requests.BulkUpdateSongPartRequest{
			{ID: partID, Rehearsals: 5},
		},
	}

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	mockSong := &model.Song{
		ID: request.SongID,
		Parts: []model.SongPart{
			{ID: uuid.New()}, // different ID
		},
	}
	txSongRepo.On("GetWithParts", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song parts not found", errCode.Error.Error())

	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongPartRepo.AssertExpectations(t)
}

func TestBulkUpdateSongParts_WhenRehearsalsDecrease_ShouldReturnConflictError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewBulkUpdateSongParts(transactionManager, nil)

	partID := uuid.New()
	request := requests.BulkUpdateSongPartsRequest{
		SongID: uuid.New(),
		Requests: []requests.BulkUpdateSongPartRequest{
			{ID: partID, Rehearsals: 5},
		},
	}

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	mockSong := &model.Song{
		ID: request.SongID,
		Parts: []model.SongPart{
			{ID: partID, Rehearsals: 10},
		},
	}
	txSongRepo.On("GetWithParts", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusConflict, errCode.Code)
	assert.Equal(t, "rehearsals can only be increased", errCode.Error.Error())

	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongPartRepo.AssertExpectations(t)
}

func TestBulkUpdateSongParts_WhenCreateRehearsalsHistoryFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	progressProcessor := new(processor.ProgressProcessorMock)
	_uut := part.NewBulkUpdateSongParts(transactionManager, progressProcessor)

	partID := uuid.New()
	request := requests.BulkUpdateSongPartsRequest{
		SongID: uuid.New(),
		Requests: []requests.BulkUpdateSongPartRequest{
			{ID: partID, Rehearsals: 15, Confidence: 30},
		},
	}

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	mockSong := &model.Song{
		ID: request.SongID,
		Parts: []model.SongPart{
			{ID: partID, Rehearsals: 10, Confidence: 30},
		},
	}
	txSongRepo.On("GetWithParts", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	internalError := errors.New("create history error")
	txSongPartRepo.On("CreateHistory", mock.IsType(new(model.SongPartHistory))).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongPartRepo.AssertExpectations(t)
}

func TestBulkUpdateSongParts_WhenGetRehearsalsHistoryFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	progressProcessor := new(processor.ProgressProcessorMock)
	_uut := part.NewBulkUpdateSongParts(transactionManager, progressProcessor)

	partID := uuid.New()
	request := requests.BulkUpdateSongPartsRequest{
		SongID: uuid.New(),
		Requests: []requests.BulkUpdateSongPartRequest{
			{ID: partID, Rehearsals: 15, Confidence: 30},
		},
	}

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	mockSong := &model.Song{
		ID: request.SongID,
		Parts: []model.SongPart{
			{ID: partID, Rehearsals: 10, Confidence: 30},
		},
	}
	txSongRepo.On("GetWithParts", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	txSongPartRepo.On("CreateHistory", mock.IsType(new(model.SongPartHistory))).
		Return(nil).
		Once()

	internalError := errors.New("get history error")
	txSongPartRepo.On("GetHistory", new([]model.SongPartHistory), partID, model.RehearsalsProperty).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongPartRepo.AssertExpectations(t)
}

func TestBulkUpdateSongParts_WhenCreateConfidenceHistoryFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	progressProcessor := new(processor.ProgressProcessorMock)
	_uut := part.NewBulkUpdateSongParts(transactionManager, progressProcessor)

	partID := uuid.New()
	request := requests.BulkUpdateSongPartsRequest{
		SongID: uuid.New(),
		Requests: []requests.BulkUpdateSongPartRequest{
			{ID: partID, Rehearsals: 10, Confidence: 50},
		},
	}

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	mockSong := &model.Song{
		ID: request.SongID,
		Parts: []model.SongPart{
			{ID: partID, Rehearsals: 10, Confidence: 30},
		},
	}
	txSongRepo.On("GetWithParts", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	internalError := errors.New("create history error")
	txSongPartRepo.On("CreateHistory", mock.IsType(new(model.SongPartHistory))).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongPartRepo.AssertExpectations(t)
}

func TestBulkUpdateSongParts_WhenGetConfidenceHistoryFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	progressProcessor := new(processor.ProgressProcessorMock)
	_uut := part.NewBulkUpdateSongParts(transactionManager, progressProcessor)

	partID := uuid.New()
	request := requests.BulkUpdateSongPartsRequest{
		SongID: uuid.New(),
		Requests: []requests.BulkUpdateSongPartRequest{
			{ID: partID, Rehearsals: 10, Confidence: 50},
		},
	}

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	mockSong := &model.Song{
		ID: request.SongID,
		Parts: []model.SongPart{
			{ID: partID, Rehearsals: 10, Confidence: 30},
		},
	}
	txSongRepo.On("GetWithParts", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	txSongPartRepo.On("CreateHistory", mock.IsType(new(model.SongPartHistory))).
		Return(nil).
		Once()

	internalError := errors.New("get history error")
	txSongPartRepo.On("GetHistory", new([]model.SongPartHistory), partID, model.ConfidenceProperty).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongPartRepo.AssertExpectations(t)
}

func TestBulkUpdateSongParts_WhenUpdateAllPartsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	progressProcessor := new(processor.ProgressProcessorMock)
	_uut := part.NewBulkUpdateSongParts(transactionManager, progressProcessor)

	partID := uuid.New()
	request := requests.BulkUpdateSongPartsRequest{
		SongID: uuid.New(),
		Requests: []requests.BulkUpdateSongPartRequest{
			{ID: partID, Rehearsals: 15, Confidence: 30}, // Confidence unchanged
		},
	}

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	mockSong := &model.Song{
		ID: request.SongID,
		Parts: []model.SongPart{
			{ID: partID, Rehearsals: 10, Confidence: 30, Progress: 0},
		},
	}
	txSongRepo.On("GetWithParts", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	txSongPartRepo.On("CreateHistory", mock.IsType(new(model.SongPartHistory))).
		Return(nil).
		Once()
	var history []model.SongPartHistory
	txSongPartRepo.On("GetHistory", new([]model.SongPartHistory), partID, model.RehearsalsProperty).
		Return(nil, &history).
		Once()
	progressProcessor.On("ComputeRehearsalsScore", history).Return(uint64(125)).Once()
	progressProcessor.On("ComputeProgress", mock.IsType(model.SongPart{})).Return(uint64(50)).Once()

	txSongRepo.On("UpdateWithAssociations", mock.IsType(mockSong)).
		Return(nil).
		Once()

	internalError := errors.New("update all error")
	txSongPartRepo.On("UpdateAll", mock.IsType(new([]model.SongPart))).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongPartRepo.AssertExpectations(t)
	progressProcessor.AssertExpectations(t)
}

func TestBulkUpdateSongParts_WhenUpdateWithAssociationsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	progressProcessor := new(processor.ProgressProcessorMock)
	_uut := part.NewBulkUpdateSongParts(transactionManager, progressProcessor)

	partID := uuid.New()
	request := requests.BulkUpdateSongPartsRequest{
		SongID: uuid.New(),
		Requests: []requests.BulkUpdateSongPartRequest{
			{ID: partID, Rehearsals: 15, Confidence: 30},
		},
	}

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	mockSong := &model.Song{
		ID: request.SongID,
		Parts: []model.SongPart{
			{ID: partID, Rehearsals: 10, Confidence: 30, Progress: 0},
		},
	}
	txSongRepo.On("GetWithParts", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	txSongPartRepo.On("CreateHistory", mock.IsType(new(model.SongPartHistory))).
		Return(nil).
		Once()
	var history []model.SongPartHistory
	txSongPartRepo.On("GetHistory", new([]model.SongPartHistory), partID, model.RehearsalsProperty).
		Return(nil, &history).
		Once()
	progressProcessor.On("ComputeRehearsalsScore", history).Return(uint64(125)).Once()
	progressProcessor.On("ComputeProgress", mock.IsType(model.SongPart{})).Return(uint64(50)).Once()

	internalError := errors.New("update song error")
	txSongRepo.On("UpdateWithAssociations", mock.IsType(mockSong)).
		Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongPartRepo.AssertExpectations(t)
	progressProcessor.AssertExpectations(t)
}

func TestBulkUpdateSongParts_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	tests := []struct {
		name                   string
		song                   model.Song
		requests               []requests.BulkUpdateSongPartRequest
		expectedPartProgress   []uint64
		expectedSongConfidence float64
		expectedSongRehearsals float64
		expectedSongProgress   float64
		rehearsalsChanged      bool
	}{
		{
			name: "No changes",
			song: model.Song{
				ID: uuid.New(),
				Parts: []model.SongPart{
					{ID: uuid.New(), Rehearsals: 5, Confidence: 30, Progress: 10},
				},
				Confidence: 30,
				Rehearsals: 5,
				Progress:   10,
			},
			requests: []requests.BulkUpdateSongPartRequest{
				{ID: uuid.New(), Rehearsals: 5, Confidence: 30},
			},
			expectedPartProgress:   nil, // no changes, so no progress computed
			expectedSongConfidence: 30,
			expectedSongRehearsals: 5,
			expectedSongProgress:   10,
			rehearsalsChanged:      false,
		},
		{
			name: "Only rehearse",
			song: model.Song{
				ID: uuid.New(),
				Parts: []model.SongPart{
					{ID: uuid.New(), Rehearsals: 5, Confidence: 30, Progress: 10},
				},
				Confidence: 30,
				Rehearsals: 5,
				Progress:   10,
			},
			requests: []requests.BulkUpdateSongPartRequest{
				{ID: uuid.New(), Rehearsals: 10, Confidence: 30},
			},
			expectedPartProgress:   []uint64{25},
			expectedSongConfidence: 30,
			expectedSongRehearsals: 10,
			expectedSongProgress:   25,
			rehearsalsChanged:      true,
		},
		{
			name: "Only confidence",
			song: model.Song{
				ID: uuid.New(),
				Parts: []model.SongPart{
					{ID: uuid.New(), Rehearsals: 5, Confidence: 30, Progress: 10},
				},
				Confidence: 30,
				Rehearsals: 5,
				Progress:   10,
			},
			requests: []requests.BulkUpdateSongPartRequest{
				{ID: uuid.New(), Rehearsals: 5, Confidence: 50},
			},
			expectedPartProgress:   []uint64{15},
			expectedSongConfidence: 50,
			expectedSongRehearsals: 5,
			expectedSongProgress:   15,
			rehearsalsChanged:      false,
		},
		{
			name: "Both",
			song: model.Song{
				ID: uuid.New(),
				Parts: []model.SongPart{
					{ID: uuid.New(), Rehearsals: 5, Confidence: 30, Progress: 10},
				},
				Confidence: 30,
				Rehearsals: 5,
				Progress:   10,
			},
			requests: []requests.BulkUpdateSongPartRequest{
				{ID: uuid.New(), Rehearsals: 10, Confidence: 50},
			},
			expectedPartProgress:   []uint64{30},
			expectedSongConfidence: 50,
			expectedSongRehearsals: 10,
			expectedSongProgress:   30,
			rehearsalsChanged:      true,
		},
		{
			name: "Multiple parts, different updates",
			song: model.Song{
				ID: uuid.New(),
				Parts: []model.SongPart{
					{ID: uuid.New(), Rehearsals: 2, Confidence: 20, Progress: 5},
					{ID: uuid.New(), Rehearsals: 4, Confidence: 40, Progress: 10},
				},
				Confidence: 30,
				Rehearsals: 3,
				Progress:   7.5,
			},
			requests: []requests.BulkUpdateSongPartRequest{
				{ID: uuid.New(), Rehearsals: 5, Confidence: 30},
				{ID: uuid.New(), Rehearsals: 6, Confidence: 50},
			},
			expectedPartProgress:   []uint64{15, 25},
			expectedSongConfidence: 40,
			expectedSongRehearsals: 5.5,
			expectedSongProgress:   20,
			rehearsalsChanged:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			transactionManager := new(transaction.ManagerMock)
			progressProcessor := new(processor.ProgressProcessorMock)
			_uut := part.NewBulkUpdateSongParts(transactionManager, progressProcessor)

			// Set up request with correct part IDs
			request := requests.BulkUpdateSongPartsRequest{
				SongID:   tt.song.ID,
				Requests: tt.requests,
			}
			for i := range request.Requests {
				request.Requests[i].ID = tt.song.Parts[i].ID
			}

			repositoryFactory := new(transaction.RepositoryFactoryMock)
			txSongRepo := new(repository.SongRepositoryMock)
			txSongPartRepo := new(repository.SongPartRepositoryMock)

			repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
			repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
			transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

			txSongRepo.On("GetWithParts", new(model.Song), request.SongID).
				Return(nil, &tt.song).Once()

			// Determine which parts have changes and set up mocks
			var history []model.SongPartHistory
			historyTimes := 0

			rehearsalsScore := uint64(125)
			confidenceScore := uint(88)

			for i, p := range tt.song.Parts {
				req := tt.requests[i]
				hasRehearsalsChanged := req.Rehearsals != p.Rehearsals
				hasConfidenceChanged := req.Confidence != p.Confidence
				if hasRehearsalsChanged || hasConfidenceChanged {
					if hasRehearsalsChanged {
						historyTimes++
						progressProcessor.On("ComputeRehearsalsScore", history).Return(rehearsalsScore).Once()
					}
					if hasConfidenceChanged {
						historyTimes++
						progressProcessor.On("ComputeConfidenceScore", history).Return(confidenceScore).Once()
					}
					// ComputeProgress should be called once per modified part
					progressProcessor.On("ComputeProgress", mock.IsType(p)).
						Run(func(args mock.Arguments) {
							partArg := args.Get(0).(model.SongPart)
							// We can't easily compare IDs here because the part is passed by value,
							// but we can check that the fields match the expected request.
							assert.Equal(t, tt.requests[i].ID, partArg.ID)
						}).
						Return(tt.expectedPartProgress[i]).
						Once()
				}
			}

			if historyTimes > 0 {
				txSongPartRepo.On("CreateHistory", mock.IsType(new(model.SongPartHistory))).
					Return(nil).
					Times(historyTimes)
				txSongPartRepo.
					On(
						"GetHistory",
						new([]model.SongPartHistory),
						mock.IsType(uuid.UUID{}),
						mock.IsType(model.ConfidenceProperty),
					).
					Return(nil, &history).
					Times(historyTimes)

				// Update song
				txSongRepo.On("UpdateWithAssociations", mock.IsType(&tt.song)).
					Run(func(args mock.Arguments) {
						updatedSong := args.Get(0).(*model.Song)
						assert.Equal(t, tt.expectedSongConfidence, updatedSong.Confidence)
						assert.Equal(t, tt.expectedSongRehearsals, updatedSong.Rehearsals)
						assert.Equal(t, tt.expectedSongProgress, updatedSong.Progress)
						if tt.rehearsalsChanged {
							assert.WithinDuration(t, time.Now(), *updatedSong.LastTimePlayed, 1*time.Minute)
						}
					}).
					Return(nil).Once()

				txSongPartRepo.On("UpdateAll", mock.IsType(&[]model.SongPart{})).
					Run(func(args mock.Arguments) {
						partsSlice := args.Get(0).(*[]model.SongPart)
						parts := *partsSlice

						assert.Len(t, parts, len(request.Requests))

						for i, expectedPart := range tt.song.Parts {
							idx := slices.IndexFunc(request.Requests, func(r requests.BulkUpdateSongPartRequest) bool {
								return r.ID == expectedPart.ID
							})
							if idx == -1 {
								continue
							}

							req := tt.requests[idx]

							// The part should have the new values
							assert.Equal(t, req.Rehearsals, parts[i].Rehearsals)
							assert.Equal(t, req.Confidence, parts[i].Confidence)
							assert.Equal(t, tt.expectedPartProgress[idx], parts[i].Progress)

							if req.Rehearsals != expectedPart.Rehearsals {
								assert.Equal(t, rehearsalsScore, parts[i].RehearsalsScore)
							}
							if req.Confidence != expectedPart.Confidence {
								assert.Equal(t, confidenceScore, parts[i].ConfidenceScore)
							}
						}
					}).
					Return(nil).Once()
			}

			// when
			errCode := _uut.Handle(request)

			// then
			assert.Nil(t, errCode)

			transactionManager.AssertExpectations(t)
			repositoryFactory.AssertExpectations(t)
			txSongRepo.AssertExpectations(t)
			txSongPartRepo.AssertExpectations(t)
			progressProcessor.AssertExpectations(t)
		})
	}
}
