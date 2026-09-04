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
	_uut := part.NewCreateSongPart(songSectionRepository, nil, nil, nil)

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
	_uut := part.NewCreateSongPart(songSectionRepository, nil, nil, nil)

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
	_uut := part.NewCreateSongPart(songSectionRepository, nil, nil, nil)

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

func TestCreateSongPart_WhenSongIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
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

	mockSong := &model.Song{ID: uuid.New()}
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

func TestCreateSongPart_WhenBandMemberIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := part.NewCreateSongPart(nil, songRepository, artistRepository, nil)

	request := requests.CreateSongPartRequest{
		SongID:       uuid.New(),
		Name:         "Some Part",
		BandMemberID: &[]uuid.UUID{uuid.New()}[0],
	}

	mockSong := &model.Song{ID: uuid.New()}
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
	transactionManager.On("Execute", mock.Anything).Return(internalError).
		Once()

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

	// given - mocking
	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongArrangementRepo := new(repository.SongArrangementRepositoryMock)

	mockSong := &model.Song{ID: uuid.New()}
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
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
	txSongArrangementRepo.AssertExpectations(t)
}

func TestCreateSongPart_WhenCountBySectionIDsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := part.NewCreateSongPart(songSectionRepository, songRepository, nil, transactionManager)

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
	txSongRepo := new(repository.SongRepositoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)
	txSongArrangementRepo := new(repository.SongArrangementRepositoryMock)

	mockSong := &model.Song{ID: uuid.New()}
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongArrangementRepository").Return(txSongArrangementRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txSongPartRepo.On("CountAllBySong", new(int64), request.SongID).
		Return(nil, &[]int64{0}[0]).
		Once()

	internalError := errors.New("count sections error")
	txSongPartRepo.On("CountBySectionIDs", request.SectionIDs).
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
	txSongPartRepo.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
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

	// given - mocking
	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongArrangementRepo := new(repository.SongArrangementRepositoryMock)

	mockSong := &model.Song{ID: uuid.New()}
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
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

	// given - mocking
	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongArrangementRepo := new(repository.SongArrangementRepositoryMock)

	mockSong := &model.Song{ID: uuid.New()}
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongArrangementRepository").Return(txSongArrangementRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txSongPartRepo.On("CountAllBySong", new(int64), request.SongID).
		Return(nil, &[]int64{2}[0]).
		Once()

	txSongPartRepo.On("Create", mock.IsType(new(model.SongPart))).
		Return(nil).
		Once()

	// updateSong
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

	// given - mocking
	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongArrangementRepo := new(repository.SongArrangementRepositoryMock)

	mockSong := &model.Song{ID: uuid.New()}
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongArrangementRepository").Return(txSongArrangementRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txSongPartRepo.On("CountAllBySong", new(int64), request.SongID).
		Return(nil, &[]int64{0}[0]).
		Once()

	txSongPartRepo.On("Create", mock.IsType(new(model.SongPart))).
		Return(nil).
		Once()

	// updateSong
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

	// given - mocking
	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongArrangementRepo := new(repository.SongArrangementRepositoryMock)

	mockSong := &model.Song{ID: uuid.New()}
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongArrangementRepository").Return(txSongArrangementRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txSongPartRepo.On("CountAllBySong", new(int64), request.SongID).
		Return(nil, &[]int64{0}[0]).
		Once()

	txSongPartRepo.On("Create", mock.IsType(new(model.SongPart))).
		Return(nil).
		Once()

	// updateSong
	txSongRepo.On("Update", mock.IsType(mockSong)).
		Return(nil).
		Once()

	// updateArrangements
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
	txSongArrangementRepo.AssertExpectations(t)
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
				ArtistID:   &[]uuid.UUID{uuid.New()}[0],
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
			artistRepository := new(repository.ArtistRepositoryMock)
			transactionManager := new(transaction.ManagerMock)
			_uut := part.NewCreateSongPart(songSectionRepository, songRepository, artistRepository, transactionManager)

			// Transaction mocks
			repositoryFactory := new(transaction.RepositoryFactoryMock)
			txSongPartRepo := new(repository.SongPartRepositoryMock)
			txSongRepo := new(repository.SongRepositoryMock)
			txSongArrangementRepo := new(repository.SongArrangementRepositoryMock)

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

			songRepository.On("Get", new(model.Song), tt.request.SongID).
				Return(nil, &tt.song).
				Once()

			// Mock band member validation
			if tt.request.BandMemberID != nil {
				member := model.BandMember{ArtistID: *tt.song.ArtistID}
				artistRepository.
					On(
						"GetBandMember",
						new(model.BandMember),
						*tt.request.BandMemberID,
					).
					Return(nil, &member).
					Once()
			}

			repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
			repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
			repositoryFactory.On("NewSongArrangementRepository").Return(txSongArrangementRepo).Once()
			transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

			txSongPartRepo.On("CountAllBySong", new(int64), tt.request.SongID).
				Return(nil, &tt.partsCount).
				Once()

			// Build sectionCounts with the actual section IDs from the request
			sectionCounts := make(map[uuid.UUID]int64, len(tt.request.SectionIDs))
			for _, secID := range tt.request.SectionIDs {
				// Use the provided count from tt.sectionCounts
				if count, ok := tt.sectionCounts[secID]; ok {
					sectionCounts[secID] = count
				}
			}

			// createSectionParts
			if len(tt.request.SectionIDs) > 0 {
				txSongPartRepo.On("CountBySectionIDs", tt.request.SectionIDs).
					Return(sectionCounts, nil).
					Once()
			}

			var newPartID uuid.UUID
			txSongPartRepo.On("Create", mock.IsType(new(model.SongPart))).
				Run(func(args mock.Arguments) {
					newPart := args.Get(0).(*model.SongPart)
					newPartID = newPart.ID
					assertCreatedSongPart(t, tt.request, *newPart, tt.partsCount, sectionCounts)
				}).
				Return(nil).Once()

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
			txSongArrangementRepo.AssertExpectations(t)
		})
	}
}

func assertCreatedSongPart(
	t *testing.T,
	request requests.CreateSongPartRequest,
	part model.SongPart,
	partsCount int64,
	sectionCounts map[uuid.UUID]int64,
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

	assert.Len(t, part.SectionParts, len(request.SectionIDs))
	for i, sp := range part.SectionParts {
		assert.Equal(t, part.ID, sp.PartID)
		assert.Equal(t, request.SectionIDs[i], sp.SectionID)
		assert.Equal(t, uint(sectionCounts[request.SectionIDs[i]]), sp.Order)
	}
}
