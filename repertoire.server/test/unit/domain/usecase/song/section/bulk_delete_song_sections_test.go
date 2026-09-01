package section

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song/section"
	"repertoire/server/internal/wrapper"
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

func TestBulkDeleteSongSections_WhenTransactionExecuteFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	_uut := section.NewBulkDeleteSongSections(transactionManager, nil)

	request := requests.BulkDeleteSongSectionsRequest{
		SongID: uuid.New(),
		IDs:    []uuid.UUID{uuid.New()},
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

func TestBulkDeleteSongSections_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	_uut := section.NewBulkDeleteSongSections(transactionManager, nil)

	request := requests.BulkDeleteSongSectionsRequest{
		SongID: uuid.New(),
		IDs:    []uuid.UUID{uuid.New()},
	}

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepository := new(repository.SongRepositoryMock)
	txSongSectionRepository := new(repository.SongSectionRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepository).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	internalError := errors.New("get song error")
	txSongRepository.On("GetWithSections", new(model.Song), request.SongID).
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
}

func TestBulkDeleteSongSections_WhenSongNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	_uut := section.NewBulkDeleteSongSections(transactionManager, nil)

	request := requests.BulkDeleteSongSectionsRequest{
		SongID: uuid.New(),
		IDs:    []uuid.UUID{uuid.New()},
	}

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepository := new(repository.SongRepositoryMock)
	txSongSectionRepository := new(repository.SongSectionRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepository).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txSongRepository.On("GetWithSections", new(model.Song), request.SongID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())

	repositoryFactory.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	txSongRepository.AssertExpectations(t)
	txSongSectionRepository.AssertExpectations(t)
}

func TestBulkDeleteSongSections_WhenSectionsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	_uut := section.NewBulkDeleteSongSections(transactionManager, nil)

	request := requests.BulkDeleteSongSectionsRequest{
		SongID: uuid.New(),
		IDs:    []uuid.UUID{uuid.New()},
	}

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepository := new(repository.SongRepositoryMock)
	txSongSectionRepository := new(repository.SongSectionRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepository).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	mockSong := &model.Song{
		ID: request.SongID,
		Sections: []model.SongSection{
			{ID: uuid.New()},
		},
	}
	txSongRepository.On("GetWithSections", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song sections not found", errCode.Error.Error())

	repositoryFactory.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	txSongRepository.AssertExpectations(t)
	txSongSectionRepository.AssertExpectations(t)
}

func TestBulkDeleteSongSections_WhenUpdateSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	_uut := section.NewBulkDeleteSongSections(transactionManager, nil)

	request := requests.BulkDeleteSongSectionsRequest{
		SongID: uuid.New(),
		IDs:    []uuid.UUID{uuid.New()},
	}

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepository := new(repository.SongRepositoryMock)
	txSongSectionRepository := new(repository.SongSectionRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepository).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	mockSong := &model.Song{
		ID: request.SongID,
		Sections: []model.SongSection{
			{ID: request.IDs[0]},
		},
	}
	txSongRepository.On("GetWithSections", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	internalError := errors.New("update song error")
	txSongRepository.On("UpdateWithAssociations", mock.IsType(mockSong)).
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
}

// With Parts

func TestBulkDeleteSongSections_WithParts_WhenGetSectionsWithPartsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	songProcessor := new(processor.SongProcessorMock)
	_uut := section.NewBulkDeleteSongSections(transactionManager, songProcessor)

	request := requests.BulkDeleteSongSectionsRequest{
		SongID:    uuid.New(),
		IDs:       []uuid.UUID{uuid.New()},
		WithParts: true,
	}

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepository := new(repository.SongRepositoryMock)
	txSongSectionRepository := new(repository.SongSectionRepositoryMock)
	txSongPartRepository := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepository).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepository).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	mockSong := &model.Song{
		ID: request.SongID,
		Sections: []model.SongSection{
			{ID: request.IDs[0]},
		},
	}
	txSongRepository.On("GetWithSections", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()
	txSongRepository.On("UpdateWithAssociations", mock.IsType(mockSong)).
		Return(nil).
		Once()

	internalError := errors.New("get sections with parts error")
	txSongSectionRepository.On("GetAllByIDsWithSectionParts", new([]model.SongSection), request.IDs).
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

func TestBulkDeleteSongSections_WithParts_WhenUpdateSongAfterPartsDeletionFails_ShouldReturnError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	songProcessor := new(processor.SongProcessorMock)
	_uut := section.NewBulkDeleteSongSections(transactionManager, songProcessor)

	sectionID := uuid.New()
	partID := uuid.New()
	request := requests.BulkDeleteSongSectionsRequest{
		SongID:    uuid.New(),
		IDs:       []uuid.UUID{sectionID},
		WithParts: true,
	}

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepository := new(repository.SongRepositoryMock)
	txSongSectionRepository := new(repository.SongSectionRepositoryMock)
	txSongPartRepository := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepository).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepository).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	mockSong := &model.Song{
		ID: request.SongID,
		Sections: []model.SongSection{
			{
				ID: sectionID,
				SectionParts: []model.SongSectionPart{
					{PartID: partID},
				},
			},
		},
	}
	txSongRepository.On("GetWithSections", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()
	txSongRepository.On("UpdateWithAssociations", mock.IsType(mockSong)).
		Return(nil).
		Once()

	mockSections := []model.SongSection{
		{
			ID: sectionID,
			SectionParts: []model.SongSectionPart{
				{PartID: partID},
			},
			SongID: request.SongID,
		},
	}
	txSongSectionRepository.On("GetAllByIDsWithSectionParts", new([]model.SongSection), request.IDs).
		Return(nil, &mockSections).
		Once()

	expectedError := wrapper.InternalServerError(errors.New("processor error"))
	songProcessor.On("UpdateSongAfterPartsDeletion", txSongRepository, request.SongID, []uuid.UUID{partID}).
		Return(expectedError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, expectedError, errCode)

	repositoryFactory.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	txSongRepository.AssertExpectations(t)
	txSongSectionRepository.AssertExpectations(t)
	txSongPartRepository.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
}

func TestBulkDeleteSongSections_WithParts_WhenDeletePartsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	songProcessor := new(processor.SongProcessorMock)
	_uut := section.NewBulkDeleteSongSections(transactionManager, songProcessor)

	sectionID := uuid.New()
	partID := uuid.New()
	request := requests.BulkDeleteSongSectionsRequest{
		SongID:    uuid.New(),
		IDs:       []uuid.UUID{sectionID},
		WithParts: true,
	}

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepository := new(repository.SongRepositoryMock)
	txSongSectionRepository := new(repository.SongSectionRepositoryMock)
	txSongPartRepository := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepository).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepository).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	mockSong := &model.Song{
		ID: request.SongID,
		Sections: []model.SongSection{
			{
				ID: sectionID,
				SectionParts: []model.SongSectionPart{
					{PartID: partID},
				},
			},
		},
	}
	txSongRepository.On("GetWithSections", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()
	txSongRepository.On("UpdateWithAssociations", mock.IsType(mockSong)).
		Return(nil).
		Once()

	mockSections := []model.SongSection{
		{
			ID: sectionID,
			SectionParts: []model.SongSectionPart{
				{PartID: partID},
			},
			SongID: request.SongID,
		},
	}
	txSongSectionRepository.On("GetAllByIDsWithSectionParts", new([]model.SongSection), request.IDs).
		Return(nil, &mockSections).
		Once()

	songProcessor.On("UpdateSongAfterPartsDeletion", txSongRepository, request.SongID, []uuid.UUID{partID}).
		Return(nil).
		Once()

	internalError := errors.New("delete parts error")
	txSongPartRepository.On("Delete", []uuid.UUID{partID}).
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

func TestBulkDeleteSongSections_WhenDeleteSectionsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	_uut := section.NewBulkDeleteSongSections(transactionManager, nil)

	request := requests.BulkDeleteSongSectionsRequest{
		SongID:    uuid.New(),
		IDs:       []uuid.UUID{uuid.New()},
		WithParts: false,
	}

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepository := new(repository.SongRepositoryMock)
	txSongSectionRepository := new(repository.SongSectionRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepository).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	mockSong := &model.Song{
		ID: request.SongID,
		Sections: []model.SongSection{
			{ID: request.IDs[0]},
		},
	}
	txSongRepository.On("GetWithSections", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()
	txSongRepository.On("UpdateWithAssociations", mock.IsType(mockSong)).
		Return(nil).
		Once()

	internalError := errors.New("delete sections error")
	txSongSectionRepository.On("Delete", request.IDs).
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
}

func TestBulkDeleteSongSections_WhenSuccessful_ShouldDeleteSections(t *testing.T) {
	tests := []struct {
		name          string
		song          model.Song
		deleteIndices []int
		withParts     bool
	}{
		{
			name: "Delete single middle section",
			song: model.Song{
				ID: uuid.New(),
				Sections: []model.SongSection{
					{ID: uuid.New(), Order: 0, SectionParts: []model.SongSectionPart{{PartID: uuid.New()}}},
					{ID: uuid.New(), Order: 1, SectionParts: []model.SongSectionPart{{PartID: uuid.New()}}},
					{ID: uuid.New(), Order: 2, SectionParts: []model.SongSectionPart{{PartID: uuid.New()}}},
				},
			},
			deleteIndices: []int{1},
		},
		{
			name: "Delete multiple sections",
			song: model.Song{
				ID: uuid.New(),
				Sections: []model.SongSection{
					{ID: uuid.New(), Order: 0},
					{ID: uuid.New(), Order: 1},
					{ID: uuid.New(), Order: 2},
					{ID: uuid.New(), Order: 3},
					{ID: uuid.New(), Order: 4},
				},
			},
			deleteIndices: []int{1, 3},
		},
		{
			name: "Delete first section",
			song: model.Song{
				ID: uuid.New(),
				Sections: []model.SongSection{
					{ID: uuid.New(), Order: 0},
					{ID: uuid.New(), Order: 1},
					{ID: uuid.New(), Order: 2},
				},
			},
			deleteIndices: []int{0},
		},
		{
			name: "Delete last section",
			song: model.Song{
				ID: uuid.New(),
				Sections: []model.SongSection{
					{ID: uuid.New(), Order: 0},
					{ID: uuid.New(), Order: 1},
					{ID: uuid.New(), Order: 2},
				},
			},
			deleteIndices: []int{2},
		},
		{
			name: "Delete all sections",
			song: model.Song{
				ID: uuid.New(),
				Sections: []model.SongSection{
					{ID: uuid.New(), Order: 0},
					{ID: uuid.New(), Order: 1},
				},
			},
			deleteIndices: []int{0, 1},
		},
		{
			name: "Delete sections with parts",
			song: model.Song{
				ID: uuid.New(),
				Sections: []model.SongSection{
					{ID: uuid.New(), Order: 0, SectionParts: []model.SongSectionPart{{PartID: uuid.New()}}},
					{ID: uuid.New(), Order: 1, SectionParts: []model.SongSectionPart{{PartID: uuid.New()}}},
				},
			},
			deleteIndices: []int{1},
			withParts:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			transactionManager := new(transaction.ManagerMock)
			var songProcessor *processor.SongProcessorMock
			var _uut section.BulkDeleteSongSections

			if tt.withParts {
				songProcessor = new(processor.SongProcessorMock)
				_uut = section.NewBulkDeleteSongSections(transactionManager, songProcessor)
			} else {
				_uut = section.NewBulkDeleteSongSections(transactionManager, nil)
			}

			idsToDelete := make([]uuid.UUID, len(tt.deleteIndices))
			for i, idx := range tt.deleteIndices {
				idsToDelete[i] = tt.song.Sections[idx].ID
			}
			request := requests.BulkDeleteSongSectionsRequest{
				SongID:    tt.song.ID,
				IDs:       idsToDelete,
				WithParts: tt.withParts,
			}

			// given - mocking
			repositoryFactory := new(transaction.RepositoryFactoryMock)
			txSongRepository := new(repository.SongRepositoryMock)
			txSongSectionRepository := new(repository.SongSectionRepositoryMock)
			txSongPartRepository := new(repository.SongPartRepositoryMock)

			repositoryFactory.On("NewSongRepository").Return(txSongRepository).Once()
			repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepository).Once()
			if tt.withParts {
				repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepository).Once()
			}
			transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

			txSongRepository.On("GetWithSections", new(model.Song), request.SongID).
				Return(nil, &tt.song).
				Once()

			txSongRepository.On("UpdateWithAssociations", mock.IsType(&tt.song)).
				Run(func(args mock.Arguments) {
					updatedSong := args.Get(0).(*model.Song)
					// sections reordered
					order := uint(0)
					for _, s := range updatedSong.Sections {
						if slices.Contains(request.IDs, s.ID) {
							continue
						}
						assert.Equal(t, order, s.Order)
						order++
					}
				}).
				Return(nil).Once()

			if tt.withParts {
				var mockSections []model.SongSection
				for _, idx := range tt.deleteIndices {
					sec := tt.song.Sections[idx]
					sec.SongID = tt.song.ID
					mockSections = append(mockSections, sec)
				}
				txSongSectionRepository.On("GetAllByIDsWithSectionParts", new([]model.SongSection), idsToDelete).
					Return(nil, &mockSections).
					Once()

				var partIDsToDelete []uuid.UUID
				partSet := make(map[uuid.UUID]bool)
				for _, sec := range mockSections {
					for _, sp := range sec.SectionParts {
						if !partSet[sp.PartID] {
							partSet[sp.PartID] = true
							partIDsToDelete = append(partIDsToDelete, sp.PartID)
						}
					}
				}

				if len(partIDsToDelete) > 0 {
					songProcessor.On("UpdateSongAfterPartsDeletion", txSongRepository, tt.song.ID, partIDsToDelete).
						Return(nil).
						Once()
					txSongPartRepository.On("Delete", partIDsToDelete).
						Return(nil).
						Once()
				}
			}

			// Delete sections
			txSongSectionRepository.On("Delete", idsToDelete).
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
			if tt.withParts {
				txSongPartRepository.AssertExpectations(t)
				songProcessor.AssertExpectations(t)
			}
		})
	}
}
