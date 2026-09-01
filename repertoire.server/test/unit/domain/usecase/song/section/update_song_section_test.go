package section

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song/section"
	"repertoire/server/internal/deduplicate"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/database/transaction"
	"repertoire/server/test/unit/data/repository"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateSongSection_WhenGetSectionFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	_uut := section.NewUpdateSongSection(
		songSectionRepository,
		nil,
		nil,
		nil,
	)

	request := requests.UpdateSongSectionRequest{
		ID: uuid.New(),
	}

	internalError := errors.New("get error")
	songSectionRepository.On("GetWithSectionParts", new(model.SongSection), request.ID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songSectionRepository.AssertExpectations(t)
}

func TestUpdateSongSection_WhenSectionNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	_uut := section.NewUpdateSongSection(
		songSectionRepository,
		nil,
		nil,
		nil,
	)

	request := requests.UpdateSongSectionRequest{
		ID: uuid.New(),
	}

	songSectionRepository.On("GetWithSectionParts", new(model.SongSection), request.ID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song section not found", errCode.Error.Error())

	songSectionRepository.AssertExpectations(t)
}

func TestUpdateSongSection_WhenGetPartsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songPartRepository := new(repository.SongPartRepositoryMock)
	_uut := section.NewUpdateSongSection(
		songSectionRepository,
		songPartRepository,
		nil,
		nil,
	)

	sectionID := uuid.New()
	partID := uuid.New()
	request := requests.UpdateSongSectionRequest{
		ID:      sectionID,
		PartIDs: []uuid.UUID{partID},
	}

	mockSection := &model.SongSection{ID: sectionID}
	songSectionRepository.On("GetWithSectionParts", new(model.SongSection), request.ID).
		Return(nil, mockSection).
		Once()

	internalError := errors.New("get parts error")
	songPartRepository.On("GetAllByIDs", new([]model.SongPart), request.PartIDs).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songSectionRepository.AssertExpectations(t)
	songPartRepository.AssertExpectations(t)
}

func TestUpdateSongSection_WhenSomePartsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songPartRepository := new(repository.SongPartRepositoryMock)
	_uut := section.NewUpdateSongSection(
		songSectionRepository,
		songPartRepository,
		nil,
		nil,
	)

	sectionID := uuid.New()
	partID := uuid.New()
	request := requests.UpdateSongSectionRequest{
		ID:      sectionID,
		PartIDs: []uuid.UUID{partID, uuid.New()},
	}

	mockSection := &model.SongSection{ID: sectionID}
	songSectionRepository.On("GetWithSectionParts", new(model.SongSection), request.ID).
		Return(nil, mockSection).
		Once()

	// Return fewer parts than requested
	mockParts := []model.SongPart{{ID: partID}}
	songPartRepository.On("GetAllByIDs", new([]model.SongPart), request.PartIDs).
		Return(nil, &mockParts).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "some parts not found", errCode.Error.Error())

	songSectionRepository.AssertExpectations(t)
	songPartRepository.AssertExpectations(t)
}

func TestUpdateSongSection_WhenPartDoesNotBelongToSong_ShouldReturnConflictError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songPartRepository := new(repository.SongPartRepositoryMock)
	_uut := section.NewUpdateSongSection(
		songSectionRepository,
		songPartRepository,
		nil,
		nil,
	)

	sectionID := uuid.New()
	songID := uuid.New()
	partID := uuid.New()
	request := requests.UpdateSongSectionRequest{
		ID:      sectionID,
		PartIDs: []uuid.UUID{partID},
	}

	mockSection := &model.SongSection{ID: sectionID, SongID: songID}
	songSectionRepository.On("GetWithSectionParts", new(model.SongSection), request.ID).
		Return(nil, mockSection).
		Once()

	mockParts := []model.SongPart{{ID: partID, SongID: uuid.New()}}
	songPartRepository.On("GetAllByIDs", new([]model.SongPart), request.PartIDs).
		Return(nil, &mockParts).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusConflict, errCode.Code)
	assert.Equal(t, "song part does not belong to the same song as the section", errCode.Error.Error())

	songSectionRepository.AssertExpectations(t)
	songPartRepository.AssertExpectations(t)
}

func TestUpdateSongSection_WhenTransactionExecuteFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songPartRepository := new(repository.SongPartRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := section.NewUpdateSongSection(
		songSectionRepository,
		songPartRepository,
		nil,
		transactionManager,
	)

	sectionID := uuid.New()
	songID := uuid.New()
	partID := uuid.New()
	request := requests.UpdateSongSectionRequest{
		ID:      sectionID,
		PartIDs: []uuid.UUID{partID},
	}

	mockSection := &model.SongSection{ID: sectionID, SongID: songID}
	songSectionRepository.On("GetWithSectionParts", new(model.SongSection), request.ID).
		Return(nil, mockSection).
		Once()

	mockParts := []model.SongPart{{ID: partID, SongID: songID}}
	songPartRepository.On("GetAllByIDs", new([]model.SongPart), request.PartIDs).
		Return(nil, &mockParts).
		Once()

	internalError := errors.New("tx error")
	transactionManager.On("Execute", mock.Anything).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	transactionManager.AssertExpectations(t)
}

func TestUpdateSongSection_WhenUpdateSectionFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := section.NewUpdateSongSection(songSectionRepository, nil, nil, transactionManager)

	sectionID := uuid.New()
	request := requests.UpdateSongSectionRequest{
		ID:      sectionID,
		Name:    "New Name",
		TypeID:  uuid.New(),
		PartIDs: []uuid.UUID{},
	}

	mockSection := &model.SongSection{
		ID:           sectionID,
		SectionParts: []model.SongSectionPart{},
	}
	songSectionRepository.On("GetWithSectionParts", new(model.SongSection), request.ID).Return(nil, mockSection).
		Once()

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongSectionRepository := new(repository.SongSectionRepositoryMock)

	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	internalError := errors.New("update section error")
	txSongSectionRepository.On("Update", mock.IsType(mockSection)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songSectionRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongSectionRepository.AssertExpectations(t)
}

func TestUpdateSongSection_WhenDeleteSectionPartsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := section.NewUpdateSongSection(
		songSectionRepository,
		nil,
		nil,
		transactionManager,
	)

	sectionID := uuid.New()
	partID := uuid.New()
	request := requests.UpdateSongSectionRequest{
		ID:      sectionID,
		PartIDs: []uuid.UUID{}, // empty -> all existing should be deleted
	}

	mockSection := &model.SongSection{
		ID: sectionID,
		SectionParts: []model.SongSectionPart{
			{PartID: partID, Order: 0},
		},
	}
	songSectionRepository.On("GetWithSectionParts", new(model.SongSection), request.ID).
		Return(nil, mockSection).
		Once()

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongSectionRepository := new(repository.SongSectionRepositoryMock)

	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txSongSectionRepository.On("Update", mock.IsType(mockSection)).
		Return(nil).
		Once()

	internalError := errors.New("delete error")
	txSongSectionRepository.On("DeleteSectionParts", &mockSection.SectionParts).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songSectionRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongSectionRepository.AssertExpectations(t)
}

func TestUpdateSongSection_WhenUpdateAllSectionPartsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songPartRepository := new(repository.SongPartRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := section.NewUpdateSongSection(
		songSectionRepository,
		songPartRepository,
		nil,
		transactionManager,
	)

	sectionID := uuid.New()
	partID := uuid.New()
	songID := uuid.New()
	request := requests.UpdateSongSectionRequest{
		ID:      sectionID,
		PartIDs: []uuid.UUID{partID}, // keep the same part
	}

	mockSection := &model.SongSection{
		ID: sectionID,
		SectionParts: []model.SongSectionPart{
			{PartID: partID, Order: 0},
		},
		SongID: songID,
	}
	songSectionRepository.On("GetWithSectionParts", new(model.SongSection), request.ID).
		Return(nil, mockSection).
		Once()

	mockParts := &[]model.SongPart{
		{ID: partID, SongID: songID},
	}
	songPartRepository.On("GetAllByIDs", new([]model.SongPart), request.PartIDs).
		Return(nil, mockParts).
		Once()

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongSectionRepository := new(repository.SongSectionRepositoryMock)

	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txSongSectionRepository.On("Update", mock.IsType(mockSection)).
		Return(nil).
		Once()

	internalError := errors.New("update error")
	txSongSectionRepository.On("UpdateAllSectionParts", mock.IsType(&[]model.SongSectionPart{})).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songSectionRepository.AssertExpectations(t)
	songPartRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongSectionRepository.AssertExpectations(t)
}

func TestUpdateSongSection_WhenCreateAllSectionPartsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songPartRepository := new(repository.SongPartRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := section.NewUpdateSongSection(
		songSectionRepository,
		songPartRepository,
		nil,
		transactionManager,
	)

	sectionID := uuid.New()
	partID := uuid.New()
	songID := uuid.New()
	request := requests.UpdateSongSectionRequest{
		ID:      sectionID,
		PartIDs: []uuid.UUID{partID},
	}

	mockSection := &model.SongSection{
		ID:           sectionID,
		SectionParts: []model.SongSectionPart{}, // no existing parts
		SongID:       songID,
	}
	songSectionRepository.On("GetWithSectionParts", new(model.SongSection), request.ID).
		Return(nil, mockSection).
		Once()

	mockParts := &[]model.SongPart{
		{ID: partID, SongID: songID},
	}
	songPartRepository.On("GetAllByIDs", new([]model.SongPart), request.PartIDs).
		Return(nil, mockParts).
		Once()

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongSectionRepository := new(repository.SongSectionRepositoryMock)

	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txSongSectionRepository.On("Update", mock.IsType(mockSection)).
		Return(nil).
		Once()

	internalError := errors.New("create error")
	txSongSectionRepository.On("CreateAllSectionParts", mock.IsType(&[]model.SongSectionPart{})).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songSectionRepository.AssertExpectations(t)
	songPartRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongSectionRepository.AssertExpectations(t)
}

func TestUpdateSongSection_WhenSuccessful_ShouldUpdateSectionAndParts(t *testing.T) {
	tests := []struct {
		name               string
		section            model.SongSection
		request            requests.UpdateSongSectionRequest
		partsToKeepIndices []int
	}{
		{
			name: "Update name and type only (keep parts)",
			section: model.SongSection{
				ID:                uuid.New(),
				Name:              "Old Name",
				SongSectionTypeID: uuid.New(),
				SongID:            uuid.New(),
				SectionParts: []model.SongSectionPart{
					{PartID: uuid.New(), Order: 0},
					{PartID: uuid.New(), Order: 1},
				},
			},
			request: requests.UpdateSongSectionRequest{
				ID:      uuid.New(),
				Name:    "New Name",
				TypeID:  uuid.New(),
				PartIDs: []uuid.UUID{uuid.New(), uuid.New()},
			},
			partsToKeepIndices: []int{0, 1},
		},
		{
			name: "Add 2 new parts",
			section: model.SongSection{
				ID:                uuid.New(),
				Name:              "Old Name",
				SongSectionTypeID: uuid.New(),
				SongID:            uuid.New(),
				SectionParts:      []model.SongSectionPart{},
			},
			request: requests.UpdateSongSectionRequest{
				ID:      uuid.New(),
				Name:    "New Name",
				TypeID:  uuid.New(),
				PartIDs: []uuid.UUID{uuid.New(), uuid.New()},
			},
		},
		{
			name: "Remove parts",
			section: model.SongSection{
				ID:                uuid.New(),
				Name:              "Old Name",
				SongSectionTypeID: uuid.New(),
				SongID:            uuid.New(),
				SectionParts: []model.SongSectionPart{
					{PartID: uuid.New(), Order: 0},
					{PartID: uuid.New(), Order: 1},
					{PartID: uuid.New(), Order: 2},
				},
			},
			request: requests.UpdateSongSectionRequest{
				ID:      uuid.New(),
				Name:    "New Name",
				TypeID:  uuid.New(),
				PartIDs: []uuid.UUID{},
			},
		},
		{
			name: "Remove/Add/Keep parts",
			section: model.SongSection{
				ID:                uuid.New(),
				Name:              "Old Name",
				SongSectionTypeID: uuid.New(),
				SongID:            uuid.New(),
				SectionParts: []model.SongSectionPart{
					{PartID: uuid.New(), Order: 0},
					{PartID: uuid.New(), Order: 1},
					{PartID: uuid.New(), Order: 2},
				},
			},
			request: requests.UpdateSongSectionRequest{
				ID:      uuid.New(),
				Name:    "New Name",
				TypeID:  uuid.New(),
				PartIDs: []uuid.UUID{uuid.New(), uuid.New(), uuid.New()},
			},
			partsToKeepIndices: []int{0, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			songSectionRepository := new(repository.SongSectionRepositoryMock)
			songPartRepository := new(repository.SongPartRepositoryMock)
			transactionManager := new(transaction.ManagerMock)
			_uut := section.NewUpdateSongSection(
				songSectionRepository,
				songPartRepository,
				nil,
				transactionManager,
			)

			// prepare
			for i, idx := range tt.partsToKeepIndices {
				tt.request.PartIDs[i] = tt.section.SectionParts[idx].PartID
			}

			oldMap := make(map[uuid.UUID]model.SongSectionPart)
			for _, sp := range tt.section.SectionParts {
				oldMap[sp.PartID] = sp
			}
			partIDs := deduplicate.Deduplicate(tt.request.PartIDs)

			var expectedPartsToDelete []model.SongSectionPart
			for pid, sp := range oldMap {
				if !partIDs[pid] {
					expectedPartsToDelete = append(expectedPartsToDelete, sp)
				}
			}

			var expectedPartsToUpdate []model.SongSectionPart
			var expectedPartsToCreate []model.SongSectionPart
			order := 0
			for pid := range partIDs {
				if sp, exists := oldMap[pid]; exists {
					sp.Order = uint(order)
					expectedPartsToUpdate = append(expectedPartsToUpdate, sp)
				} else {
					expectedPartsToCreate = append(expectedPartsToCreate, model.SongSectionPart{
						PartID:    pid,
						SectionID: tt.section.ID,
						Order:     uint(order),
					})
				}
				order++
			}

			// given - mocking
			songSectionRepository.On("GetWithSectionParts", new(model.SongSection), tt.request.ID).
				Return(nil, &tt.section).
				Once()

			if len(tt.request.PartIDs) > 0 {
				mockParts := make([]model.SongPart, len(tt.request.PartIDs))
				for i, pid := range tt.request.PartIDs {
					mockParts[i] = model.SongPart{ID: pid, SongID: tt.section.SongID}
				}
				songPartRepository.On("GetAllByIDs", new([]model.SongPart), tt.request.PartIDs).
					Return(nil, &mockParts).
					Once()
			}

			repositoryFactory := new(transaction.RepositoryFactoryMock)
			txSongSectionRepository := new(repository.SongSectionRepositoryMock)

			repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepository).Once()
			transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

			txSongSectionRepository.On("Update", mock.IsType(&tt.section)).
				Run(func(args mock.Arguments) {
					updatedSection := args.Get(0).(*model.SongSection)
					assert.Equal(t, tt.request.Name, updatedSection.Name)
					assert.Equal(t, tt.request.TypeID, updatedSection.SongSectionTypeID)
				}).
				Return(nil).
				Once()

			if len(expectedPartsToDelete) > 0 {
				txSongSectionRepository.On("DeleteSectionParts", &expectedPartsToDelete).
					Return(nil).
					Once()
			}
			if len(expectedPartsToUpdate) > 0 {
				txSongSectionRepository.On("UpdateAllSectionParts", mock.IsType(&expectedPartsToUpdate)).
					Return(nil).
					Once()
			}
			if len(expectedPartsToCreate) > 0 {
				txSongSectionRepository.On("CreateAllSectionParts", mock.IsType(&expectedPartsToCreate)).
					Return(nil).
					Once()
			}

			// when
			errCode := _uut.Handle(tt.request)

			// then
			assert.Nil(t, errCode)

			songSectionRepository.AssertExpectations(t)
			songPartRepository.AssertExpectations(t)
			transactionManager.AssertExpectations(t)
			repositoryFactory.AssertExpectations(t)
			txSongSectionRepository.AssertExpectations(t)
		})
	}
}
