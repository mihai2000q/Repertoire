package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/models"
	"repertoire/utils"
	"testing"
)

// Get
func TestSongService_Get_WhenSongRepositoryReturnsError_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := &songService{
		repository: songRepository,
	}
	id := uuid.New()

	internalError := errors.New("internal error")
	songRepository.On("Get", new(models.Song), id).Return(internalError).Once()

	// when
	song, errCode := _uut.Get(id)

	// then
	assert.Empty(t, song)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestSongService_Get_WhenSongIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := &songService{
		repository: songRepository,
	}
	id := uuid.New()

	songRepository.On("Get", new(models.Song), id).Return(nil).Once()

	// when
	song, errCode := _uut.Get(id)

	// then
	assert.Empty(t, song)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestSongService_Get_WhenSuccessful_ShouldReturnSong(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := &songService{
		repository: songRepository,
	}
	id := uuid.New()

	expectedSong := &models.Song{
		ID:    id,
		Title: "Some Song",
	}

	songRepository.On("Get", new(models.Song), id).Return(nil, expectedSong).Once()

	// when
	song, errCode := _uut.Get(id)

	// then
	assert.NotEmpty(t, song)
	assert.Equal(t, expectedSong, &song)
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
}

// GetAll
func TestSongService_GetAll_WhenSongRepositoryReturnsError_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := &songService{
		repository: songRepository,
	}
	request := requests.GetSongsRequest{
		UserID: uuid.New(),
	}

	internalError := errors.New("internal error")
	songRepository.On("GetAllByUser", mock.Anything, request.UserID).
		Return(internalError).
		Once()

	// when
	songs, errCode := _uut.GetAll(request)

	// then
	assert.Empty(t, songs)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestSongService_GetAll_WhenSuccessful_ShouldReturnSongs(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := &songService{
		repository: songRepository,
	}
	request := requests.GetSongsRequest{
		UserID: uuid.New(),
	}

	expectedSongs := &[]models.Song{
		{Title: "Some Song"},
		{Title: "Some other Song"},
	}

	songRepository.On("GetAllByUser", mock.IsType(expectedSongs), request.UserID).
		Return(nil, expectedSongs).
		Once()

	// when
	songs, errCode := _uut.GetAll(request)

	// then
	assert.Equal(t, expectedSongs, &songs)
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
}

// Create
func TestSongService_Create_WhenJwtServiceReturnsErrorCode_ShouldReturnUnauthorizedError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	jwtService := new(service.JwtServiceMock)
	_uut := &songService{
		repository: songRepository,
		jwtService: jwtService,
	}
	request := requests.CreateSongRequest{
		Title: "Some Song",
	}
	token := "this is a token"

	unauthorizedError := utils.UnauthorizedError(errors.New("not authorized"))
	jwtService.On("GetUserIdFromJwt", token).Return(uuid.Nil, unauthorizedError).Once()

	// when
	errCode := _uut.Create(request, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, unauthorizedError, errCode)

	jwtService.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestSongService_Create_WhenSongRepositoryReturnsError_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	jwtService := new(service.JwtServiceMock)
	_uut := &songService{
		repository: songRepository,
		jwtService: jwtService,
	}
	request := requests.CreateSongRequest{
		Title:      "Some Song",
		IsRecorded: &[]bool{false}[0],
	}
	token := "this is a token"
	userID := uuid.New()

	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()
	internalError := errors.New("internal error")
	songRepository.On("Create", mock.IsType(new(models.Song))).
		Run(func(args mock.Arguments) {
			newSong := args.Get(0).(*models.Song)
			assert.Equal(t, request.Title, newSong.Title)
			assert.Equal(t, request.IsRecorded, newSong.IsRecorded)
			assert.Equal(t, userID, newSong.UserID)
		}).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Create(request, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestSongService_Create_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	jwtService := new(service.JwtServiceMock)
	_uut := &songService{
		repository: songRepository,
		jwtService: jwtService,
	}
	request := requests.CreateSongRequest{
		Title:      "Some Song",
		IsRecorded: &[]bool{true}[0],
	}
	token := "this is a token"
	userID := uuid.New()

	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()
	songRepository.On("Create", mock.IsType(new(models.Song))).
		Run(func(args mock.Arguments) {
			newSong := args.Get(0).(*models.Song)
			assert.Equal(t, request.Title, newSong.Title)
			assert.Equal(t, request.IsRecorded, newSong.IsRecorded)
			assert.Equal(t, userID, newSong.UserID)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Create(request, token)

	// then
	assert.Nil(t, errCode)

	jwtService.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

// Update
func TestSongService_Update_WhenSongRepositoryGetReturnsError_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := &songService{
		repository: songRepository,
	}
	request := requests.UpdateSongRequest{
		ID:    uuid.New(),
		Title: "New Song",
	}

	internalError := errors.New("internal error")
	songRepository.On("Get", new(models.Song), request.ID).Return(internalError).Once()

	// when
	errCode := _uut.Update(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestSongService_Update_WhenSongIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := &songService{
		repository: songRepository,
	}
	request := requests.UpdateSongRequest{
		ID:    uuid.New(),
		Title: "New Song",
	}

	songRepository.On("Get", new(models.Song), request.ID).Return(nil).Once()

	// when
	errCode := _uut.Update(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestSongService_Update_WhenSongRepositoryUpdateReturnsError_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := &songService{
		repository: songRepository,
	}
	request := requests.UpdateSongRequest{
		ID:    uuid.New(),
		Title: "New Song",
	}

	song := &models.Song{
		ID:    request.ID,
		Title: "Some Song",
	}

	songRepository.On("Get", new(models.Song), request.ID).Return(nil, song).Once()
	internalError := errors.New("internal error")
	songRepository.On("Update", mock.IsType(song)).
		Run(func(args mock.Arguments) {
			newSong := args.Get(0).(*models.Song)
			assert.Equal(t, request.Title, newSong.Title)
			assert.Equal(t, request.IsRecorded, newSong.IsRecorded)
		}).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Update(request)

	// then
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestSongService_Update_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := &songService{
		repository: songRepository,
	}
	request := requests.UpdateSongRequest{
		ID:    uuid.New(),
		Title: "New Song",
	}

	song := &models.Song{
		ID:    request.ID,
		Title: "Some Song",
	}

	songRepository.On("Get", new(models.Song), request.ID).Return(nil, song).Once()
	songRepository.On("Update", mock.IsType(song)).
		Run(func(args mock.Arguments) {
			newSong := args.Get(0).(*models.Song)
			assert.Equal(t, request.Title, newSong.Title)
			assert.Equal(t, request.IsRecorded, newSong.IsRecorded)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Update(request)

	// then
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
}

// Delete
func TestSongService_Delete_WhenSongRepositoryReturnsError_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := &songService{
		repository: songRepository,
	}
	request := requests.GetSongsRequest{
		UserID: uuid.New(),
	}

	internalError := errors.New("internal error")
	songRepository.On("GetAllByUser", mock.Anything, request.UserID).
		Return(internalError).
		Once()

	// when
	songs, errCode := _uut.GetAll(request)

	// then
	assert.Empty(t, songs)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestSongService_Delete_WhenSuccessful_ShouldReturnSongs(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := &songService{
		repository: songRepository,
	}
	id := uuid.New()

	songRepository.On("Delete", id).Return(nil).Once()

	// when
	errCode := _uut.Delete(id)

	// then
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
}
