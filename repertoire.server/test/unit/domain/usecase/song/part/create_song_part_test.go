package part

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
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

func TestCreateSongPart_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := part.NewCreateSongPart(nil, songRepository, nil, nil)

	request := requests.CreateSongPartRequest{
		SongID: uuid.New(),
	}

	internalError := errors.New("internal error")
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)
	songRepository.AssertExpectations(t)
}

func TestCreateSongPart_WhenSongNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := part.NewCreateSongPart(nil, songRepository, nil, nil)

	request := requests.CreateSongPartRequest{
		SongID: uuid.New(),
	}

	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())
	songRepository.AssertExpectations(t)
}

func TestCreateSongPart_WhenGetSectionFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := part.NewCreateSongPart(songSectionRepository, songRepository, nil, nil)

	request := requests.CreateSongPartRequest{
		SongID:    uuid.New(),
		SectionID: &[]uuid.UUID{uuid.New()}[0],
	}

	mockSong := &model.Song{ID: uuid.New()}
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	internalError := errors.New("internal error")
	songSectionRepository.On("Get", new(model.SongSection), *request.SectionID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songSectionRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestCreateSongPart_WhenSectionNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := part.NewCreateSongPart(songSectionRepository, songRepository, nil, nil)

	request := requests.CreateSongPartRequest{
		SongID:    uuid.New(),
		SectionID: &[]uuid.UUID{uuid.New()}[0],
	}

	mockSong := &model.Song{ID: uuid.New()}
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	songSectionRepository.On("Get", new(model.SongSection), *request.SectionID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "section not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
	songSectionRepository.AssertExpectations(t)
}

func TestCreateSongPart_WhenSectionBelongsToDifferentSong_ShouldReturnConflictError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := part.NewCreateSongPart(songSectionRepository, songRepository, nil, nil)

	request := requests.CreateSongPartRequest{
		SongID:    uuid.New(),
		SectionID: &[]uuid.UUID{uuid.New()}[0],
	}

	mockSong := &model.Song{ID: uuid.New()}
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	mockSection := &model.SongSection{
		ID:     *request.SectionID,
		SongID: uuid.New(), // different song
	}
	songSectionRepository.On("Get", new(model.SongSection), *request.SectionID).
		Return(nil, mockSection).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusConflict, errCode.Code)
	assert.Equal(t, "section does not belong to the same song", errCode.Error.Error())
	songSectionRepository.AssertExpectations(t)
}

func TestCreateSongPart_WhenGetBandMemberFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := part.NewCreateSongPart(nil, songRepository, artistRepository, nil)

	request := requests.CreateSongPartRequest{
		SongID:       uuid.New(),
		Name:         "Some Part",
		BandMemberID: &[]uuid.UUID{uuid.New()}[0],
	}

	artistID := uuid.New()
	mockSong := &model.Song{ID: uuid.New(), ArtistID: &artistID}
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	internalError := errors.New("internal error")
	artistRepository.On("GetBandMember", new(model.BandMember), *request.BandMemberID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
	artistRepository.AssertExpectations(t)
}

func TestCreateSongPart_WhenBandMemberNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := part.NewCreateSongPart(nil, songRepository, artistRepository, nil)

	request := requests.CreateSongPartRequest{
		SongID:       uuid.New(),
		Name:         "Some Part",
		BandMemberID: &[]uuid.UUID{uuid.New()}[0],
	}

	mockSong := &model.Song{ID: uuid.New(), ArtistID: &[]uuid.UUID{uuid.New()}[0]}
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	artistRepository.On("GetBandMember", new(model.BandMember), *request.BandMemberID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "band member not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
	artistRepository.AssertExpectations(t)
}

func TestCreateSongPart_WhenBandMemberNotAssociated_ShouldReturnConflictError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := part.NewCreateSongPart(nil, songRepository, artistRepository, nil)

	request := requests.CreateSongPartRequest{
		SongID:       uuid.New(),
		Name:         "Some Part",
		BandMemberID: &[]uuid.UUID{uuid.New()}[0],
	}

	artistID := uuid.New()
	mockSong := &model.Song{ID: uuid.New(), ArtistID: &artistID}
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	bandMember := &model.BandMember{ID: *request.BandMemberID, ArtistID: uuid.New()}
	artistRepository.On("GetBandMember", new(model.BandMember), *request.BandMemberID).
		Return(nil, bandMember).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusConflict, errCode.Code)
	assert.Equal(t, "band member is not part of the artist associated with this song", errCode.Error.Error())

	songRepository.AssertExpectations(t)
	artistRepository.AssertExpectations(t)
}

func TestCreateSongPart_WhenTransactionExecuteFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewCreateSongPart(nil, songRepository, nil, transactionManager)

	request := requests.CreateSongPartRequest{
		SongID: uuid.New(),
		Name:   "Some Part",
	}

	mockSong := &model.Song{ID: uuid.New()}
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	internalError := errors.New("transaction failed")
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

func TestCreateSongPart_WhenCountAllBySongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewCreateSongPart(nil, songRepository, nil, transactionManager)

	request := requests.CreateSongPartRequest{
		SongID: uuid.New(),
		Name:   "Some Part",
	}

	mockSong := &model.Song{ID: uuid.New()}
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongSectionRepo := new(repository.SongSectionRepositoryMock)
	txSongArrangementRepo := new(repository.SongArrangementRepositoryMock)

	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepo).Once()
	repositoryFactory.On("NewSongArrangementRepository").Return(txSongArrangementRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	internalError := errors.New("count error")
	txSongPartRepo.On("CountAllBySong", new(int64), request.SongID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongPartRepo.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongSectionRepo.AssertExpectations(t)
	txSongArrangementRepo.AssertExpectations(t)
}

func TestCreateSongPart_WhenCountAllSectionPartsBySectionIDFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewCreateSongPart(songSectionRepository, songRepository, nil, transactionManager)

	request := requests.CreateSongPartRequest{
		SongID:    uuid.New(),
		Name:      "Some Part",
		SectionID: &[]uuid.UUID{uuid.New()}[0],
	}

	mockSong := &model.Song{ID: request.SongID}
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	// section validation
	mockSection := &model.SongSection{ID: *request.SectionID, SongID: request.SongID}
	songSectionRepository.On("Get", new(model.SongSection), *request.SectionID).
		Return(nil, mockSection).
		Once()

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)
	txSongSectionRepo := new(repository.SongSectionRepositoryMock)
	txSongArrangementRepo := new(repository.SongArrangementRepositoryMock)

	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepo).Once()
	repositoryFactory.On("NewSongArrangementRepository").Return(txSongArrangementRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txSongPartRepo.On("CountAllBySong", new(int64), request.SongID).
		Return(nil, &[]int64{0}[0]).
		Once()

	internalError := errors.New("count section parts error")
	txSongSectionRepo.On("CountAllSectionPartsBySectionID", new(int64), *request.SectionID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songSectionRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongPartRepo.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongSectionRepo.AssertExpectations(t)
	txSongArrangementRepo.AssertExpectations(t)
}

func TestCreateSongPart_WhenCreatePartFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewCreateSongPart(nil, songRepository, nil, transactionManager)

	request := requests.CreateSongPartRequest{
		SongID: uuid.New(),
		Name:   "Some Part",
	}

	mockSong := &model.Song{ID: uuid.New()}
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongSectionRepo := new(repository.SongSectionRepositoryMock)
	txSongArrangementRepo := new(repository.SongArrangementRepositoryMock)

	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepo).Once()
	repositoryFactory.On("NewSongArrangementRepository").Return(txSongArrangementRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txSongPartRepo.On("CountAllBySong", new(int64), request.SongID).
		Return(nil, &[]int64{0}[0]).
		Once()

	internalError := errors.New("create error")
	txSongPartRepo.On("Create", mock.IsType(new(model.SongPart))).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongPartRepo.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongSectionRepo.AssertExpectations(t)
	txSongArrangementRepo.AssertExpectations(t)
}

func TestCreateSongPart_WhenUpdateSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewCreateSongPart(nil, songRepository, nil, transactionManager)

	request := requests.CreateSongPartRequest{
		SongID: uuid.New(),
		Name:   "Some Part",
	}

	mockSong := &model.Song{ID: uuid.New()}
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongSectionRepo := new(repository.SongSectionRepositoryMock)
	txSongArrangementRepo := new(repository.SongArrangementRepositoryMock)

	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepo).Once()
	repositoryFactory.On("NewSongArrangementRepository").Return(txSongArrangementRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txSongPartRepo.On("CountAllBySong", new(int64), request.SongID).
		Return(nil, &[]int64{2}[0]).
		Once()

	txSongPartRepo.On("Create", mock.IsType(new(model.SongPart))).
		Return(nil).
		Once()

	internalError := errors.New("update error")
	txSongRepo.On("Update", mock.IsType(mockSong)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongPartRepo.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongSectionRepo.AssertExpectations(t)
	txSongArrangementRepo.AssertExpectations(t)
}

func TestCreateSongPart_WhenGetArrangementsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewCreateSongPart(nil, songRepository, nil, transactionManager)

	request := requests.CreateSongPartRequest{
		SongID: uuid.New(),
		Name:   "Some Part",
	}

	mockSong := &model.Song{ID: uuid.New()}
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongSectionRepo := new(repository.SongSectionRepositoryMock)
	txSongArrangementRepo := new(repository.SongArrangementRepositoryMock)

	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepo).Once()
	repositoryFactory.On("NewSongArrangementRepository").Return(txSongArrangementRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txSongPartRepo.On("CountAllBySong", new(int64), request.SongID).
		Return(nil, &[]int64{0}[0]).
		Once()

	txSongPartRepo.On("Create", mock.IsType(new(model.SongPart))).
		Return(nil).
		Once()

	txSongRepo.On("Update", mock.IsType(mockSong)).
		Return(nil).
		Once()

	internalError := errors.New("get arrangements error")
	txSongArrangementRepo.On("GetAllBySong", new([]model.SongArrangement), request.SongID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongPartRepo.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongSectionRepo.AssertExpectations(t)
	txSongArrangementRepo.AssertExpectations(t)
}

func TestCreateSongPart_WhenUpdateArrangementsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewCreateSongPart(nil, songRepository, nil, transactionManager)

	request := requests.CreateSongPartRequest{
		SongID: uuid.New(),
		Name:   "Some Part",
	}

	mockSong := &model.Song{ID: uuid.New()}
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongSectionRepo := new(repository.SongSectionRepositoryMock)
	txSongArrangementRepo := new(repository.SongArrangementRepositoryMock)

	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepo).Once()
	repositoryFactory.On("NewSongArrangementRepository").Return(txSongArrangementRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txSongPartRepo.On("CountAllBySong", new(int64), request.SongID).
		Return(nil, &[]int64{0}[0]).
		Once()

	txSongPartRepo.On("Create", mock.IsType(new(model.SongPart))).
		Return(nil).
		Once()

	txSongRepo.On("Update", mock.IsType(mockSong)).
		Return(nil).
		Once()

	arrangements := []model.SongArrangement{{ID: uuid.New()}}
	txSongArrangementRepo.On("GetAllBySong", new([]model.SongArrangement), request.SongID).
		Return(nil, &arrangements).
		Once()

	internalError := errors.New("update arrangements error")
	txSongArrangementRepo.On("UpdateAllWithAssociations", mock.IsType(&arrangements)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongPartRepo.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongSectionRepo.AssertExpectations(t)
	txSongArrangementRepo.AssertExpectations(t)
}

func TestCreateSongPart_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	tests := []struct {
		name                   string
		request                requests.CreateSongPartRequest
		song                   model.Song
		partsCount             int64
		sectionPartsCount      int64
		expectedSongConfidence float64
		expectedSongRehearsals float64
		expectedSongProgress   float64
	}{
		{
			name: "No prior parts, no section, no band member",
			request: requests.CreateSongPartRequest{
				SongID: uuid.New(),
				Name:   "Part 1",
			},
			song:                   model.Song{ID: uuid.New()},
			partsCount:             0,
			expectedSongConfidence: 0,
			expectedSongRehearsals: 0,
			expectedSongProgress:   0,
		},
		{
			name: "Prior parts with stats",
			request: requests.CreateSongPartRequest{
				SongID: uuid.New(),
				Name:   "Part X",
			},
			song: model.Song{
				ID:         uuid.New(),
				Confidence: 50,
				Rehearsals: 10,
				Progress:   54,
			},
			partsCount:             1,
			expectedSongConfidence: 25,
			expectedSongRehearsals: 5,
			expectedSongProgress:   27,
		},
		{
			name: "With one section, no existing parts",
			request: requests.CreateSongPartRequest{
				SongID:    uuid.New(),
				Name:      "Part with section",
				SectionID: &[]uuid.UUID{uuid.New()}[0],
			},
			song: model.Song{
				ID:         uuid.New(),
				Confidence: 0,
				Rehearsals: 0,
				Progress:   0,
			},
			partsCount:             0,
			sectionPartsCount:      0,
			expectedSongConfidence: 0,
			expectedSongRehearsals: 0,
			expectedSongProgress:   0,
		},
		{
			name: "With one section that already has parts",
			request: requests.CreateSongPartRequest{
				SongID:    uuid.New(),
				Name:      "Part in existing section",
				SectionID: &[]uuid.UUID{uuid.New()}[0],
			},
			song: model.Song{
				ID:         uuid.New(),
				Confidence: 0,
				Rehearsals: 0,
				Progress:   0,
			},
			partsCount:             0,
			sectionPartsCount:      3,
			expectedSongConfidence: 0,
			expectedSongRehearsals: 0,
			expectedSongProgress:   0,
		},
		{
			name: "With section and band member",
			request: requests.CreateSongPartRequest{
				SongID:       uuid.New(),
				Name:         "Part with band member",
				SectionID:    &[]uuid.UUID{uuid.New()}[0],
				BandMemberID: &[]uuid.UUID{uuid.New()}[0],
			},
			song: model.Song{
				ID:         uuid.New(),
				Confidence: 0,
				Rehearsals: 0,
				Progress:   0,
				ArtistID:   &[]uuid.UUID{uuid.New()}[0],
			},
			partsCount:             0,
			sectionPartsCount:      2,
			expectedSongConfidence: 0,
			expectedSongRehearsals: 0,
			expectedSongProgress:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			songSectionRepository := new(repository.SongSectionRepositoryMock)
			songRepository := new(repository.SongRepositoryMock)
			artistRepository := new(repository.ArtistRepositoryMock)
			transactionManager := new(transaction.ManagerMock)
			_uut := part.NewCreateSongPart(
				songSectionRepository,
				songRepository,
				artistRepository,
				transactionManager,
			)

			// Set the song ID to match the request
			tt.song.ID = tt.request.SongID

			// Mock section validation
			if tt.request.SectionID != nil {
				mockSection := &model.SongSection{
					ID:     *tt.request.SectionID,
					SongID: tt.request.SongID,
				}
				songSectionRepository.On("Get", new(model.SongSection), *tt.request.SectionID).
					Return(nil, mockSection).
					Once()
			}

			songRepository.On("Get", new(model.Song), tt.request.SongID).
				Return(nil, &tt.song).
				Once()

			// Mock band member validation
			if tt.request.BandMemberID != nil {
				member := &model.BandMember{ArtistID: *tt.song.ArtistID}
				artistRepository.On("GetBandMember", new(model.BandMember), *tt.request.BandMemberID).
					Return(nil, member).
					Once()
			}

			repositoryFactory := new(transaction.RepositoryFactoryMock)
			txSongPartRepo := new(repository.SongPartRepositoryMock)
			txSongRepo := new(repository.SongRepositoryMock)
			txSongSectionRepo := new(repository.SongSectionRepositoryMock)
			txSongArrangementRepo := new(repository.SongArrangementRepositoryMock)

			repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
			repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
			repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepo).Once()
			repositoryFactory.On("NewSongArrangementRepository").Return(txSongArrangementRepo).Once()
			transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

			txSongPartRepo.On("CountAllBySong", new(int64), tt.request.SongID).
				Return(nil, &tt.partsCount).
				Once()

			// Count section parts only when a section is provided
			if tt.request.SectionID != nil {
				txSongSectionRepo.On("CountAllSectionPartsBySectionID", new(int64), *tt.request.SectionID).
					Return(nil, &tt.sectionPartsCount).
					Once()
			}

			var newPartID uuid.UUID
			txSongPartRepo.On("Create", mock.IsType(new(model.SongPart))).
				Run(func(args mock.Arguments) {
					newPart := args.Get(0).(*model.SongPart)
					newPartID = newPart.ID
					assertCreatedSongPart(t, tt.request, *newPart, tt.partsCount, tt.sectionPartsCount)
				}).
				Return(nil).
				Once()

			// updateSong
			txSongRepo.On("Update", mock.IsType(&tt.song)).
				Run(func(args mock.Arguments) {
					updatedSong := args.Get(0).(*model.Song)
					assert.Equal(t, tt.expectedSongConfidence, updatedSong.Confidence)
					assert.Equal(t, tt.expectedSongRehearsals, updatedSong.Rehearsals)
					assert.Equal(t, tt.expectedSongProgress, updatedSong.Progress)
				}).
				Return(nil).
				Once()

			// updateArrangements
			arrangements := []model.SongArrangement{
				{
					ID:              uuid.New(),
					Name:            "Arrangement 1",
					Order:           0,
					PartOccurrences: []model.SongPartOccurrences{{PartID: uuid.New(), Occurrences: 1}},
				},
				{ID: uuid.New(), Name: "Arrangement 2", Order: 1},
			}
			txSongArrangementRepo.On("GetAllBySong", new([]model.SongArrangement), tt.request.SongID).
				Return(nil, &arrangements).
				Once()

			oldArrangements := slices.Clone(arrangements)
			txSongArrangementRepo.On("UpdateAllWithAssociations", mock.IsType(&arrangements)).
				Run(func(args mock.Arguments) {
					newArrangements := args.Get(0).(*[]model.SongArrangement)
					for i, arr := range *newArrangements {
						assert.Len(t, arr.PartOccurrences, len(oldArrangements[i].PartOccurrences)+1)
						newOccurrence := arr.PartOccurrences[len(arr.PartOccurrences)-1]
						assert.Equal(t, newPartID, newOccurrence.PartID)
						assert.Equal(t, arr.ID, newOccurrence.ArrangementID)
						assert.Zero(t, newOccurrence.Occurrences)
					}
				}).
				Return(nil).
				Once()

			// when
			errCode := _uut.Handle(tt.request)

			// then
			assert.Nil(t, errCode)

			songSectionRepository.AssertExpectations(t)
			songRepository.AssertExpectations(t)
			artistRepository.AssertExpectations(t)
			transactionManager.AssertExpectations(t)
			repositoryFactory.AssertExpectations(t)
			txSongPartRepo.AssertExpectations(t)
			txSongRepo.AssertExpectations(t)
			txSongSectionRepo.AssertExpectations(t)
			txSongArrangementRepo.AssertExpectations(t)
		})
	}
}

func assertCreatedSongPart(
	t *testing.T,
	request requests.CreateSongPartRequest,
	part model.SongPart,
	partsCount int64,
	sectionPartsCount int64,
) {
	assert.NotEmpty(t, part.ID)
	assert.Equal(t, request.Name, part.Name)
	assert.Zero(t, part.Rehearsals)
	assert.Equal(t, model.DefaultSongPartConfidence, part.Confidence)
	assert.Zero(t, part.RehearsalsScore)
	assert.Zero(t, part.ConfidenceScore)
	assert.Zero(t, part.Progress)
	assert.Equal(t, uint(partsCount), part.SongOrder)
	assert.Equal(t, request.InstrumentID, part.InstrumentID)
	assert.Equal(t, request.SongID, part.SongID)

	if request.SectionID == nil {
		assert.Empty(t, part.SectionParts)
		return
	}

	assert.Len(t, part.SectionParts, 1)
	sp := part.SectionParts[0]
	assert.Equal(t, part.ID, sp.PartID)
	assert.Equal(t, *request.SectionID, sp.SectionID)
	assert.Equal(t, uint(sectionPartsCount), sp.Order)
	assert.Equal(t, request.BandMemberID, sp.BandMemberID)
}
