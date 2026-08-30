package part

import (
	"errors"
	"net/http"
	"repertoire/server/domain/usecase/song/part"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/database/transaction"
	"repertoire/server/test/unit/data/repository"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDeleteSongPart_WhenTransactionExecuteFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewDeleteSongPart(transactionManager)

	id := uuid.New()
	songID := uuid.New()

	internalError := errors.New("transaction error")
	transactionManager.On("Execute", mock.Anything).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	transactionManager.AssertExpectations(t)
}

func TestDeleteSongPart_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewDeleteSongPart(transactionManager)

	id := uuid.New()
	songID := uuid.New()

	// given - mocking
	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepository := new(repository.SongRepositoryMock)
	txSongSectionRepository := new(repository.SongSectionRepositoryMock)
	txSongPartRepository := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepository).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepository).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	internalError := errors.New("get error")
	txSongRepository.On("GetWithParts", new(model.Song), songID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	repositoryFactory.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	txSongRepository.AssertExpectations(t)
	txSongSectionRepository.AssertExpectations(t)
	txSongPartRepository.AssertExpectations(t)
}

func TestDeleteSongPart_WhenSongNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewDeleteSongPart(transactionManager)

	id := uuid.New()
	songID := uuid.New()

	// given - mocking
	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepository := new(repository.SongRepositoryMock)
	txSongSectionRepository := new(repository.SongSectionRepositoryMock)
	txSongPartRepository := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepository).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepository).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txSongRepository.On("GetWithParts", new(model.Song), songID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())

	repositoryFactory.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	txSongRepository.AssertExpectations(t)
	txSongSectionRepository.AssertExpectations(t)
	txSongPartRepository.AssertExpectations(t)
}

func TestDeleteSongPart_WhenPartNotFoundInSong_ShouldReturnNotFoundError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewDeleteSongPart(transactionManager)

	id := uuid.New()
	songID := uuid.New()

	// given - mocking
	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepository := new(repository.SongRepositoryMock)
	txSongSectionRepository := new(repository.SongSectionRepositoryMock)
	txSongPartRepository := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepository).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepository).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	mockSong := &model.Song{
		ID:    songID,
		Parts: []model.SongPart{{ID: uuid.New()}},
	}
	txSongRepository.On("GetWithParts", new(model.Song), songID).
		Return(nil, mockSong).
		Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song part not found", errCode.Error.Error())

	repositoryFactory.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	txSongRepository.AssertExpectations(t)
	txSongSectionRepository.AssertExpectations(t)
	txSongPartRepository.AssertExpectations(t)
}

func TestDeleteSongPart_WhenUpdateSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewDeleteSongPart(transactionManager)

	id := uuid.New()
	songID := uuid.New()

	// given - mocking
	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepository := new(repository.SongRepositoryMock)
	txSongSectionRepository := new(repository.SongSectionRepositoryMock)
	txSongPartRepository := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepository).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepository).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	mockSong := &model.Song{
		ID:    songID,
		Parts: []model.SongPart{{ID: id}},
	}
	txSongRepository.On("GetWithParts", new(model.Song), songID).
		Return(nil, mockSong).
		Once()

	internalError := errors.New("update error")
	txSongRepository.On("UpdateWithAssociations", mock.IsType(mockSong)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	repositoryFactory.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	txSongRepository.AssertExpectations(t)
	txSongSectionRepository.AssertExpectations(t)
	txSongPartRepository.AssertExpectations(t)
}

func TestDeleteSongPart_WhenGetSectionsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewDeleteSongPart(transactionManager)

	id := uuid.New()
	songID := uuid.New()

	// given - mocking
	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepository := new(repository.SongRepositoryMock)
	txSongSectionRepository := new(repository.SongSectionRepositoryMock)
	txSongPartRepository := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepository).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepository).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	// updateSong
	mockSong := &model.Song{
		ID:    songID,
		Parts: []model.SongPart{{ID: id}},
	}
	txSongRepository.On("GetWithParts", new(model.Song), songID).
		Return(nil, mockSong).
		Once()
	txSongRepository.On("UpdateWithAssociations", mock.IsType(mockSong)).
		Return(nil).
		Once()

	// updateSections
	internalError := errors.New("get sections error")
	txSongSectionRepository.On("GetAllByPartWithSectionParts", new([]model.SongSection), id).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	repositoryFactory.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	txSongRepository.AssertExpectations(t)
	txSongSectionRepository.AssertExpectations(t)
	txSongPartRepository.AssertExpectations(t)
}

func TestDeleteSongPart_WhenUpdateAllSectionPartsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewDeleteSongPart(transactionManager)

	id := uuid.New()
	songID := uuid.New()
	sectionID := uuid.New()

	// given - mocking
	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepository := new(repository.SongRepositoryMock)
	txSongSectionRepository := new(repository.SongSectionRepositoryMock)
	txSongPartRepository := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepository).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepository).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	// updateSong
	mockSong := &model.Song{
		ID:    songID,
		Parts: []model.SongPart{{ID: id}},
	}
	txSongRepository.On("GetWithParts", new(model.Song), songID).
		Return(nil, mockSong).
		Once()
	txSongRepository.On("UpdateWithAssociations", mock.IsType(mockSong)).
		Return(nil).
		Once()

	// updateSections
	mockSections := []model.SongSection{
		{
			ID: sectionID,
			SectionParts: []model.SongSectionPart{
				{PartID: id, SectionID: sectionID, Order: 0},
				{PartID: uuid.New(), SectionID: sectionID, Order: 1},
			},
		},
	}
	txSongSectionRepository.On("GetAllByPartWithSectionParts", new([]model.SongSection), id).
		Return(nil, &mockSections).
		Once()

	internalError := errors.New("update section parts error")
	txSongSectionRepository.On("UpdateAllSectionParts", mock.Anything).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	repositoryFactory.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	txSongRepository.AssertExpectations(t)
	txSongSectionRepository.AssertExpectations(t)
	txSongPartRepository.AssertExpectations(t)
}

func TestDeleteSongPart_WhenDeletePartFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewDeleteSongPart(transactionManager)

	id := uuid.New()
	songID := uuid.New()

	// given - mocking
	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepository := new(repository.SongRepositoryMock)
	txSongSectionRepository := new(repository.SongSectionRepositoryMock)
	txSongPartRepository := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepository).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepository).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	// updateSong
	mockSong := &model.Song{
		ID:    songID,
		Parts: []model.SongPart{{ID: id}},
	}
	txSongRepository.On("GetWithParts", new(model.Song), songID).
		Return(nil, mockSong).
		Once()
	txSongRepository.On("UpdateWithAssociations", mock.IsType(mockSong)).
		Return(nil).
		Once()

	// updateSections
	txSongSectionRepository.On("GetAllByPartWithSectionParts", new([]model.SongSection), id).
		Return(nil, &[]model.SongSection{}).
		Once()

	internalError := errors.New("delete error")
	txSongPartRepository.On("Delete", []uuid.UUID{id}).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	repositoryFactory.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	txSongRepository.AssertExpectations(t)
	txSongSectionRepository.AssertExpectations(t)
	txSongPartRepository.AssertExpectations(t)
}

func TestDeleteSongPart_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	tests := []struct {
		name                   string
		song                   model.Song
		partIndex              int
		sections               []model.SongSection
		sectionPartIndices     []int
		expectedSongConfidence float64
		expectedSongRehearsals float64
		expectedSongProgress   float64
	}{
		{
			name: "Only part",
			song: model.Song{
				ID: uuid.New(),
				Parts: []model.SongPart{
					{ID: uuid.New(), Confidence: 0, Rehearsals: 0, Progress: 0},
				},
			},
			partIndex:              0,
			expectedSongConfidence: 0,
			expectedSongRehearsals: 0,
			expectedSongProgress:   0,
		},
		{
			name: "Only part with stats",
			song: model.Song{
				ID: uuid.New(),
				Parts: []model.SongPart{
					{ID: uuid.New(), Confidence: 55, Rehearsals: 25, Progress: 39},
				},
				Confidence: 55,
				Rehearsals: 25,
				Progress:   39,
			},
			partIndex:              0,
			expectedSongConfidence: 0,
			expectedSongRehearsals: 0,
			expectedSongProgress:   0,
		},
		{
			name: "Multiple parts, middle one deleted",
			song: model.Song{
				ID: uuid.New(),
				Parts: []model.SongPart{
					{ID: uuid.New(), SongOrder: 0, Confidence: 0, Rehearsals: 0, Progress: 0},
					{ID: uuid.New(), SongOrder: 1, Confidence: 0, Rehearsals: 0, Progress: 0},
					{ID: uuid.New(), SongOrder: 2, Confidence: 0, Rehearsals: 0, Progress: 0},
					{ID: uuid.New(), SongOrder: 3, Confidence: 0, Rehearsals: 0, Progress: 0},
				},
				Confidence: 0,
				Rehearsals: 0,
				Progress:   0,
			},
			partIndex:              2,
			expectedSongConfidence: 0,
			expectedSongRehearsals: 0,
			expectedSongProgress:   0,
		},
		{
			name: "Multiple parts with stats, middle one deleted",
			song: model.Song{
				ID: uuid.New(),
				Parts: []model.SongPart{
					{ID: uuid.New(), SongOrder: 0, Confidence: 55, Rehearsals: 12, Progress: 45},
					{ID: uuid.New(), SongOrder: 1, Confidence: 23, Rehearsals: 5, Progress: 15},
					{ID: uuid.New(), SongOrder: 2, Confidence: 78, Rehearsals: 25, Progress: 100},
					{ID: uuid.New(), SongOrder: 3, Confidence: 40, Rehearsals: 6, Progress: 63},
					{ID: uuid.New(), SongOrder: 4, Confidence: 80, Rehearsals: 19, Progress: 170},
				},
				Confidence: 55.2,
				Rehearsals: 13.4,
				Progress:   78.6,
			},
			partIndex:              2,
			expectedSongConfidence: 49.5,
			expectedSongRehearsals: 10.5,
			expectedSongProgress:   73.25,
		},
		{
			name: "With one section, deleted part is second in section",
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
			partIndex: 1,
			sections: []model.SongSection{
				{
					ID: uuid.New(),
					SectionParts: []model.SongSectionPart{
						{PartID: uuid.New(), Order: 0},
						{PartID: uuid.New(), Order: 1},
						{PartID: uuid.New(), Order: 2},
					},
				},
			},
			sectionPartIndices:     []int{1}, // deleted part is at index 1
			expectedSongConfidence: 30,
			expectedSongRehearsals: 8,
			expectedSongProgress:   20,
		},
		{
			name: "With multiple sections, deleted part has different positions",
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
			partIndex: 1,
			sections: []model.SongSection{
				{
					ID: uuid.New(),
					SectionParts: []model.SongSectionPart{
						{PartID: uuid.New(), Order: 0},
						{PartID: uuid.New(), Order: 1},
						{PartID: uuid.New(), Order: 2},
					},
				},
				{
					ID: uuid.New(),
					SectionParts: []model.SongSectionPart{
						{PartID: uuid.New(), Order: 0},
						{PartID: uuid.New(), Order: 1},
					},
				},
			},
			sectionPartIndices:     []int{1, 0}, // in first section: index 1; in second: index 0
			expectedSongConfidence: 37.5,
			expectedSongRehearsals: 9,
			expectedSongProgress:   35,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			transactionManager := new(transaction.ManagerMock)
			_uut := part.NewDeleteSongPart(transactionManager)

			song := tt.song
			partToDelete := song.Parts[tt.partIndex]
			id := partToDelete.ID
			songID := song.ID

			// given - mocking
			repositoryFactory := new(transaction.RepositoryFactoryMock)
			txSongRepository := new(repository.SongRepositoryMock)
			txSongSectionRepository := new(repository.SongSectionRepositoryMock)
			txSongPartRepository := new(repository.SongPartRepositoryMock)

			repositoryFactory.On("NewSongRepository").Return(txSongRepository).Once()
			repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepository).Once()
			repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepository).Once()
			transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

			// updateSong
			txSongRepository.On("GetWithParts", new(model.Song), songID).
				Return(nil, &song).
				Once()

			txSongRepository.On("UpdateWithAssociations", mock.IsType(&song)).
				Run(func(args mock.Arguments) {
					updatedSong := args.Get(0).(*model.Song)

					assert.Equal(t, tt.expectedSongConfidence, updatedSong.Confidence)
					assert.Equal(t, tt.expectedSongRehearsals, updatedSong.Rehearsals)
					assert.Equal(t, tt.expectedSongProgress, updatedSong.Progress)

					parts := slices.Clone(updatedSong.Parts)
					parts = slices.DeleteFunc(parts, func(p model.SongPart) bool {
						return p.ID == id
					})
					for i, p := range parts {
						assert.Equal(t, uint(i), p.SongOrder)
					}
				}).
				Return(nil).Once()

			// updateSections
			if len(tt.sections) > 0 {
				var mockSections []model.SongSection
				var expectedSectionParts []model.SongSectionPart
				for i, sec := range tt.sections {
					idx := 0
					if i < len(tt.sectionPartIndices) {
						idx = tt.sectionPartIndices[i]
					}

					// Prepare the section with the deleted part ID set at the correct index
					var newSectionParts []model.SongSectionPart
					for j, sp := range sec.SectionParts {
						newSectionParts = append(newSectionParts, sp)
						if j == idx {
							newSectionParts[j].PartID = id
						}
					}
					sec.SectionParts = newSectionParts
					mockSections = append(mockSections, sec)

					// Build expected slice for this section: entries after idx with Order decremented by 1
					for j := idx + 1; j < len(sec.SectionParts); j++ {
						expectedPart := sec.SectionParts[j]
						expectedPart.Order = expectedPart.Order - 1
						expectedSectionParts = append(expectedSectionParts, expectedPart)
					}
				}

				txSongSectionRepository.On("GetAllByPartWithSectionParts", new([]model.SongSection), id).
					Return(nil, &mockSections).
					Once()

				txSongSectionRepository.On("UpdateAllSectionParts", mock.IsType(&expectedSectionParts)).
					Run(func(args mock.Arguments) {
						updatedSectionParts := args.Get(0).(*[]model.SongSectionPart)
						assert.Equal(t, expectedSectionParts, *updatedSectionParts)
					}).
					Return(nil).
					Once()
			} else {
				txSongSectionRepository.On("GetAllByPartWithSectionParts", new([]model.SongSection), id).
					Return(nil, &[]model.SongSection{}).
					Once()
			}

			txSongPartRepository.On("Delete", []uuid.UUID{id}).
				Return(nil).
				Once()

			// when
			errCode := _uut.Handle(id, songID)

			// then
			assert.Nil(t, errCode)

			repositoryFactory.AssertExpectations(t)
			transactionManager.AssertExpectations(t)
			txSongRepository.AssertExpectations(t)
			txSongSectionRepository.AssertExpectations(t)
			txSongPartRepository.AssertExpectations(t)
		})
	}
}
