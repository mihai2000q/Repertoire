package section

import (
	"errors"
	"net/http"
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

func TestDeleteSongSection_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewDeleteSongSection(songRepository, nil, nil)

	id := uuid.New()
	songID := uuid.New()

	internalError := errors.New("internal error")
	songRepository.On("GetWithSections", new(model.Song), songID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, songID, false)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestDeleteSongSection_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewDeleteSongSection(songRepository, nil, nil)

	id := uuid.New()
	songID := uuid.New()

	songRepository.On("GetWithSections", new(model.Song), songID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(id, songID, false)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestDeleteSongSection_WhenSectionIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewDeleteSongSection(songRepository, nil, nil)

	id := uuid.New()
	songID := uuid.New()

	// given - mocking
	mockSong := &model.Song{
		ID: songID,
		Sections: []model.SongSection{
			{ID: uuid.New(), Order: 0},
		},
	}
	songRepository.On("GetWithSections", new(model.Song), songID).
		Return(nil, mockSong).
		Once()

	// when
	errCode := _uut.Handle(id, songID, false)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song section not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestDeleteSongSection_WhenTransactionExecuteFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := section.NewDeleteSongSection(songRepository, transactionManager, nil)

	id := uuid.New()
	songID := uuid.New()

	mockSong := &model.Song{
		ID: songID,
		Sections: []model.SongSection{
			{ID: id, Order: 0},
		},
	}
	songRepository.On("GetWithSections", new(model.Song), songID).
		Return(nil, mockSong).
		Once()

	internalError := errors.New("transaction error")
	transactionManager.On("Execute", mock.Anything).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, songID, false)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
}

func TestDeleteSongSection_WhenUpdateSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := section.NewDeleteSongSection(songRepository, transactionManager, nil)

	id := uuid.New()
	songID := uuid.New()

	mockSong := &model.Song{
		ID: songID,
		Sections: []model.SongSection{
			{ID: id, Order: 0},
		},
	}
	songRepository.On("GetWithSections", new(model.Song), songID).
		Return(nil, mockSong).
		Once()

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongSectionRepo := new(repository.SongSectionRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	internalError := errors.New("update error")
	txSongRepo.On("UpdateWithAssociations", mock.IsType(mockSong)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, songID, false)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongSectionRepo.AssertExpectations(t)
}

func TestDeleteSongSection_WhenDeleteSectionFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := section.NewDeleteSongSection(songRepository, transactionManager, nil)

	id := uuid.New()
	songID := uuid.New()

	mockSong := &model.Song{
		ID: songID,
		Sections: []model.SongSection{
			{ID: id, Order: 0},
		},
	}
	songRepository.On("GetWithSections", new(model.Song), songID).
		Return(nil, mockSong).
		Once()

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongSectionRepo := new(repository.SongSectionRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txSongRepo.On("UpdateWithAssociations", mock.IsType(mockSong)).
		Return(nil).
		Once()

	internalError := errors.New("delete section error")
	txSongSectionRepo.On("Delete", []uuid.UUID{id}).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, songID, false)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongSectionRepo.AssertExpectations(t)
}

// With Parts

func TestDeleteSongSection_WithParts_WhenGetSectionWithPartsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	songProcessor := new(processor.SongProcessorMock)
	_uut := section.NewDeleteSongSection(songRepository, transactionManager, songProcessor)

	id := uuid.New()
	songID := uuid.New()

	mockSong := &model.Song{
		ID: songID,
		Sections: []model.SongSection{
			{ID: id, Order: 0},
		},
	}
	songRepository.On("GetWithSections", new(model.Song), songID).
		Return(nil, mockSong).
		Once()

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongSectionRepo := new(repository.SongSectionRepositoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepo).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txSongRepo.On("UpdateWithAssociations", mock.IsType(mockSong)).
		Return(nil).
		Once()

	internalError := errors.New("get section parts error")
	txSongSectionRepo.On("GetWithSectionParts", new(model.SongSection), id).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, songID, true)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongSectionRepo.AssertExpectations(t)
	txSongPartRepo.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
}

func TestDeleteSongSection_WithParts_WhenUpdateSongAfterPartsDeletionFails_ShouldReturnError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	songProcessor := new(processor.SongProcessorMock)
	_uut := section.NewDeleteSongSection(songRepository, transactionManager, songProcessor)

	id := uuid.New()
	songID := uuid.New()
	partID := uuid.New()

	mockSong := &model.Song{
		ID: songID,
		Sections: []model.SongSection{
			{
				ID:    id,
				Order: 0,
				SectionParts: []model.SongSectionPart{
					{PartID: partID},
				},
				SongID: songID,
			},
		},
	}
	songRepository.On("GetWithSections", new(model.Song), songID).
		Return(nil, mockSong).
		Once()

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongSectionRepo := new(repository.SongSectionRepositoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepo).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txSongRepo.On("UpdateWithAssociations", mock.IsType(mockSong)).
		Return(nil).
		Once()
	txSongSectionRepo.On("GetWithSectionParts", new(model.SongSection), id).
		Return(nil, &mockSong.Sections[0]).
		Once()

	expectedError := httperror.InternalServerError(errors.New("processor error"))
	songProcessor.On("UpdateSongAfterPartsDeletion", txSongRepo, songID, []uuid.UUID{partID}).
		Return(expectedError).
		Once()

	// when
	errCode := _uut.Handle(id, songID, true)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, expectedError, errCode)

	songRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongSectionRepo.AssertExpectations(t)
	txSongPartRepo.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
}

func TestDeleteSongSection_WithParts_WhenDeletePartsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	songProcessor := new(processor.SongProcessorMock)
	_uut := section.NewDeleteSongSection(songRepository, transactionManager, songProcessor)

	id := uuid.New()
	songID := uuid.New()
	partID := uuid.New()

	mockSong := &model.Song{
		ID: songID,
		Sections: []model.SongSection{
			{
				ID:           id,
				Order:        0,
				SectionParts: []model.SongSectionPart{{PartID: partID}},
				SongID:       songID,
			},
		},
	}
	songRepository.On("GetWithSections", new(model.Song), songID).
		Return(nil, mockSong).
		Once()

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongSectionRepo := new(repository.SongSectionRepositoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepo).Once()
	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txSongRepo.On("UpdateWithAssociations", mock.IsType(mockSong)).
		Return(nil).
		Once()
	txSongSectionRepo.On("GetWithSectionParts", new(model.SongSection), id).
		Return(nil, &mockSong.Sections[0]).
		Once()

	songProcessor.On("UpdateSongAfterPartsDeletion", txSongRepo, songID, []uuid.UUID{partID}).
		Return(nil).
		Once()

	internalError := errors.New("delete parts error")
	txSongPartRepo.On("Delete", []uuid.UUID{partID}).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, songID, true)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongSectionRepo.AssertExpectations(t)
	txSongPartRepo.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
}

func TestDeleteSongSection_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	tests := []struct {
		name          string
		song          model.Song
		sectionsIndex int
		withParts     bool
		partIDs       []uuid.UUID
	}{
		{
			name: "Delete only section",
			song: model.Song{
				ID: uuid.New(),
				Sections: []model.SongSection{
					{ID: uuid.New(), Order: 0},
				},
			},
			sectionsIndex: 0,
		},
		{
			name: "Delete middle section, reorder others",
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
			sectionsIndex: 2,
		},
		{
			name: "Delete section and its parts",
			song: model.Song{
				ID: uuid.New(),
				Sections: []model.SongSection{
					{ID: uuid.New(), Order: 0},
					{ID: uuid.New(), Order: 1},
				},
			},
			sectionsIndex: 1,
			withParts:     true,
			partIDs:       []uuid.UUID{uuid.New(), uuid.New()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			songRepository := new(repository.SongRepositoryMock)
			transactionManager := new(transaction.ManagerMock)
			songProcessor := new(processor.SongProcessorMock)
			_uut := section.NewDeleteSongSection(songRepository, transactionManager, songProcessor)

			sectionToDelete := tt.song.Sections[tt.sectionsIndex]
			id := sectionToDelete.ID
			songID := sectionToDelete.SongID

			songRepository.On("GetWithSections", new(model.Song), songID).
				Return(nil, &tt.song).
				Once()

			repositoryFactory := new(transaction.RepositoryFactoryMock)
			txSongRepo := new(repository.SongRepositoryMock)
			txSongSectionRepo := new(repository.SongSectionRepositoryMock)
			txSongPartRepo := new(repository.SongPartRepositoryMock)

			repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
			repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepo).Once()
			if tt.withParts {
				repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
			}
			transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

			// Update song (reorder sections)
			txSongRepo.On("UpdateWithAssociations", mock.IsType(&tt.song)).
				Run(func(args mock.Arguments) {
					updatedSong := args.Get(0).(*model.Song)
					// Check reordering
					sections := slices.Clone(updatedSong.Sections)
					sections = slices.DeleteFunc(sections, func(s model.SongSection) bool {
						return s.ID == id
					})
					for i, s := range sections {
						assert.Equal(t, uint(i), s.Order)
					}
				}).
				Return(nil).Once()

			if tt.withParts {
				// Prepare section with parts
				mockSection := &model.SongSection{
					ID:           id,
					SectionParts: []model.SongSectionPart{},
					SongID:       songID,
				}
				for _, pid := range tt.partIDs {
					mockSection.SectionParts = append(mockSection.SectionParts, model.SongSectionPart{PartID: pid})
				}
				txSongSectionRepo.On("GetWithSectionParts", new(model.SongSection), id).
					Return(nil, mockSection).
					Once()

				songProcessor.On("UpdateSongAfterPartsDeletion", txSongRepo, songID, tt.partIDs).
					Return(nil).
					Once()

				if len(tt.partIDs) > 0 {
					txSongPartRepo.On("Delete", tt.partIDs).
						Return(nil).
						Once()
				}
			}

			// Delete section
			txSongSectionRepo.On("Delete", []uuid.UUID{id}).
				Return(nil).
				Once()

			// when
			errCode := _uut.Handle(id, songID, tt.withParts)

			// then
			assert.Nil(t, errCode)

			songRepository.AssertExpectations(t)
			transactionManager.AssertExpectations(t)
			repositoryFactory.AssertExpectations(t)
			txSongRepo.AssertExpectations(t)
			txSongSectionRepo.AssertExpectations(t)
			if tt.withParts {
				txSongPartRepo.AssertExpectations(t)
				songProcessor.AssertExpectations(t)
			}
		})
	}
}
