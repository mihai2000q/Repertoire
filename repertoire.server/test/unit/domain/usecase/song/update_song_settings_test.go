package song

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"testing"
)

func TestUpdateSongSettings_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewUpdateSongSettings(songRepository)

	request := requests.UpdateSongSettingsRequest{
		SettingsID: uuid.New(),
	}

	internalError := errors.New("internal error")
	songRepository.On("GetSettings", new(model.SongSettings), request.SettingsID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestUpdateSongSettings_WhenSettingsAreEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewUpdateSongSettings(songRepository)

	request := requests.UpdateSongSettingsRequest{
		SettingsID: uuid.New(),
	}

	songRepository.On("GetSettings", new(model.SongSettings), request.SettingsID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "settings not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestUpdateSongSettings_WhenUpdateSettingsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewUpdateSongSettings(songRepository)

	request := requests.UpdateSongSettingsRequest{
		SettingsID: uuid.New(),
	}

	mockSettings := &model.SongSettings{
		ID: request.SettingsID,
	}
	songRepository.On("GetSettings", new(model.SongSettings), request.SettingsID).
		Return(nil, mockSettings).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("UpdateSettings", mock.IsType(mockSettings)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestUpdateSongSettings_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewUpdateSongSettings(songRepository)

	request := requests.UpdateSongSettingsRequest{
		SettingsID:          uuid.New(),
		DefaultInstrumentID: &[]uuid.UUID{uuid.New()}[0],
		DefaultBandMemberID: &[]uuid.UUID{uuid.New()}[0],
	}

	mockSettings := &model.SongSettings{
		ID: request.SettingsID,
	}

	songRepository.On("GetSettings", new(model.SongSettings), request.SettingsID).
		Return(nil, mockSettings).
		Once()

	songRepository.On("UpdateSettings", mock.IsType(mockSettings)).
		Run(func(args mock.Arguments) {
			newSettings := args.Get(0).(*model.SongSettings)
			assertUpdatedSongSettings(t, request, newSettings)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
}

func assertUpdatedSongSettings(t *testing.T, request requests.UpdateSongSettingsRequest, settings *model.SongSettings) {
	assert.Equal(t, request.DefaultInstrumentID, settings.DefaultInstrumentID)
	assert.Equal(t, request.DefaultBandMemberID, settings.DefaultBandMemberID)
}
