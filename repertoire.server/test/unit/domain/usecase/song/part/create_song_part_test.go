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

func TestCreateSongPart_WhenGetAllByIDsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	_uut := part.NewCreateSongPart(songSectionRepository, nil, nil)

	sectionID := uuid.New()
	request := requests.CreateSongPartRequest{
		SongID:     uuid.New(),
		SectionIDs: []uuid.UUID{sectionID},
	}

	internalError := errors.New("internal error")
	songSectionRepository.On("GetAllByIDs", new([]model.SongSection), request.SectionIDs).
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

func TestCreateSongPart_WhenSectionsAreNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	_uut := part.NewCreateSongPart(songSectionRepository, nil, nil)

	sectionID := uuid.New()
	request := requests.CreateSongPartRequest{
		SongID:     uuid.New(),
		SectionIDs: []uuid.UUID{sectionID},
	}

	songSectionRepository.On("GetAllByIDs", new([]model.SongSection), request.SectionIDs).
		Return(nil, &[]model.SongSection{}).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "sections not found", errCode.Error.Error())

	songSectionRepository.AssertExpectations(t)
}

func TestCreateSongPart_WhenSectionBelongsToDifferentSong_ShouldReturnConflictError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	_uut := part.NewCreateSongPart(songSectionRepository, nil, nil)

	sectionID := uuid.New()
	request := requests.CreateSongPartRequest{
		SongID:     uuid.New(),
		SectionIDs: []uuid.UUID{sectionID},
	}

	mockSections := []model.SongSection{
		{ID: sectionID, SongID: uuid.New()}, // different song
	}
	songSectionRepository.On("GetAllByIDs", new([]model.SongSection), request.SectionIDs).
		Return(nil, &mockSections).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusConflict, errCode.Code)
	assert.Equal(t, "section does not belong to the same song", errCode.Error.Error())

	songSectionRepository.AssertExpectations(t)
}

func TestCreateSongPart_WhenIsBandMemberAssociatedWithSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := part.NewCreateSongPart(nil, songRepository, nil)

	request := requests.CreateSongPartRequest{
		SongID:       uuid.New(),
		Name:         "Some Part",
		BandMemberID: &[]uuid.UUID{uuid.New()}[0],
	}

	internalError := errors.New("internal error")
	songRepository.On("IsBandMemberAssociatedWithSong", request.SongID, *request.BandMemberID).
		Return(false, internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestCreateSongPart_WhenBandMemberNotAssociated_ShouldReturnConflictError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := part.NewCreateSongPart(nil, songRepository, nil)

	request := requests.CreateSongPartRequest{
		SongID:       uuid.New(),
		Name:         "Some Part",
		BandMemberID: &[]uuid.UUID{uuid.New()}[0],
	}

	songRepository.On("IsBandMemberAssociatedWithSong", request.SongID, *request.BandMemberID).
		Return(false, nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusConflict, errCode.Code)
	assert.Equal(t, "band member is not part of the artist associated with this song", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestCreateSongPart_WhenTransactionExecuteFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewCreateSongPart(songSectionRepository, songRepository, transactionManager)

	request := requests.CreateSongPartRequest{
		SongID: uuid.New(),
		Name:   "Some Part",
	}

	internalError := errors.New("transaction failed")
	transactionManager.On("Execute", mock.Anything).Return(internalError).
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
}

func TestCreateSongPart_WhenCountAllBySongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewCreateSongPart(songSectionRepository, songRepository, transactionManager)

	request := requests.CreateSongPartRequest{
		SongID: uuid.New(),
		Name:   "Some Part",
	}

	// given - mocking
	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongPartRepository := new(repository.SongPartRepositoryMock)
	txSongRepository := new(repository.SongRepositoryMock)
	txSongArrangementRepository := new(repository.SongArrangementRepositoryMock)

	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(txSongRepository).Once()
	repositoryFactory.On("NewSongArrangementRepository").Return(txSongArrangementRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	internalError := errors.New("count error")
	txSongPartRepository.On("CountAllBySong", new(int64), request.SongID).
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
	txSongPartRepository.AssertExpectations(t)
	txSongRepository.AssertExpectations(t)
	txSongArrangementRepository.AssertExpectations(t)
}

func TestCreateSongPart_WhenCreatePartFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewCreateSongPart(songSectionRepository, songRepository, transactionManager)

	request := requests.CreateSongPartRequest{
		SongID: uuid.New(),
		Name:   "Some Part",
	}

	// given - mocking
	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongPartRepository := new(repository.SongPartRepositoryMock)
	txSongRepository := new(repository.SongRepositoryMock)
	txSongArrangementRepository := new(repository.SongArrangementRepositoryMock)

	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(txSongRepository).Once()
	repositoryFactory.On("NewSongArrangementRepository").Return(txSongArrangementRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txSongPartRepository.On("CountAllBySong", new(int64), request.SongID).
		Return(nil, &[]int64{0}[0]).
		Once()

	internalError := errors.New("create error")
	txSongPartRepository.On("Create", mock.IsType(new(model.SongPart))).
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
	txSongPartRepository.AssertExpectations(t)
	txSongRepository.AssertExpectations(t)
	txSongArrangementRepository.AssertExpectations(t)
}

func TestCreateSongPart_WhenCountBySectionIDsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewCreateSongPart(songSectionRepository, songRepository, transactionManager)

	sectionID := uuid.New()
	request := requests.CreateSongPartRequest{
		SongID:     uuid.New(),
		Name:       "Some Part",
		SectionIDs: []uuid.UUID{sectionID},
	}

	// given - mocking
	mockSections := []model.SongSection{
		{ID: sectionID, SongID: request.SongID},
	}
	songSectionRepository.On("GetAllByIDs", new([]model.SongSection), request.SectionIDs).
		Return(nil, &mockSections).
		Once()

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongPartRepository := new(repository.SongPartRepositoryMock)
	txSongRepository := new(repository.SongRepositoryMock)
	txSongSectionRepository := new(repository.SongSectionRepositoryMock)
	txSongArrangementRepository := new(repository.SongArrangementRepositoryMock)

	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(txSongRepository).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepository).Once()
	repositoryFactory.On("NewSongArrangementRepository").Return(txSongArrangementRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txSongPartRepository.On("CountAllBySong", new(int64), request.SongID).
		Return(nil, &[]int64{0}[0]).
		Once()

	txSongPartRepository.On("Create", mock.IsType(new(model.SongPart))).
		Return(nil).
		Once()

	internalError := errors.New("count sections error")
	txSongPartRepository.On("CountBySectionIDs", request.SectionIDs).
		Return(make(map[uuid.UUID]int64), internalError).
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
	txSongPartRepository.AssertExpectations(t)
	txSongRepository.AssertExpectations(t)
	txSongSectionRepository.AssertExpectations(t)
	txSongArrangementRepository.AssertExpectations(t)
}

func TestCreateSongPart_WhenCreateAllSectionPartsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewCreateSongPart(songSectionRepository, songRepository, transactionManager)

	sectionID := uuid.New()
	request := requests.CreateSongPartRequest{
		SongID:     uuid.New(),
		Name:       "Some Part",
		SectionIDs: []uuid.UUID{sectionID},
	}

	// given - mocking
	mockSections := []model.SongSection{
		{ID: sectionID, SongID: request.SongID},
	}
	songSectionRepository.On("GetAllByIDs", new([]model.SongSection), request.SectionIDs).
		Return(nil, &mockSections).
		Once()

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongPartRepository := new(repository.SongPartRepositoryMock)
	txSongRepository := new(repository.SongRepositoryMock)
	txSongSectionRepository := new(repository.SongSectionRepositoryMock)
	txSongArrangementRepository := new(repository.SongArrangementRepositoryMock)

	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(txSongRepository).Once()
	repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepository).Once()
	repositoryFactory.On("NewSongArrangementRepository").Return(txSongArrangementRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txSongPartRepository.On("CountAllBySong", new(int64), request.SongID).
		Return(nil, &[]int64{0}[0]).
		Once()

	txSongPartRepository.On("Create", mock.IsType(new(model.SongPart))).
		Return(nil).
		Once()

	// createSectionParts
	counts := map[uuid.UUID]int64{sectionID: 0}
	txSongPartRepository.On("CountBySectionIDs", request.SectionIDs).
		Return(counts, nil).
		Once()

	internalError := errors.New("create section parts error")
	txSongSectionRepository.On("CreateAllSectionParts", mock.IsType(new([]model.SongSectionPart))).
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
	txSongPartRepository.AssertExpectations(t)
	txSongRepository.AssertExpectations(t)
	txSongSectionRepository.AssertExpectations(t)
	txSongArrangementRepository.AssertExpectations(t)
}

func TestCreateSongPart_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewCreateSongPart(songSectionRepository, songRepository, transactionManager)

	request := requests.CreateSongPartRequest{
		SongID: uuid.New(),
		Name:   "Some Part",
	}

	// given - mocking
	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongPartRepository := new(repository.SongPartRepositoryMock)
	txSongRepository := new(repository.SongRepositoryMock)
	txSongArrangementRepository := new(repository.SongArrangementRepositoryMock)

	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(txSongRepository).Once()
	repositoryFactory.On("NewSongArrangementRepository").Return(txSongArrangementRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txSongPartRepository.On("CountAllBySong", new(int64), request.SongID).
		Return(nil, &[]int64{0}[0]).
		Once()

	txSongPartRepository.On("Create", mock.IsType(new(model.SongPart))).
		Return(nil).
		Once()

	internalError := errors.New("get song error")
	txSongRepository.On("Get", new(model.Song), request.SongID).
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
	txSongPartRepository.AssertExpectations(t)
	txSongRepository.AssertExpectations(t)
	txSongArrangementRepository.AssertExpectations(t)
}

func TestCreateSongPart_WhenSongNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewCreateSongPart(songSectionRepository, songRepository, transactionManager)

	request := requests.CreateSongPartRequest{
		SongID: uuid.New(),
		Name:   "Some Part",
	}

	// given - mocking
	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongPartRepository := new(repository.SongPartRepositoryMock)
	txSongRepository := new(repository.SongRepositoryMock)
	txSongArrangementRepository := new(repository.SongArrangementRepositoryMock)

	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(txSongRepository).Once()
	repositoryFactory.On("NewSongArrangementRepository").Return(txSongArrangementRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txSongPartRepository.On("CountAllBySong", new(int64), request.SongID).
		Return(nil, &[]int64{0}[0]).
		Once()

	txSongPartRepository.On("Create", mock.IsType(new(model.SongPart))).
		Return(nil).
		Once()

	txSongRepository.On("Get", new(model.Song), request.SongID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())

	songSectionRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongPartRepository.AssertExpectations(t)
	txSongRepository.AssertExpectations(t)
	txSongArrangementRepository.AssertExpectations(t)
}

func TestCreateSongPart_WhenUpdateSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewCreateSongPart(songSectionRepository, songRepository, transactionManager)

	request := requests.CreateSongPartRequest{
		SongID: uuid.New(),
		Name:   "Some Part",
	}

	// given - mocking
	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongPartRepository := new(repository.SongPartRepositoryMock)
	txSongRepository := new(repository.SongRepositoryMock)
	txSongArrangementRepository := new(repository.SongArrangementRepositoryMock)

	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(txSongRepository).Once()
	repositoryFactory.On("NewSongArrangementRepository").Return(txSongArrangementRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txSongPartRepository.On("CountAllBySong", new(int64), request.SongID).
		Return(nil, &[]int64{2}[0]).
		Once()

	txSongPartRepository.On("Create", mock.IsType(new(model.SongPart))).
		Return(nil).
		Once()

	// updateSong
	mockSong := &model.Song{
		ID:         request.SongID,
		Confidence: 50,
		Rehearsals: 10,
		Progress:   55,
	}
	txSongRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	internalError := errors.New("update error")
	txSongRepository.On("Update", mock.IsType(mockSong)).
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
	txSongPartRepository.AssertExpectations(t)
	txSongRepository.AssertExpectations(t)
	txSongArrangementRepository.AssertExpectations(t)
}

func TestCreateSongPart_WhenGetArrangementsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewCreateSongPart(songSectionRepository, songRepository, transactionManager)

	request := requests.CreateSongPartRequest{
		SongID: uuid.New(),
		Name:   "Some Part",
	}

	// given - mocking
	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongPartRepository := new(repository.SongPartRepositoryMock)
	txSongRepository := new(repository.SongRepositoryMock)
	txSongArrangementRepository := new(repository.SongArrangementRepositoryMock)

	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(txSongRepository).Once()
	repositoryFactory.On("NewSongArrangementRepository").Return(txSongArrangementRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txSongPartRepository.On("CountAllBySong", new(int64), request.SongID).
		Return(nil, &[]int64{0}[0]).
		Once()

	txSongPartRepository.On("Create", mock.IsType(new(model.SongPart))).
		Return(nil).
		Once()

	// updateSong
	mockSong := &model.Song{ID: request.SongID}
	txSongRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	txSongRepository.On("Update", mock.IsType(mockSong)).
		Return(nil).
		Once()

	internalError := errors.New("get arrangements error")
	txSongArrangementRepository.On("GetAllBySong", new([]model.SongArrangement), request.SongID).
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
	txSongPartRepository.AssertExpectations(t)
	txSongRepository.AssertExpectations(t)
	txSongArrangementRepository.AssertExpectations(t)
}

func TestCreateSongPart_WhenUpdateArrangementsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewCreateSongPart(songSectionRepository, songRepository, transactionManager)

	request := requests.CreateSongPartRequest{
		SongID: uuid.New(),
		Name:   "Some Part",
	}

	// given - mocking
	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongPartRepository := new(repository.SongPartRepositoryMock)
	txSongRepository := new(repository.SongRepositoryMock)
	txSongArrangementRepository := new(repository.SongArrangementRepositoryMock)

	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(txSongRepository).Once()
	repositoryFactory.On("NewSongArrangementRepository").Return(txSongArrangementRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txSongPartRepository.On("CountAllBySong", new(int64), request.SongID).
		Return(nil, &[]int64{0}[0]).
		Once()

	txSongPartRepository.On("Create", mock.IsType(new(model.SongPart))).
		Return(nil).
		Once()

	// updateSong
	mockSong := &model.Song{ID: request.SongID}
	txSongRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	txSongRepository.On("Update", mock.IsType(mockSong)).
		Return(nil).
		Once()

	// updateArrangements
	arrangements := []model.SongArrangement{{ID: uuid.New()}}
	txSongArrangementRepository.On("GetAllBySong", new([]model.SongArrangement), request.SongID).
		Return(nil, &arrangements).
		Once()

	internalError := errors.New("update arrangements error")
	txSongArrangementRepository.On("UpdateAllWithAssociations", mock.IsType(&arrangements)).
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
	txSongPartRepository.AssertExpectations(t)
	txSongRepository.AssertExpectations(t)
	txSongArrangementRepository.AssertExpectations(t)
}

func TestCreateSongPart_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	tests := []struct {
		name                   string
		request                requests.CreateSongPartRequest
		song                   model.Song
		partsCount             int64
		sectionCounts          map[uuid.UUID]int64
		expectedSongConfidence float64
		expectedSongRehearsals float64
		expectedSongProgress   float64
	}{
		{
			name: "No prior parts, no band member, no sections",
			request: requests.CreateSongPartRequest{
				SongID: uuid.New(),
				Name:   "Part 1",
			},
			song: model.Song{
				ID:         uuid.New(),
				Confidence: 0,
				Rehearsals: 0,
				Progress:   0,
			},
			partsCount:             0,
			sectionCounts:          nil,
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
			sectionCounts:          nil,
			expectedSongConfidence: 25,
			expectedSongRehearsals: 5,
			expectedSongProgress:   27,
		},
		{
			name: "With band member",
			request: requests.CreateSongPartRequest{
				SongID:       uuid.New(),
				Name:         "Part with member",
				BandMemberID: &[]uuid.UUID{uuid.New()}[0],
			},
			song: model.Song{
				ID:         uuid.New(),
				Confidence: 0,
				Rehearsals: 0,
				Progress:   0,
			},
			partsCount:             1,
			sectionCounts:          nil,
			expectedSongConfidence: 0,
			expectedSongRehearsals: 0,
			expectedSongProgress:   0,
		},
		{
			name: "With one section, no existing parts",
			request: requests.CreateSongPartRequest{
				SongID:     uuid.New(),
				Name:       "Part with one section",
				SectionIDs: []uuid.UUID{uuid.New()},
			},
			song: model.Song{
				ID:         uuid.New(),
				Confidence: 0,
				Rehearsals: 0,
				Progress:   0,
			},
			partsCount: 0,
			sectionCounts: map[uuid.UUID]int64{
				uuid.New(): 0,
			},
			expectedSongConfidence: 0,
			expectedSongRehearsals: 0,
			expectedSongProgress:   0,
		},
		{
			name: "With one section that already has 2 parts",
			request: requests.CreateSongPartRequest{
				SongID:     uuid.New(),
				Name:       "Part with existing section",
				SectionIDs: []uuid.UUID{uuid.New()},
			},
			song: model.Song{
				ID:         uuid.New(),
				Confidence: 0,
				Rehearsals: 0,
				Progress:   0,
			},
			partsCount: 0,
			sectionCounts: map[uuid.UUID]int64{
				uuid.New(): 2,
			},
			expectedSongConfidence: 0,
			expectedSongRehearsals: 0,
			expectedSongProgress:   0,
		},
		{
			name: "With multiple sections, different counts",
			request: requests.CreateSongPartRequest{
				SongID:     uuid.New(),
				Name:       "Part with sections",
				SectionIDs: []uuid.UUID{uuid.New(), uuid.New()},
			},
			song: model.Song{
				ID:         uuid.New(),
				Confidence: 0,
				Rehearsals: 0,
				Progress:   0,
			},
			partsCount: 0,
			sectionCounts: map[uuid.UUID]int64{
				uuid.New(): 3,
				uuid.New(): 1,
			},
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
			transactionManager := new(transaction.ManagerMock)
			_uut := part.NewCreateSongPart(songSectionRepository, songRepository, transactionManager)

			// Set the song ID to match the request
			tt.song.ID = tt.request.SongID

			// Mock section validation
			if len(tt.request.SectionIDs) > 0 {
				// Build sections for validation using the request's SectionIDs
				var sections []model.SongSection
				for _, secID := range tt.request.SectionIDs {
					sec := model.SongSection{
						ID:     secID,
						SongID: tt.request.SongID,
					}
					sections = append(sections, sec)
				}

				songSectionRepository.On("GetAllByIDs", new([]model.SongSection), tt.request.SectionIDs).
					Return(nil, &sections).
					Once()
			}

			// Mock band member association if provided
			if tt.request.BandMemberID != nil {
				songRepository.
					On(
						"IsBandMemberAssociatedWithSong",
						tt.request.SongID,
						*tt.request.BandMemberID,
					).
					Return(true, nil).
					Once()
			}

			// Transaction mocks
			repositoryFactory := new(transaction.RepositoryFactoryMock)
			txSongPartRepository := new(repository.SongPartRepositoryMock)
			txSongRepository := new(repository.SongRepositoryMock)
			txSongSectionRepository := new(repository.SongSectionRepositoryMock)
			txSongArrangementRepository := new(repository.SongArrangementRepositoryMock)

			repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepository).Once()
			repositoryFactory.On("NewSongRepository").Return(txSongRepository).Once()
			if len(tt.request.SectionIDs) > 0 {
				repositoryFactory.On("NewSongSectionRepository").Return(txSongSectionRepository).Once()
			}
			repositoryFactory.On("NewSongArrangementRepository").Return(txSongArrangementRepository).Once()
			transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

			txSongPartRepository.On("CountAllBySong", new(int64), tt.request.SongID).
				Return(nil, &tt.partsCount).
				Once()

			var newPartID uuid.UUID
			txSongPartRepository.On("Create", mock.IsType(new(model.SongPart))).
				Run(func(args mock.Arguments) {
					newPart := args.Get(0).(*model.SongPart)
					newPartID = newPart.ID
					assertCreatedSongPart(t, tt.request, *newPart, tt.partsCount)
				}).
				Return(nil).Once()

			// createSectionParts
			if len(tt.request.SectionIDs) > 0 {
				// Build sectionCounts with the actual section IDs from the request
				sectionCounts := make(map[uuid.UUID]int64)
				for _, secID := range tt.request.SectionIDs {
					// Use the provided count from tt.sectionCounts
					if count, ok := tt.sectionCounts[secID]; ok {
						sectionCounts[secID] = count
					}
				}

				txSongPartRepository.On("CountBySectionIDs", tt.request.SectionIDs).
					Return(sectionCounts, nil).
					Once()

				txSongSectionRepository.On("CreateAllSectionParts", mock.IsType(new([]model.SongSectionPart))).
					Run(func(args mock.Arguments) {
						sectionParts := args.Get(0).(*[]model.SongSectionPart)
						assert.Len(t, *sectionParts, len(tt.request.SectionIDs))
						for i, sp := range *sectionParts {
							assert.Equal(t, newPartID, sp.PartID)
							assert.Equal(t, tt.request.SectionIDs[i], sp.SectionID)
							assert.Equal(t, uint(sectionCounts[tt.request.SectionIDs[i]]), sp.Order)
						}
					}).
					Return(nil).
					Once()
			}

			// updateSong
			txSongRepository.On("Get", new(model.Song), tt.request.SongID).
				Return(nil, &tt.song).
				Once()

			txSongRepository.On("Update", mock.IsType(&tt.song)).
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
			txSongArrangementRepository.On("GetAllBySong", new([]model.SongArrangement), tt.request.SongID).
				Return(nil, &arrangements).
				Once()

			oldArrangements := slices.Clone(arrangements)
			txSongArrangementRepository.On("UpdateAllWithAssociations", mock.IsType(&arrangements)).
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
			transactionManager.AssertExpectations(t)
			repositoryFactory.AssertExpectations(t)
			txSongPartRepository.AssertExpectations(t)
			txSongRepository.AssertExpectations(t)
			if len(tt.request.SectionIDs) > 0 {
				txSongSectionRepository.AssertExpectations(t)
			}
			txSongArrangementRepository.AssertExpectations(t)
		})
	}
}

func assertCreatedSongPart(
	t *testing.T,
	request requests.CreateSongPartRequest,
	part model.SongPart,
	partsCount int64,
) {
	assert.NotEmpty(t, part.ID)
	assert.Equal(t, request.Name, part.Name)
	assert.Zero(t, part.Rehearsals)
	assert.Equal(t, model.DefaultSongPartConfidence, part.Confidence)
	assert.Zero(t, part.RehearsalsScore)
	assert.Zero(t, part.ConfidenceScore)
	assert.Zero(t, part.Progress)
	assert.Equal(t, uint(partsCount), part.SongOrder)
	assert.Equal(t, request.BandMemberID, part.BandMemberID)
	assert.Equal(t, request.InstrumentID, part.InstrumentID)
	assert.Equal(t, request.SongID, part.SongID)
}
