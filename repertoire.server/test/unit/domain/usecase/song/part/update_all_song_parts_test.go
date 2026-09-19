package part

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song/part"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/database/transaction"
	"repertoire/server/test/unit/data/repository"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

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

func TestUpdateAllSongParts_WhenGetBandMemberFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := part.NewUpdateAllSongParts(songRepository, artistRepository, nil)

	memberID := uuid.New()
	request := requests.UpdateAllSongPartsRequest{
		SongID:       uuid.New(),
		BandMemberID: &memberID,
	}

	artistID := uuid.New()
	mockSong := &model.Song{ID: request.SongID, ArtistID: &artistID}
	songRepository.On("GetWithParts", new(model.Song), request.SongID).
		Return(nil, mockSong).Once()

	internalError := errors.New("get band member error")
	artistRepository.On("GetBandMember", new(model.BandMember), memberID).
		Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)
	songRepository.AssertExpectations(t)
	artistRepository.AssertExpectations(t)
}

func TestUpdateAllSongParts_WhenBandMemberNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := part.NewUpdateAllSongParts(songRepository, artistRepository, nil)

	memberID := uuid.New()
	request := requests.UpdateAllSongPartsRequest{
		SongID:       uuid.New(),
		BandMemberID: &memberID,
	}

	artistID := uuid.New()
	mockSong := &model.Song{ID: request.SongID, ArtistID: &artistID}
	songRepository.On("GetWithParts", new(model.Song), request.SongID).
		Return(nil, mockSong).Once()

	artistRepository.On("GetBandMember", new(model.BandMember), memberID).
		Return(nil).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "band member not found", errCode.Error.Error())
	songRepository.AssertExpectations(t)
	artistRepository.AssertExpectations(t)
}

func TestUpdateAllSongParts_WhenSongHasNoArtist_ShouldReturnConflictError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := part.NewUpdateAllSongParts(songRepository, artistRepository, nil)

	memberID := uuid.New()
	request := requests.UpdateAllSongPartsRequest{
		SongID:       uuid.New(),
		BandMemberID: &memberID,
	}

	mockSong := &model.Song{ID: request.SongID, ArtistID: nil}
	songRepository.On("GetWithParts", new(model.Song), request.SongID).
		Return(nil, mockSong).Once()

	artistRepository.On("GetBandMember", new(model.BandMember), memberID).
		Return(nil, &model.BandMember{ID: memberID}).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusConflict, errCode.Code)
	assert.Equal(t, "band member is not part of the artist associated with this song", errCode.Error.Error())
	songRepository.AssertExpectations(t)
	artistRepository.AssertExpectations(t)
}

func TestUpdateAllSongParts_WhenBandMemberNotAssociated_ShouldReturnConflictError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := part.NewUpdateAllSongParts(songRepository, artistRepository, nil)

	memberID := uuid.New()
	request := requests.UpdateAllSongPartsRequest{
		SongID:       uuid.New(),
		BandMemberID: &memberID,
	}

	artistID := uuid.New()
	mockSong := &model.Song{ID: request.SongID, ArtistID: &artistID}
	songRepository.On("GetWithParts", new(model.Song), request.SongID).
		Return(nil, mockSong).Once()

	differentArtistID := uuid.New()
	artistRepository.On("GetBandMember", new(model.BandMember), memberID).
		Return(nil, &model.BandMember{ID: memberID, ArtistID: differentArtistID}).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusConflict, errCode.Code)
	assert.Equal(t, "band member is not part of the artist associated with this song", errCode.Error.Error())
	songRepository.AssertExpectations(t)
	artistRepository.AssertExpectations(t)
}

func TestUpdateAllSongParts_WhenTransactionExecuteFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewUpdateAllSongParts(songRepository, nil, transactionManager)

	request := requests.UpdateAllSongPartsRequest{
		SongID:       uuid.New(),
		InstrumentID: &[]uuid.UUID{uuid.New()}[0],
	}

	mockSong := &model.Song{ID: request.SongID}
	songRepository.On("GetWithParts", new(model.Song), request.SongID).
		Return(nil, mockSong).Once()

	internalError := errors.New("transaction error")
	transactionManager.On("Execute", mock.Anything).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)
	songRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
}

func TestUpdateAllSongParts_WhenUpdateWithAssociationsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewUpdateAllSongParts(songRepository, nil, transactionManager)

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
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongSectionRepo.AssertExpectations(t)
}

func TestUpdateAllSongParts_WhenGetAllSectionPartsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	artistRepository := new(repository.ArtistRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewUpdateAllSongParts(songRepository, artistRepository, transactionManager)

	memberID := uuid.New()
	request := requests.UpdateAllSongPartsRequest{
		SongID:       uuid.New(),
		BandMemberID: &memberID,
	}

	artistID := uuid.New()
	partID := uuid.New()
	mockSong := &model.Song{
		ID:       request.SongID,
		ArtistID: &artistID,
		Parts: []model.SongPart{
			{ID: partID},
		},
	}
	songRepository.On("GetWithParts", new(model.Song), request.SongID).
		Return(nil, mockSong).Once()
	artistRepository.On("GetBandMember", new(model.BandMember), memberID).
		Return(nil, &model.BandMember{ID: memberID, ArtistID: artistID}).Once()

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
	artistRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongSectionRepo.AssertExpectations(t)
}

func TestUpdateAllSongParts_WhenUpdateAllSectionPartsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	artistRepository := new(repository.ArtistRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewUpdateAllSongParts(songRepository, artistRepository, transactionManager)

	memberID := uuid.New()
	request := requests.UpdateAllSongPartsRequest{
		SongID:       uuid.New(),
		BandMemberID: &memberID,
	}

	artistID := uuid.New()
	partID := uuid.New()
	mockSong := &model.Song{
		ID:       request.SongID,
		ArtistID: &artistID,
		Parts: []model.SongPart{
			{ID: partID},
		},
	}
	songRepository.On("GetWithParts", new(model.Song), request.SongID).
		Return(nil, mockSong).Once()
	artistRepository.On("GetBandMember", new(model.BandMember), memberID).
		Return(nil, &model.BandMember{ID: memberID, ArtistID: artistID}).Once()

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

	internalError := errors.New("update section parts error")
	txSongSectionRepo.On("UpdateAllSectionParts", mock.IsType(&sectionParts)).
		Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)
	songRepository.AssertExpectations(t)
	artistRepository.AssertExpectations(t)
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
		member       *model.BandMember // only when BandMemberID is set
		withParts    bool              // true → mock GetAllSectionPartsByPartIDs + UpdateAllSectionParts
		withSongSave bool              // true → mock UpdateWithAssociations
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
			member:    &model.BandMember{ArtistID: uuid.New()},
			withParts: true,
		},
		{
			name: "Both instrument and band member",
			request: requests.UpdateAllSongPartsRequest{
				SongID:       uuid.New(),
				InstrumentID: &[]uuid.UUID{uuid.New()}[0],
				BandMemberID: &[]uuid.UUID{uuid.New()}[0],
			},
			member:       &model.BandMember{ArtistID: uuid.New()},
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
			artistRepository := new(repository.ArtistRepositoryMock)
			transactionManager := new(transaction.ManagerMock)
			_uut := part.NewUpdateAllSongParts(songRepository, artistRepository, transactionManager)

			parts := make([]model.SongPart, len(partIDs))
			for i, pid := range partIDs {
				parts[i] = model.SongPart{ID: pid}
			}
			mockSong := &model.Song{ID: tt.request.SongID, Parts: parts}
			if tt.request.BandMemberID != nil {
				mockSong.ArtistID = &tt.member.ArtistID
			}

			songRepository.On("GetWithParts", new(model.Song), tt.request.SongID).
				Return(nil, mockSong).Once()

			if tt.request.BandMemberID != nil {
				artistRepository.On("GetBandMember", new(model.BandMember), *tt.request.BandMemberID).
					Return(nil, tt.member).Once()
			}

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

			if tt.withParts {
				sectionParts := make([]model.SongSectionPart, len(partIDs))
				for i, pid := range partIDs {
					sectionParts[i] = model.SongSectionPart{PartID: pid, SectionID: uuid.New()}
				}
				txSongSectionRepo.On("GetAllSectionPartsByPartIDs", new([]model.SongSectionPart), partIDs).
					Return(nil, &sectionParts).Once()

				txSongSectionRepo.On("UpdateAllSectionParts", mock.IsType(&sectionParts)).
					Run(func(args mock.Arguments) {
						updated := args.Get(0).(*[]model.SongSectionPart)
						for _, sp := range *updated {
							assert.Equal(t, tt.request.BandMemberID, sp.BandMemberID)
						}
					}).
					Return(nil).Once()
			}

			// when
			errCode := _uut.Handle(tt.request)

			// then
			assert.Nil(t, errCode)

			songRepository.AssertExpectations(t)
			artistRepository.AssertExpectations(t)
			transactionManager.AssertExpectations(t)
			repositoryFactory.AssertExpectations(t)
			txSongRepo.AssertExpectations(t)
			txSongSectionRepo.AssertExpectations(t)
		})
	}
}
