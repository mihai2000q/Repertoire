package section

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song/section"
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
	txSongRepo := new(repository.SongRepositoryMock)
	txSongSectionRepo := new(repository.SongSectionRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	internalError := errors.New("get song error")
	txSongRepo.On("GetWithSections", new(model.Song), request.SongID).
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
	txSongRepo.AssertExpectations(t)
	txSongSectionRepo.AssertExpectations(t)
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
	txSongRepo := new(repository.SongRepositoryMock)
	txSongSectionRepo := new(repository.SongSectionRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txSongRepo.On("GetWithSections", new(model.Song), request.SongID).
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
	txSongRepo.AssertExpectations(t)
	txSongSectionRepo.AssertExpectations(t)
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
	txSongRepo := new(repository.SongRepositoryMock)
	txSongSectionRepo := new(repository.SongSectionRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	mockSong := &model.Song{
		ID: request.SongID,
		Sections: []model.SongSection{
			{ID: uuid.New()},
		},
	}
	txSongRepo.On("GetWithSections", new(model.Song), request.SongID).
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
	txSongRepo.AssertExpectations(t)
	txSongSectionRepo.AssertExpectations(t)
}

func TestBulkDeleteSongSections_WhenGetSectionsWithPartsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	songProcessor := new(processor.SongProcessorMock)
	_uut := section.NewBulkDeleteSongSections(transactionManager, songProcessor)

	request := requests.BulkDeleteSongSectionsRequest{
		SongID:  uuid.New(),
		IDs:     []uuid.UUID{uuid.New()},
		PartIDs: []uuid.UUID{uuid.New()},
	}

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongSectionRepo := new(repository.SongSectionRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	mockSong := &model.Song{
		ID: request.SongID,
		Sections: []model.SongSection{
			{ID: request.IDs[0]},
		},
	}
	txSongRepo.On("GetWithSections", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	internalError := errors.New("get sections with parts error")
	txSongSectionRepo.On("GetAllByIDsWithSectionParts", new([]model.SongSection), request.IDs).
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
	txSongRepo.AssertExpectations(t)
	txSongSectionRepo.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
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
	txSongRepo := new(repository.SongRepositoryMock)
	txSongSectionRepo := new(repository.SongSectionRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	mockSong := &model.Song{
		ID: request.SongID,
		Sections: []model.SongSection{
			{ID: request.IDs[0]},
		},
	}
	txSongRepo.On("GetWithSections", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	internalError := errors.New("update song error")
	txSongRepo.On("UpdateWithAssociations", mock.IsType(mockSong)).
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
	txSongRepo.AssertExpectations(t)
	txSongSectionRepo.AssertExpectations(t)
}

func TestBulkDeleteSongSections_WhenUpdateSongAfterPartsDeletionFails_ShouldReturnError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	songProcessor := new(processor.SongProcessorMock)
	_uut := section.NewBulkDeleteSongSections(transactionManager, songProcessor)

	sectionID := uuid.New()
	partID := uuid.New()
	request := requests.BulkDeleteSongSectionsRequest{
		SongID:  uuid.New(),
		IDs:     []uuid.UUID{sectionID},
		PartIDs: []uuid.UUID{partID},
	}

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongSectionRepo := new(repository.SongSectionRepositoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepo).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
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
	txSongRepo.On("GetWithSections", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()
	txSongRepo.On("UpdateWithAssociations", mock.IsType(mockSong)).
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
	txSongSectionRepo.On("GetAllByIDsWithSectionParts", new([]model.SongSection), request.IDs).
		Return(nil, &mockSections).
		Once()

	expectedError := httperror.InternalServerError(errors.New("processor error"))
	songProcessor.On("UpdateSongAfterPartsDeletion", txSongRepo, request.SongID, []uuid.UUID{partID}).
		Return(expectedError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, expectedError, errCode)

	repositoryFactory.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongSectionRepo.AssertExpectations(t)
	txSongPartRepo.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
}

func TestBulkDeleteSongSections_WhenDeletePartsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	songProcessor := new(processor.SongProcessorMock)
	_uut := section.NewBulkDeleteSongSections(transactionManager, songProcessor)

	sectionID := uuid.New()
	partID := uuid.New()
	request := requests.BulkDeleteSongSectionsRequest{
		SongID:  uuid.New(),
		IDs:     []uuid.UUID{sectionID},
		PartIDs: []uuid.UUID{partID},
	}

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongSectionRepo := new(repository.SongSectionRepositoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepo).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
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
	txSongRepo.On("GetWithSections", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()
	txSongRepo.On("UpdateWithAssociations", mock.IsType(mockSong)).
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
	txSongSectionRepo.On("GetAllByIDsWithSectionParts", new([]model.SongSection), request.IDs).
		Return(nil, &mockSections).
		Once()

	songProcessor.On("UpdateSongAfterPartsDeletion", txSongRepo, request.SongID, request.PartIDs).
		Return(nil).
		Once()

	internalError := errors.New("delete parts error")
	txSongPartRepo.On("Delete", request.PartIDs).
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
	txSongRepo.AssertExpectations(t)
	txSongSectionRepo.AssertExpectations(t)
	txSongPartRepo.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
}

func TestBulkDeleteSongSections_WhenDeleteSectionsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	_uut := section.NewBulkDeleteSongSections(transactionManager, nil)

	request := requests.BulkDeleteSongSectionsRequest{
		SongID: uuid.New(),
		IDs:    []uuid.UUID{uuid.New()},
	}

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongSectionRepo := new(repository.SongSectionRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	mockSong := &model.Song{
		ID: request.SongID,
		Sections: []model.SongSection{
			{ID: request.IDs[0]},
		},
	}
	txSongRepo.On("GetWithSections", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()
	txSongRepo.On("UpdateWithAssociations", mock.IsType(mockSong)).
		Return(nil).
		Once()

	internalError := errors.New("delete sections error")
	txSongSectionRepo.On("Delete", request.IDs).
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
	txSongRepo.AssertExpectations(t)
	txSongSectionRepo.AssertExpectations(t)
}
func TestBulkDeleteSongSections_WhenSuccessful_ShouldDeleteSections(t *testing.T) {
	// partRef points at song.Sections[section].SectionParts[part]
	type partRef struct{ section, part int }
	sharedPartID := uuid.New()

	tests := []struct {
		name          string
		song          model.Song
		deleteIndices []int
		partsToDelete []partRef
	}{
		{
			name: "Delete sections with all their parts",
			song: model.Song{
				ID: uuid.New(),
				Sections: []model.SongSection{
					{ID: uuid.New(), Order: 0, SectionParts: []model.SongSectionPart{{PartID: uuid.New()}}},
					{ID: uuid.New(), Order: 1, SectionParts: []model.SongSectionPart{{PartID: uuid.New()}}},
				},
			},
			deleteIndices: []int{1},
			partsToDelete: []partRef{{1, 0}},
		},
		{
			name: "Delete sections with only some of their parts",
			song: model.Song{
				ID: uuid.New(),
				Sections: []model.SongSection{
					{ID: uuid.New(), Order: 0, SectionParts: []model.SongSectionPart{{PartID: uuid.New()}, {PartID: uuid.New()}}},
					{ID: uuid.New(), Order: 1, SectionParts: []model.SongSectionPart{{PartID: uuid.New()}}},
				},
			},
			deleteIndices: []int{0},
			partsToDelete: []partRef{{0, 0}},
		},
		{
			name: "Delete sections sharing a part",
			song: model.Song{
				ID: uuid.New(),
				Sections: []model.SongSection{
					{ID: uuid.New(), Order: 0, SectionParts: []model.SongSectionPart{{PartID: sharedPartID}, {PartID: uuid.New()}}},
					{ID: uuid.New(), Order: 1, SectionParts: []model.SongSectionPart{{PartID: sharedPartID}}},
					{ID: uuid.New(), Order: 2, SectionParts: []model.SongSectionPart{{PartID: uuid.New()}}},
				},
			},
			deleteIndices: []int{0, 1},
			partsToDelete: []partRef{{0, 0}, {0, 1}, {1, 0}}, // the shared part is sent twice
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			idsToDelete := make([]uuid.UUID, len(tt.deleteIndices))
			for i, idx := range tt.deleteIndices {
				idsToDelete[i] = tt.song.Sections[idx].ID
			}
			partIDsToDelete := make([]uuid.UUID, len(tt.partsToDelete))
			for i, ref := range tt.partsToDelete {
				partIDsToDelete[i] = tt.song.Sections[ref.section].SectionParts[ref.part].PartID
			}
			withParts := len(partIDsToDelete) > 0

			transactionManager := new(transaction.ManagerMock)
			var songProcessor *processor.SongProcessorMock
			var _uut section.BulkDeleteSongSections
			if withParts {
				songProcessor = new(processor.SongProcessorMock)
				_uut = section.NewBulkDeleteSongSections(transactionManager, songProcessor)
			} else {
				_uut = section.NewBulkDeleteSongSections(transactionManager, nil)
			}

			request := requests.BulkDeleteSongSectionsRequest{
				SongID:  tt.song.ID,
				IDs:     idsToDelete,
				PartIDs: partIDsToDelete,
			}

			// given - mocking
			repositoryFactory := new(transaction.RepositoryFactoryMock)
			txSongRepo := new(repository.SongRepositoryMock)
			txSongSectionRepo := new(repository.SongSectionRepositoryMock)
			txSongPartRepo := new(repository.SongPartRepositoryMock)

			repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
			repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepo).Once()
			if withParts {
				repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
			}
			transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

			txSongRepo.On("GetWithSections", new(model.Song), request.SongID).
				Return(nil, &tt.song).
				Once()

			txSongRepo.On("UpdateWithAssociations", mock.IsType(&tt.song)).
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

			if withParts {
				var mockSections []model.SongSection
				for _, idx := range tt.deleteIndices {
					mockSections = append(mockSections, tt.song.Sections[idx])
				}
				txSongSectionRepo.On("GetAllByIDsWithSectionParts", new([]model.SongSection), idsToDelete).
					Return(nil, &mockSections).
					Once()

				songProcessor.On("UpdateSongAfterPartsDeletion", txSongRepo, tt.song.ID, partIDsToDelete).
					Return(nil).
					Once()
				txSongPartRepo.On("Delete", partIDsToDelete).
					Return(nil).
					Once()
			}

			// Delete sections
			txSongSectionRepo.On("Delete", idsToDelete).
				Return(nil).
				Once()

			// when
			errCode := _uut.Handle(request)

			// then
			assert.Nil(t, errCode)

			repositoryFactory.AssertExpectations(t)
			transactionManager.AssertExpectations(t)
			txSongRepo.AssertExpectations(t)
			txSongSectionRepo.AssertExpectations(t)
			if withParts {
				txSongPartRepo.AssertExpectations(t)
				songProcessor.AssertExpectations(t)
			}
		})
	}
}
