package validator

import (
	"errors"
	"net/http"
	"repertoire/server/domain/validator"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSongPartValidator_WhenNoIds_ShouldReturnNothing(t *testing.T) {
	// given
	songPartRepository := new(repository.SongPartRepositoryMock)
	_uut := validator.NewSongPartValidator(songPartRepository)

	// when
	errCode := _uut.Validate(nil, uuid.New())

	// then
	assert.Nil(t, errCode)
	songPartRepository.AssertExpectations(t)
}

func TestSongPartValidator_WhenGetAllByIDsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songPartRepository := new(repository.SongPartRepositoryMock)
	_uut := validator.NewSongPartValidator(songPartRepository)

	ids := []uuid.UUID{uuid.New(), uuid.New()}
	songID := uuid.New()

	internalError := errors.New("get parts error")
	songPartRepository.On("GetAllByIDs", new([]model.SongPart), ids).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Validate(ids, songID)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songPartRepository.AssertExpectations(t)
}

func TestSongPartValidator_WhenSomePartsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songPartRepository := new(repository.SongPartRepositoryMock)
	_uut := validator.NewSongPartValidator(songPartRepository)

	ids := []uuid.UUID{uuid.New(), uuid.New()}
	songID := uuid.New()

	// only one of the two requested parts exists
	found := []model.SongPart{{ID: ids[0], SongID: songID}}
	songPartRepository.On("GetAllByIDs", new([]model.SongPart), ids).
		Return(nil, &found).
		Once()

	// when
	errCode := _uut.Validate(ids, songID)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "parts not found", errCode.Error.Error())

	songPartRepository.AssertExpectations(t)
}

func TestSongPartValidator_WhenAnyPartBelongsToDifferentSong_ShouldReturnConflictError(t *testing.T) {
	// given
	songPartRepository := new(repository.SongPartRepositoryMock)
	_uut := validator.NewSongPartValidator(songPartRepository)

	ids := []uuid.UUID{uuid.New(), uuid.New()}
	songID := uuid.New()

	found := []model.SongPart{
		{ID: ids[0], SongID: songID},
		{ID: ids[1], SongID: uuid.New()}, // different song
	}
	songPartRepository.On("GetAllByIDs", new([]model.SongPart), ids).
		Return(nil, &found).
		Once()

	// when
	errCode := _uut.Validate(ids, songID)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusConflict, errCode.Code)
	assert.Equal(t, "song part does not belong to the same song as the section", errCode.Error.Error())

	songPartRepository.AssertExpectations(t)
}

func TestSongPartValidator_WhenSuccessful_ShouldReturnNil(t *testing.T) {
	// given
	songPartRepository := new(repository.SongPartRepositoryMock)
	_uut := validator.NewSongPartValidator(songPartRepository)

	ids := []uuid.UUID{uuid.New(), uuid.New()}
	songID := uuid.New()

	found := []model.SongPart{
		{ID: ids[0], SongID: songID},
		{ID: ids[1], SongID: songID},
	}
	songPartRepository.On("GetAllByIDs", new([]model.SongPart), ids).
		Return(nil, &found).
		Once()

	// when
	errCode := _uut.Validate(ids, songID)

	// then
	assert.Nil(t, errCode)

	songPartRepository.AssertExpectations(t)
}
