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
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateSongPart_WhenGetPartFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songPartRepository := new(repository.SongPartRepositoryMock)
	_uut := part.NewUpdateSongPart(songPartRepository, nil, nil, nil)

	request := requests.UpdateSongPartRequest{
		ID:   uuid.New(),
		Name: "Some Part",
	}

	internalError := errors.New("get error")
	songPartRepository.On("Get", new(model.SongPart), request.ID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songPartRepository.AssertExpectations(t)
}

func TestUpdateSongPart_WhenPartNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songPartRepository := new(repository.SongPartRepositoryMock)
	_uut := part.NewUpdateSongPart(songPartRepository, nil, nil, nil)

	request := requests.UpdateSongPartRequest{
		ID:   uuid.New(),
		Name: "Some Part",
	}

	songPartRepository.On("Get", new(model.SongPart), request.ID).
		Return(nil). // not found
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song part not found", errCode.Error.Error())

	songPartRepository.AssertExpectations(t)
}

func TestUpdateSongPart_WhenRehearsalsDecreasing_ShouldReturnConflictError(t *testing.T) {
	// given
	songPartRepository := new(repository.SongPartRepositoryMock)
	_uut := part.NewUpdateSongPart(songPartRepository, nil, nil, nil)

	request := requests.UpdateSongPartRequest{
		ID:         uuid.New(),
		Name:       "Some Part",
		Rehearsals: 10,
	}

	mockPart := &model.SongPart{
		ID:         request.ID,
		Rehearsals: 20, // higher than request
	}
	songPartRepository.On("Get", new(model.SongPart), request.ID).
		Return(nil, mockPart).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusConflict, errCode.Code)
	assert.Equal(t, "rehearsals can only be increased", errCode.Error.Error())

	songPartRepository.AssertExpectations(t)
}

func TestUpdateSongPart_WhenIsBandMemberAssociatedWithSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songPartRepository := new(repository.SongPartRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := part.NewUpdateSongPart(songPartRepository, songRepository, nil, nil)

	request := requests.UpdateSongPartRequest{
		ID:           uuid.New(),
		Name:         "Some Part",
		BandMemberID: &[]uuid.UUID{uuid.New()}[0],
	}

	mockPart := &model.SongPart{
		ID:     request.ID,
		Name:   "Old",
		SongID: uuid.New(),
	}
	songPartRepository.On("Get", new(model.SongPart), request.ID).
		Return(nil, mockPart).
		Once()

	internalError := errors.New("association check error")
	songRepository.On("IsBandMemberAssociatedWithSong", mockPart.SongID, *request.BandMemberID).
		Return(false, internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songPartRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestUpdateSongPart_WhenBandMemberNotAssociated_ShouldReturnConflictError(t *testing.T) {
	// given
	songPartRepository := new(repository.SongPartRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := part.NewUpdateSongPart(songPartRepository, songRepository, nil, nil)

	request := requests.UpdateSongPartRequest{
		ID:           uuid.New(),
		Name:         "Some Part",
		BandMemberID: &[]uuid.UUID{uuid.New()}[0],
	}

	mockPart := &model.SongPart{
		ID:     request.ID,
		Name:   "Old",
		SongID: uuid.New(),
	}
	songPartRepository.On("Get", new(model.SongPart), request.ID).
		Return(nil, mockPart).
		Once()

	songRepository.On("IsBandMemberAssociatedWithSong", mockPart.SongID, *request.BandMemberID).
		Return(false, nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusConflict, errCode.Code)
	assert.Equal(t, "band member is not part of the artist associated with this song", errCode.Error.Error())

	songPartRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestUpdateSongPart_WhenTransactionExecuteFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songPartRepository := new(repository.SongPartRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewUpdateSongPart(songPartRepository, nil, nil, transactionManager)

	request := requests.UpdateSongPartRequest{
		ID:   uuid.New(),
		Name: "Some Part",
	}

	mockPart := &model.SongPart{
		ID:         request.ID,
		Name:       "Old",
		SongID:     uuid.New(),
		Confidence: 30,
	}
	songPartRepository.On("Get", new(model.SongPart), request.ID).
		Return(nil, mockPart).
		Once()

	internalError := errors.New("transaction error")
	transactionManager.On("Execute", mock.Anything).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songPartRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
}

func TestUpdateSongPart_WhenCreateHistoryFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songPartRepository := new(repository.SongPartRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewUpdateSongPart(songPartRepository, nil, nil, transactionManager)

	request := requests.UpdateSongPartRequest{
		ID:         uuid.New(),
		Name:       "Some Part",
		Confidence: 50,
	}

	// given - mocking
	mockPart := &model.SongPart{
		ID:         request.ID,
		Name:       "Old",
		SongID:     uuid.New(),
		Confidence: 30,
	}
	songPartRepository.On("Get", new(model.SongPart), request.ID).
		Return(nil, mockPart).
		Once()

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

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

	songPartRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongPartRepo.AssertExpectations(t)
}

func TestUpdateSongPart_WhenGetHistoryFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songPartRepository := new(repository.SongPartRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewUpdateSongPart(songPartRepository, nil, nil, transactionManager)

	request := requests.UpdateSongPartRequest{
		ID:         uuid.New(),
		Name:       "Some Part",
		Confidence: 50,
	}

	// given - mocking
	mockPart := &model.SongPart{
		ID:         request.ID,
		Name:       "Old",
		SongID:     uuid.New(),
		Confidence: 30,
	}
	songPartRepository.On("Get", new(model.SongPart), request.ID).
		Return(nil, mockPart).
		Once()

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txSongPartRepo.On("CreateHistory", mock.IsType(new(model.SongPartHistory))).
		Return(nil).
		Once()

	internalError := errors.New("get history error")
	txSongPartRepo.On("GetHistory", mock.IsType(new([]model.SongPartHistory)), mockPart.ID, model.ConfidenceProperty).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songPartRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongPartRepo.AssertExpectations(t)
}

func TestUpdateSongPart_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songPartRepository := new(repository.SongPartRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	progressProcessor := new(processor.ProgressProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewUpdateSongPart(songPartRepository, songRepository, progressProcessor, transactionManager)

	request := requests.UpdateSongPartRequest{
		ID:         uuid.New(),
		Name:       "Some Part",
		Confidence: 50,
	}

	// given - mocking
	mockPart := &model.SongPart{
		ID:         request.ID,
		Name:       "Old",
		SongID:     uuid.New(),
		Confidence: 30,
	}
	songPartRepository.On("Get", new(model.SongPart), request.ID).
		Return(nil, mockPart).
		Once()

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	// updateConfidence mocks
	txSongPartRepo.On("CreateHistory", mock.IsType(new(model.SongPartHistory))).
		Return(nil).
		Once()
	var history []model.SongPartHistory
	txSongPartRepo.On("GetHistory", mock.IsType(new([]model.SongPartHistory)), mockPart.ID, model.ConfidenceProperty).
		Return(nil, &history).
		Once()
	progressProcessor.On("ComputeConfidenceScore", history).Return(uint(88)).Once()

	progressProcessor.On("ComputeProgress", mock.IsType(*mockPart)).Return(uint64(8)).Once()

	internalError := errors.New("get song error")
	txSongRepo.On("Get", new(model.Song), mockPart.SongID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songPartRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
	progressProcessor.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongPartRepo.AssertExpectations(t)
}

func TestUpdateSongPart_WhenCountAllBySongFailsInsideUpdateSongStats_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songPartRepository := new(repository.SongPartRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	progressProcessor := new(processor.ProgressProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewUpdateSongPart(songPartRepository, songRepository, progressProcessor, transactionManager)

	request := requests.UpdateSongPartRequest{
		ID:         uuid.New(),
		Name:       "Some Part",
		Confidence: 50,
	}

	// given - mocking
	mockPart := &model.SongPart{
		ID:         request.ID,
		Name:       "Old",
		SongID:     uuid.New(),
		Confidence: 30,
	}
	songPartRepository.On("Get", new(model.SongPart), request.ID).
		Return(nil, mockPart).
		Once()

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	// updateConfidence mocks
	txSongPartRepo.On("CreateHistory", mock.IsType(new(model.SongPartHistory))).
		Return(nil).
		Once()
	var history []model.SongPartHistory
	txSongPartRepo.On("GetHistory", mock.IsType(new([]model.SongPartHistory)), mockPart.ID, model.ConfidenceProperty).
		Return(nil, &history).
		Once()
	progressProcessor.On("ComputeConfidenceScore", history).Return(uint(88)).Once()

	progressProcessor.On("ComputeProgress", mock.IsType(*mockPart)).Return(uint64(8)).Once()

	mockSong := &model.Song{ID: mockPart.SongID}
	txSongRepo.On("Get", new(model.Song), mockPart.SongID).
		Return(nil, mockSong).
		Once()

	internalError := errors.New("count error")
	txSongPartRepo.On("CountAllBySong", new(int64), mockPart.SongID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songPartRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
	progressProcessor.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongPartRepo.AssertExpectations(t)
}

func TestUpdateSongPart_WhenUpdateSongFailsInsideUpdateSongStats_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songPartRepository := new(repository.SongPartRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	progressProcessor := new(processor.ProgressProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewUpdateSongPart(songPartRepository, songRepository, progressProcessor, transactionManager)

	request := requests.UpdateSongPartRequest{
		ID:         uuid.New(),
		Name:       "Some Part",
		Confidence: 50,
	}

	// given - mocking
	mockPart := &model.SongPart{
		ID:         request.ID,
		Name:       "Old",
		SongID:     uuid.New(),
		Confidence: 30,
	}
	songPartRepository.On("Get", new(model.SongPart), request.ID).
		Return(nil, mockPart).
		Once()

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	// updateConfidence mocks
	txSongPartRepo.On("CreateHistory", mock.IsType(new(model.SongPartHistory))).
		Return(nil).
		Once()
	var history []model.SongPartHistory
	txSongPartRepo.On("GetHistory", mock.IsType(new([]model.SongPartHistory)), mockPart.ID, model.ConfidenceProperty).
		Return(nil, &history).
		Once()
	progressProcessor.On("ComputeConfidenceScore", history).Return(uint(88)).Once()

	progressProcessor.On("ComputeProgress", mock.IsType(*mockPart)).Return(uint64(8)).Once()

	// updateSongStats mocks
	mockSong := &model.Song{
		ID:         mockPart.SongID,
		Confidence: 45.5,
		Progress:   13.5,
	}
	txSongRepo.On("Get", new(model.Song), mockPart.SongID).
		Return(nil, mockSong).
		Once()

	partsCount := int64(4)
	txSongPartRepo.On("CountAllBySong", new(int64), mockPart.SongID).
		Return(nil, &partsCount).
		Once()

	internalError := errors.New("update song error")
	txSongRepo.On("Update", mock.IsType(mockSong)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songPartRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
	progressProcessor.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongPartRepo.AssertExpectations(t)
}

func TestUpdateSongPart_WhenUpdatePartFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songPartRepository := new(repository.SongPartRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	progressProcessor := new(processor.ProgressProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewUpdateSongPart(songPartRepository, songRepository, progressProcessor, transactionManager)

	request := requests.UpdateSongPartRequest{
		ID:   uuid.New(),
		Name: "New Name",
	}

	// given - mocking
	mockPart := &model.SongPart{
		ID:     request.ID,
		Name:   "Old",
		SongID: uuid.New(),
	}
	songPartRepository.On("Get", new(model.SongPart), request.ID).
		Return(nil, mockPart).
		Once()

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	internalError := errors.New("update part error")
	txSongPartRepo.On("Update", mock.IsType(new(model.SongPart))).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songPartRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
	progressProcessor.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongPartRepo.AssertExpectations(t)
}

func TestUpdateSongPart_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	id := uuid.New()
	songID := uuid.New()

	tests := []struct {
		name                   string
		part                   *model.SongPart
		request                requests.UpdateSongPartRequest
		song                   *model.Song
		partsCount             int64
		expectedSongConfidence float64
		expectedSongRehearsals float64
		expectedSongProgress   float64
		progress               uint64
		confidenceScore        uint
		rehearsalsScore        uint64
	}{
		{
			name: "Only name change",
			part: &model.SongPart{
				ID:     id,
				Name:   "Old Name",
				SongID: songID,
			},
			request: requests.UpdateSongPartRequest{
				ID:   id,
				Name: "New Name",
			},
			song:                   nil,
			partsCount:             0,
			expectedSongConfidence: 0,
			expectedSongRehearsals: 0,
			expectedSongProgress:   0,
			progress:               0,
			confidenceScore:        0,
			rehearsalsScore:        0,
		},
		{
			name: "Only confidence change",
			part: &model.SongPart{
				ID:         id,
				Name:       "Old",
				SongID:     songID,
				Confidence: 30,
				Progress:   7,
			},
			request: requests.UpdateSongPartRequest{
				ID:         id,
				Name:       "New",
				Confidence: 50,
			},
			song: &model.Song{
				ID:         songID,
				Confidence: 45.5,
				Progress:   13.5,
			},
			partsCount:             4,
			expectedSongConfidence: 50.5,
			expectedSongRehearsals: 0,
			expectedSongProgress:   13.75,
			progress:               8,
			confidenceScore:        88,
			rehearsalsScore:        0,
		},
		{
			name: "Only rehearsals change",
			part: &model.SongPart{
				ID:         id,
				Name:       "Old",
				SongID:     songID,
				Rehearsals: 35,
				Progress:   7,
			},
			request: requests.UpdateSongPartRequest{
				ID:         id,
				Name:       "New",
				Rehearsals: 40,
			},
			song: &model.Song{
				ID:         songID,
				Rehearsals: 27.5,
				Progress:   13.5,
			},
			partsCount:             4,
			expectedSongConfidence: 0,
			expectedSongRehearsals: 28.75,
			expectedSongProgress:   13.75,
			progress:               8,
			confidenceScore:        0,
			rehearsalsScore:        125,
		},
		{
			name: "Both confidence and rehearsals change",
			part: &model.SongPart{
				ID:         id,
				Name:       "Old",
				SongID:     songID,
				Confidence: 30,
				Rehearsals: 35,
				Progress:   7,
			},
			request: requests.UpdateSongPartRequest{
				ID:         id,
				Name:       "New",
				Confidence: 50,
				Rehearsals: 40,
			},
			song: &model.Song{
				ID:         songID,
				Confidence: 45.5,
				Rehearsals: 27.5,
				Progress:   13.5,
			},
			partsCount:             4,
			expectedSongConfidence: 50.5,
			expectedSongRehearsals: 28.75,
			expectedSongProgress:   13.75,
			progress:               8,
			confidenceScore:        88,
			rehearsalsScore:        125,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			songPartRepository := new(repository.SongPartRepositoryMock)
			songRepository := new(repository.SongRepositoryMock)
			progressProcessor := new(processor.ProgressProcessorMock)
			transactionManager := new(transaction.ManagerMock)
			_uut := part.NewUpdateSongPart(songPartRepository, songRepository, progressProcessor, transactionManager)

			// given - mocking
			songPartRepository.On("Get", new(model.SongPart), tt.request.ID).
				Return(nil, tt.part).
				Once()

			repositoryFactory := new(transaction.RepositoryFactoryMock)
			txSongRepo := new(repository.SongRepositoryMock)
			txSongPartRepo := new(repository.SongPartRepositoryMock)

			repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
			repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
			transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

			hasRehearsalsChanged := tt.part.Rehearsals != tt.request.Rehearsals
			hasConfidenceChanged := tt.part.Confidence != tt.request.Confidence

			if hasRehearsalsChanged || hasConfidenceChanged {
				var history []model.SongPartHistory
				historyTimes := 0
				if hasConfidenceChanged {
					historyTimes++
					progressProcessor.On("ComputeConfidenceScore", history).Return(tt.confidenceScore).Once()
				}
				if hasRehearsalsChanged {
					historyTimes++
					progressProcessor.On("ComputeRehearsalsScore", history).Return(tt.rehearsalsScore).Once()
				}

				txSongPartRepo.On("CreateHistory", mock.IsType(new(model.SongPartHistory))).
					Run(func(args mock.Arguments) {
						h := args.Get(0).(*model.SongPartHistory)
						assert.NotEmpty(t, h.ID)
						assert.Equal(t, tt.part.ID, h.PartID)
						switch h.Property {
						case model.ConfidenceProperty:
							assert.Equal(t, tt.part.Confidence, h.From)
							assert.Equal(t, tt.request.Confidence, h.To)
						case model.RehearsalsProperty:
							assert.Equal(t, tt.part.Rehearsals, h.From)
							assert.Equal(t, tt.request.Rehearsals, h.To)
						default:
							assert.Fail(t, "unexpected property")
						}
					}).
					Return(nil).
					Times(historyTimes)

				txSongPartRepo.
					On(
						"GetHistory",
						new([]model.SongPartHistory),
						tt.part.ID,
						mock.IsType(model.ConfidenceProperty),
					).
					Return(nil, &history).
					Times(historyTimes)

				progressProcessor.On("ComputeProgress", mock.IsType(*tt.part)).
					Return(tt.progress).
					Once()

				// updateSongStats
				txSongRepo.On("Get", new(model.Song), tt.part.SongID).
					Return(nil, tt.song).
					Once()
				txSongPartRepo.On("CountAllBySong", new(int64), tt.part.SongID).
					Return(nil, &tt.partsCount).
					Once()

				txSongRepo.On("Update", mock.IsType(tt.song)).
					Run(func(args mock.Arguments) {
						newSong := args.Get(0).(*model.Song)

						assert.Equal(t, tt.expectedSongConfidence, newSong.Confidence)
						assert.Equal(t, tt.expectedSongRehearsals, newSong.Rehearsals)
						assert.Equal(t, tt.expectedSongProgress, newSong.Progress)

						if hasRehearsalsChanged {
							assert.WithinDuration(t, time.Now(), *newSong.LastTimePlayed, 1*time.Minute)
						}
					}).
					Return(nil).
					Once()
			}

			txSongPartRepo.On("Update", mock.IsType(new(model.SongPart))).
				Run(func(args mock.Arguments) {
					newPart := args.Get(0).(*model.SongPart)
					assertUpdatedSongPart(t, tt.request, newPart, tt.confidenceScore, tt.rehearsalsScore, tt.progress)
				}).
				Return(nil).
				Once()

			// when
			errCode := _uut.Handle(tt.request)

			// then
			assert.Nil(t, errCode)

			songPartRepository.AssertExpectations(t)
			songRepository.AssertExpectations(t)
			progressProcessor.AssertExpectations(t)
			transactionManager.AssertExpectations(t)
			repositoryFactory.AssertExpectations(t)
			txSongRepo.AssertExpectations(t)
			txSongPartRepo.AssertExpectations(t)
		})
	}
}

func assertUpdatedSongPart(
	t *testing.T,
	request requests.UpdateSongPartRequest,
	part *model.SongPart,
	confidenceScore uint,
	rehearsalsScore uint64,
	progress uint64,
) {
	assert.Equal(t, request.Name, part.Name)
	assert.Equal(t, request.Confidence, part.Confidence)
	assert.Equal(t, request.Rehearsals, part.Rehearsals)
	assert.Equal(t, request.BandMemberID, part.BandMemberID)
	assert.Equal(t, request.InstrumentID, part.InstrumentID)
	assert.Equal(t, confidenceScore, part.ConfidenceScore)
	assert.Equal(t, rehearsalsScore, part.RehearsalsScore)
	assert.Equal(t, progress, part.Progress)
}
