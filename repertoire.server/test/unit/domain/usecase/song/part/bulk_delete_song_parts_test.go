package part

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song/part"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/database/transaction"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/domain/processor"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBulkDeleteSongParts_WhenTransactionExecuteFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewBulkDeleteSongParts(transactionManager, nil)

	request := requests.BulkDeleteSongPartsRequest{
		SongID: uuid.New(),
		IDs:    []uuid.UUID{uuid.New()},
	}

	internalError := errors.New("transaction error")
	transactionManager.On("Execute", mock.Anything).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	transactionManager.AssertExpectations(t)
}

func TestBulkDeleteSongParts_WhenSongProcessorUpdateFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	songProcessor := new(processor.SongProcessorMock)
	_uut := part.NewBulkDeleteSongParts(transactionManager, songProcessor)

	request := requests.BulkDeleteSongPartsRequest{
		SongID: uuid.New(),
		IDs:    []uuid.UUID{uuid.New()},
	}

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepository := new(repository.SongRepositoryMock)
	txSongSectionRepository := new(repository.SongSectionRepositoryMock)
	txSongPartRepository := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepository).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepository).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	internalError := httperror.InternalServerError(errors.New("internal error"))
	songProcessor.On("UpdateSongAfterPartsDeletion", txSongRepository, request.SongID, request.IDs).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, internalError, errCode)

	repositoryFactory.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	txSongRepository.AssertExpectations(t)
	txSongSectionRepository.AssertExpectations(t)
	txSongPartRepository.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
}

func TestBulkDeleteSongParts_WhenGetSectionsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	songProcessor := new(processor.SongProcessorMock)
	_uut := part.NewBulkDeleteSongParts(transactionManager, songProcessor)

	request := requests.BulkDeleteSongPartsRequest{
		SongID: uuid.New(),
		IDs:    []uuid.UUID{uuid.New()},
	}

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepository := new(repository.SongRepositoryMock)
	txSongSectionRepository := new(repository.SongSectionRepositoryMock)
	txSongPartRepository := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepository).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepository).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	// Processor succeeds
	songProcessor.On("UpdateSongAfterPartsDeletion", txSongRepository, request.SongID, request.IDs).
		Return(nil).
		Once()

	// GetSections fails
	internalError := errors.New("get sections error")
	txSongSectionRepository.On("GetAllByPartIDsWithSectionParts", new([]model.SongSection), request.IDs).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	repositoryFactory.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	txSongRepository.AssertExpectations(t)
	txSongSectionRepository.AssertExpectations(t)
	txSongPartRepository.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
}

func TestBulkDeleteSongParts_WhenUpdateSectionsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	songProcessor := new(processor.SongProcessorMock)
	_uut := part.NewBulkDeleteSongParts(transactionManager, songProcessor)

	sectionID := uuid.New()
	id := uuid.New()
	remainingPartID := uuid.New()
	request := requests.BulkDeleteSongPartsRequest{
		SongID: uuid.New(),
		IDs:    []uuid.UUID{id},
	}

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepository := new(repository.SongRepositoryMock)
	txSongSectionRepository := new(repository.SongSectionRepositoryMock)
	txSongPartRepository := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepository).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepository).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	// Processor succeeds
	songProcessor.On("UpdateSongAfterPartsDeletion", txSongRepository, request.SongID, request.IDs).
		Return(nil).Once()

	// Get sections
	mockSection := &model.SongSection{
		ID: sectionID,
		SectionParts: []model.SongSectionPart{
			{PartID: id, SectionID: sectionID, Order: 0},
			{PartID: remainingPartID, SectionID: sectionID, Order: 1},
		},
	}
	txSongSectionRepository.On("GetAllByPartIDsWithSectionParts", new([]model.SongSection), request.IDs).
		Return(nil, &[]model.SongSection{*mockSection}).
		Once()

	// Update sections fails
	internalError := errors.New("update sections error")
	txSongSectionRepository.On("UpdateAllSectionParts", mock.IsType(new([]model.SongSectionPart))).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	repositoryFactory.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	txSongRepository.AssertExpectations(t)
	txSongSectionRepository.AssertExpectations(t)
	txSongPartRepository.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
}

func TestBulkDeleteSongParts_WhenDeletePartsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	songProcessor := new(processor.SongProcessorMock)
	_uut := part.NewBulkDeleteSongParts(transactionManager, songProcessor)

	request := requests.BulkDeleteSongPartsRequest{
		SongID: uuid.New(),
		IDs:    []uuid.UUID{uuid.New()},
	}

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepository := new(repository.SongRepositoryMock)
	txSongSectionRepository := new(repository.SongSectionRepositoryMock)
	txSongPartRepository := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepository).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepository).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	// Processor succeeds
	songProcessor.On("UpdateSongAfterPartsDeletion", txSongRepository, request.SongID, request.IDs).
		Return(nil).
		Once()

	// No sections
	txSongSectionRepository.On("GetAllByPartIDsWithSectionParts", new([]model.SongSection), request.IDs).
		Return(nil, &[]model.SongSection{}).
		Once()

	// Delete parts fails
	internalError := errors.New("delete error")
	txSongPartRepository.On("Delete", request.IDs).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	repositoryFactory.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	txSongRepository.AssertExpectations(t)
	txSongSectionRepository.AssertExpectations(t)
	txSongPartRepository.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
}

func TestBulkDeleteSongParts_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	tests := []struct {
		name                   string
		song                   model.Song
		partIndexes            []int
		sectionParts           map[int][]model.SongSectionPart
		expectedSongConfidence float64
		expectedSongRehearsals float64
		expectedSongProgress   float64
	}{
		{
			name: "Delete single part, no sections",
			song: model.Song{
				ID: uuid.New(),
				Parts: []model.SongPart{
					{ID: uuid.New(), SongOrder: 0, Confidence: 55, Rehearsals: 12, Progress: 45},
					{ID: uuid.New(), SongOrder: 1, Confidence: 23, Rehearsals: 5, Progress: 15},
				},
				Confidence: 39,
				Rehearsals: 8.5,
				Progress:   30,
			},
			partIndexes:            []int{0},
			expectedSongConfidence: 23,
			expectedSongRehearsals: 5,
			expectedSongProgress:   15,
		},
		{
			name: "Delete multiple parts, no sections",
			song: model.Song{
				ID: uuid.New(),
				Parts: []model.SongPart{
					{ID: uuid.New(), SongOrder: 0, Confidence: 10, Rehearsals: 2, Progress: 5},
					{ID: uuid.New(), SongOrder: 1, Confidence: 20, Rehearsals: 4, Progress: 10},
					{ID: uuid.New(), SongOrder: 2, Confidence: 30, Rehearsals: 6, Progress: 15},
					{ID: uuid.New(), SongOrder: 3, Confidence: 40, Rehearsals: 8, Progress: 20},
					{ID: uuid.New(), SongOrder: 4, Confidence: 50, Rehearsals: 10, Progress: 25},
				},
				Confidence: 30,
				Rehearsals: 6,
				Progress:   15,
			},
			partIndexes:            []int{1, 3},
			expectedSongConfidence: 30,
			expectedSongRehearsals: 6,
			expectedSongProgress:   15,
		},
		{
			name: "Delete all parts",
			song: model.Song{
				ID: uuid.New(),
				Parts: []model.SongPart{
					{ID: uuid.New(), SongOrder: 0, Confidence: 55, Rehearsals: 12, Progress: 45},
				},
				Confidence: 55,
				Rehearsals: 12,
				Progress:   45,
			},
			partIndexes:            []int{0},
			expectedSongConfidence: 0,
			expectedSongRehearsals: 0,
			expectedSongProgress:   0,
		},
		{
			name: "Delete part from one section",
			song: model.Song{
				ID: uuid.New(),
				Parts: []model.SongPart{
					{ID: uuid.New(), SongOrder: 0, Confidence: 30, Rehearsals: 8, Progress: 20},
					{ID: uuid.New(), SongOrder: 1, Confidence: 60, Rehearsals: 15, Progress: 80},
				},
				Confidence: 45,
				Rehearsals: 11.5,
				Progress:   50,
			},
			partIndexes: []int{1},
			sectionParts: map[int][]model.SongSectionPart{
				1: {
					{PartID: uuid.New(), SectionID: uuid.New(), Order: 1},
				},
			},
			expectedSongConfidence: 30,
			expectedSongRehearsals: 8,
			expectedSongProgress:   20,
		},
		{
			name: "Delete parts from multiple sections",
			song: model.Song{
				ID: uuid.New(),
				Parts: []model.SongPart{
					{ID: uuid.New(), SongOrder: 0, Confidence: 30, Rehearsals: 8, Progress: 20},
					{ID: uuid.New(), SongOrder: 1, Confidence: 60, Rehearsals: 15, Progress: 80},
					{ID: uuid.New(), SongOrder: 2, Confidence: 45, Rehearsals: 10, Progress: 50},
				},
				Confidence: 45,
				Rehearsals: 11,
				Progress:   50,
			},
			partIndexes: []int{0, 2},
			sectionParts: map[int][]model.SongSectionPart{
				0: {
					{PartID: uuid.New(), SectionID: uuid.New(), Order: 0},
					{PartID: uuid.New(), SectionID: uuid.New(), Order: 2},
				},
				2: {
					{PartID: uuid.New(), SectionID: uuid.New(), Order: 1},
				},
			},
			expectedSongConfidence: 60,
			expectedSongRehearsals: 15,
			expectedSongProgress:   80,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			transactionManager := new(transaction.ManagerMock)
			songProcessor := new(processor.SongProcessorMock)
			_uut := part.NewBulkDeleteSongParts(transactionManager, songProcessor)

			request := requests.BulkDeleteSongPartsRequest{
				SongID: tt.song.ID,
				IDs:    []uuid.UUID{},
			}
			for _, idx := range tt.partIndexes {
				request.IDs = append(request.IDs, tt.song.Parts[idx].ID)
			}

			repositoryFactory := new(transaction.RepositoryFactoryMock)
			txSongRepository := new(repository.SongRepositoryMock)
			txSongSectionRepository := new(repository.SongSectionRepositoryMock)
			txSongPartRepository := new(repository.SongPartRepositoryMock)

			repositoryFactory.On("NewSongRepository").Return(txSongRepository).Once()
			repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepository).Once()
			repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepository).Once()
			transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

			// Processor call
			songProcessor.On("UpdateSongAfterPartsDeletion", txSongRepository, request.SongID, request.IDs).
				Return(nil).
				Once()

			// Reorder sections
			sectionMap := make(map[uuid.UUID][]model.SongSectionPart)
			for _, p := range tt.song.Parts {
				for _, sp := range p.SectionParts {
					sectionMap[sp.SectionID] = append(sectionMap[sp.SectionID], sp)
				}
			}

			if len(sectionMap) > 0 {
				var mockSections []model.SongSection
				for secID, sps := range sectionMap {
					mockSections = append(mockSections, model.SongSection{
						ID:           secID,
						SectionParts: sps,
					})
				}
				txSongSectionRepository.
					On(
						"GetAllByPartIDsWithSectionParts",
						new([]model.SongSection),
						request.IDs,
					).
					Return(nil, &mockSections).
					Once()

				// Compute expected updates for UpdateAllSectionParts
				var expectedSectionParts []model.SongSectionPart
				for _, sec := range mockSections {
					var shift uint
					for _, sp := range sec.SectionParts {
						if slices.Contains(request.IDs, sp.PartID) {
							shift++
							continue
						}
						sp.Order -= shift
						expectedSectionParts = append(expectedSectionParts, sp)
					}
				}
				txSongSectionRepository.On("UpdateAllSectionParts", mock.IsType(&expectedSectionParts)).
					Run(func(args mock.Arguments) {
						sectionParts := args.Get(0).(*[]model.SongSectionPart)
						assert.Equal(t, expectedSectionParts, *sectionParts)
					}).
					Return(nil).
					Once()
			} else {
				txSongSectionRepository.On("GetAllByPartIDsWithSectionParts", new([]model.SongSection), request.IDs).
					Return(nil, &[]model.SongSection{}).
					Once()
			}

			// Delete parts
			txSongPartRepository.On("Delete", request.IDs).
				Return(nil).
				Once()

			// when
			errCode := _uut.Handle(request)

			// then
			assert.Nil(t, errCode)

			repositoryFactory.AssertExpectations(t)
			transactionManager.AssertExpectations(t)
			txSongRepository.AssertExpectations(t)
			txSongSectionRepository.AssertExpectations(t)
			txSongPartRepository.AssertExpectations(t)
			songProcessor.AssertExpectations(t)
		})
	}
}
