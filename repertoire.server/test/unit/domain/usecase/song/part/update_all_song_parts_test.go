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
	"repertoire/server/test/unit/domain/validator"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// expectedValidatedBandMemberIDs is the ids argument the use case hands to the band member validator:
// only the request's band member, or nothing at all (nil or empty) when the request has none.
func expectedValidatedBandMemberIDs(request requests.UpdateAllSongPartsRequest) any {
	if request.BandMemberID == nil {
		return mock.MatchedBy(func(ids []uuid.UUID) bool { return len(ids) == 0 })
	}
	return []uuid.UUID{*request.BandMemberID}
}

func TestUpdateAllSongParts_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := part.NewUpdateAllSongParts(songRepository, nil, nil)

	request := requests.UpdateAllSongPartsRequest{SongID: uuid.New()}

	internalError := errors.New("internal error")
	songRepository.On("GetWithParts", new(model.Song), request.SongID).
		Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)
	songRepository.AssertExpectations(t)
}

func TestUpdateAllSongParts_WhenSongNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := part.NewUpdateAllSongParts(songRepository, nil, nil)

	request := requests.UpdateAllSongPartsRequest{SongID: uuid.New()}

	songRepository.On("GetWithParts", new(model.Song), request.SongID).
		Return(nil).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())
	songRepository.AssertExpectations(t)
}

func TestUpdateAllSongParts_WhenBandMemberValidationFails_ShouldReturnItsError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	bandMemberValidator := new(validator.BandMemberValidatorMock)
	_uut := part.NewUpdateAllSongParts(songRepository, bandMemberValidator, nil)

	request := requests.UpdateAllSongPartsRequest{
		SongID:       uuid.New(),
		BandMemberID: &[]uuid.UUID{uuid.New()}[0],
	}

	mockSong := &model.Song{ID: request.SongID}
	songRepository.On("GetWithParts", new(model.Song), request.SongID).
		Return(nil, mockSong).Once()

	validationError := httperror.ConflictError(errors.New("band member is not valid"))
	bandMemberValidator.On("Validate", expectedValidatedBandMemberIDs(request), *mockSong).
		Return(nil, validationError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, validationError, errCode)
	songRepository.AssertExpectations(t)
	bandMemberValidator.AssertExpectations(t)
}

func TestUpdateAllSongParts_WhenTransactionExecuteFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	bandMemberValidator := new(validator.BandMemberValidatorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewUpdateAllSongParts(songRepository, bandMemberValidator, transactionManager)

	request := requests.UpdateAllSongPartsRequest{
		SongID:       uuid.New(),
		InstrumentID: &[]uuid.UUID{uuid.New()}[0],
	}

	mockSong := &model.Song{ID: request.SongID}
	songRepository.On("GetWithParts", new(model.Song), request.SongID).
		Return(nil, mockSong).Once()

	bandMemberValidator.On("Validate", expectedValidatedBandMemberIDs(request), *mockSong).
		Return(nil, nil).Once()

	internalError := errors.New("transaction error")
	transactionManager.On("Execute", mock.Anything).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)
	songRepository.AssertExpectations(t)
	bandMemberValidator.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
}

func TestUpdateAllSongParts_WhenUpdateWithAssociationsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	bandMemberValidator := new(validator.BandMemberValidatorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewUpdateAllSongParts(songRepository, bandMemberValidator, transactionManager)

	request := requests.UpdateAllSongPartsRequest{
		SongID:       uuid.New(),
		InstrumentID: &[]uuid.UUID{uuid.New()}[0],
	}

	mockSong := &model.Song{
		ID: request.SongID,
		Parts: []model.SongPart{
			{ID: uuid.New()},
		},
	}
	songRepository.On("GetWithParts", new(model.Song), request.SongID).
		Return(nil, mockSong).Once()

	bandMemberValidator.On("Validate", expectedValidatedBandMemberIDs(request), *mockSong).
		Return(nil, nil).Once()

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongSectionRepo := new(repository.SongSectionRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	internalError := errors.New("update song error")
	txSongRepo.On("UpdateWithAssociations", mock.IsType(mockSong)).
		Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)
	songRepository.AssertExpectations(t)
	bandMemberValidator.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongSectionRepo.AssertExpectations(t)
}

func TestUpdateAllSongParts_WhenGetAllSectionPartsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	bandMemberValidator := new(validator.BandMemberValidatorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewUpdateAllSongParts(songRepository, bandMemberValidator, transactionManager)

	request := requests.UpdateAllSongPartsRequest{
		SongID:       uuid.New(),
		BandMemberID: &[]uuid.UUID{uuid.New()}[0],
	}

	partID := uuid.New()
	mockSong := &model.Song{
		ID: request.SongID,
		Parts: []model.SongPart{
			{ID: partID},
		},
	}
	songRepository.On("GetWithParts", new(model.Song), request.SongID).
		Return(nil, mockSong).Once()

	bandMembers := []model.BandMember{{ID: *request.BandMemberID}}
	bandMemberValidator.On("Validate", expectedValidatedBandMemberIDs(request), *mockSong).
		Return(bandMembers, nil).Once()

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongSectionRepo := new(repository.SongSectionRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	internalError := errors.New("get section parts error")
	txSongSectionRepo.On("GetAllSectionPartsByPartIDs", new([]model.SongSectionPart), []uuid.UUID{partID}).
		Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)
	songRepository.AssertExpectations(t)
	bandMemberValidator.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongSectionRepo.AssertExpectations(t)
}

func TestUpdateAllSongParts_WhenReplaceSectionPartBandMembersFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	bandMemberValidator := new(validator.BandMemberValidatorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewUpdateAllSongParts(songRepository, bandMemberValidator, transactionManager)

	request := requests.UpdateAllSongPartsRequest{
		SongID:       uuid.New(),
		BandMemberID: &[]uuid.UUID{uuid.New()}[0],
	}

	partID := uuid.New()
	mockSong := &model.Song{
		ID: request.SongID,
		Parts: []model.SongPart{
			{ID: partID},
		},
	}
	songRepository.On("GetWithParts", new(model.Song), request.SongID).
		Return(nil, mockSong).Once()

	bandMembers := []model.BandMember{{ID: *request.BandMemberID}}
	bandMemberValidator.On("Validate", expectedValidatedBandMemberIDs(request), *mockSong).
		Return(bandMembers, nil).Once()

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongSectionRepo := new(repository.SongSectionRepositoryMock)

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	sectionParts := []model.SongSectionPart{
		{PartID: partID, SectionID: uuid.New()},
	}
	txSongSectionRepo.On("GetAllSectionPartsByPartIDs", new([]model.SongSectionPart), []uuid.UUID{partID}).
		Return(nil, &sectionParts).Once()

	internalError := errors.New("replace band members error")
	txSongSectionRepo.On("ReplaceSectionPartBandMembers", mock.IsType(new(model.SongSectionPart)), bandMembers).
		Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)
	songRepository.AssertExpectations(t)
	bandMemberValidator.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongSectionRepo.AssertExpectations(t)
}

func TestUpdateAllSongParts_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	partIDs := []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}

	tests := []struct {
		name         string
		request      requests.UpdateAllSongPartsRequest
		withParts    bool // true → mock GetAllSectionPartsByPartIDs + ReplaceSectionPartBandMembers
		withSongSave bool // true → mock UpdateWithAssociations
	}{
		{
			name: "Only instrument",
			request: requests.UpdateAllSongPartsRequest{
				SongID:       uuid.New(),
				InstrumentID: &[]uuid.UUID{uuid.New()}[0],
			},
			withSongSave: true,
		},
		{
			name: "Only band member",
			request: requests.UpdateAllSongPartsRequest{
				SongID:       uuid.New(),
				BandMemberID: &[]uuid.UUID{uuid.New()}[0],
			},
			withParts: true,
		},
		{
			name: "Both instrument and band member",
			request: requests.UpdateAllSongPartsRequest{
				SongID:       uuid.New(),
				InstrumentID: &[]uuid.UUID{uuid.New()}[0],
				BandMemberID: &[]uuid.UUID{uuid.New()}[0],
			},
			withParts:    true,
			withSongSave: true,
		},
		{
			name: "Neither instrument nor band member (no-op)",
			request: requests.UpdateAllSongPartsRequest{
				SongID: uuid.New(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			songRepository := new(repository.SongRepositoryMock)
			bandMemberValidator := new(validator.BandMemberValidatorMock)
			transactionManager := new(transaction.ManagerMock)
			_uut := part.NewUpdateAllSongParts(songRepository, bandMemberValidator, transactionManager)

			parts := make([]model.SongPart, len(partIDs))
			for i, pid := range partIDs {
				parts[i] = model.SongPart{ID: pid}
			}
			mockSong := &model.Song{ID: tt.request.SongID, Parts: parts}

			songRepository.On("GetWithParts", new(model.Song), tt.request.SongID).
				Return(nil, mockSong).Once()

			// validation
			var bandMembers []model.BandMember
			if tt.request.BandMemberID != nil {
				bandMembers = []model.BandMember{{ID: *tt.request.BandMemberID}}
			}
			bandMemberValidator.On("Validate", expectedValidatedBandMemberIDs(tt.request), *mockSong).
				Return(bandMembers, nil).Once()

			repositoryFactory := new(transaction.RepositoryFactoryMock)
			txSongRepo := new(repository.SongRepositoryMock)
			txSongSectionRepo := new(repository.SongSectionRepositoryMock)

			repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
			repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepo).Once()
			transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

			if tt.withSongSave {
				txSongRepo.On("UpdateWithAssociations", mock.IsType(mockSong)).
					Run(func(args mock.Arguments) {
						updatedSong := args.Get(0).(*model.Song)
						for _, p := range updatedSong.Parts {
							assert.Equal(t, tt.request.InstrumentID, p.InstrumentID)
						}
					}).
					Return(nil).Once()
			}

			var replacedPartIDs []uuid.UUID
			if tt.withParts {
				sectionParts := make([]model.SongSectionPart, len(partIDs))
				for i, pid := range partIDs {
					sectionParts[i] = model.SongSectionPart{PartID: pid, SectionID: uuid.New()}
				}
				txSongSectionRepo.On("GetAllSectionPartsByPartIDs", new([]model.SongSectionPart), partIDs).
					Return(nil, &sectionParts).Once()

				txSongSectionRepo.On("ReplaceSectionPartBandMembers", mock.IsType(new(model.SongSectionPart)), mock.Anything).
					Run(func(args mock.Arguments) {
						sp := args.Get(0).(*model.SongSectionPart)
						replacedPartIDs = append(replacedPartIDs, sp.PartID)

						// every section part ends up with exactly the requested band member
						newBandMembers := args.Get(1).([]model.BandMember)
						require.Len(t, newBandMembers, 1)
						assert.Equal(t, *tt.request.BandMemberID, newBandMembers[0].ID)
					}).
					Return(nil).Times(len(partIDs))
			}

			// when
			errCode := _uut.Handle(tt.request)

			// then
			assert.Nil(t, errCode)
			if tt.withParts {
				assert.ElementsMatch(t, partIDs, replacedPartIDs)
			}

			songRepository.AssertExpectations(t)
			bandMemberValidator.AssertExpectations(t)
			transactionManager.AssertExpectations(t)
			repositoryFactory.AssertExpectations(t)
			txSongRepo.AssertExpectations(t)
			txSongSectionRepo.AssertExpectations(t)
		})
	}
}
