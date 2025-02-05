package artist

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/artist"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateArtist_WhenGetArtistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := artist.NewUpdateArtist(artistRepository)

	request := requests.UpdateArtistRequest{
		ID:   uuid.New(),
		Name: "New Artist",
	}

	internalError := errors.New("internal error")
	artistRepository.On("Get", new(model.Artist), request.ID).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
}

func TestUpdateArtist_WhenArtistIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := artist.NewUpdateArtist(artistRepository)

	request := requests.UpdateArtistRequest{
		ID:   uuid.New(),
		Name: "New Artist",
	}

	artistRepository.On("Get", new(model.Artist), request.ID).Return(nil).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "artist not found", errCode.Error.Error())

	artistRepository.AssertExpectations(t)
}

func TestUpdateArtist_WhenUpdateArtistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := artist.NewUpdateArtist(artistRepository)

	request := requests.UpdateArtistRequest{
		ID:   uuid.New(),
		Name: "New Artist",
	}

	mockArtist := &model.Artist{
		ID:   request.ID,
		Name: "Some Artist",
	}

	artistRepository.On("Get", new(model.Artist), request.ID).Return(nil, mockArtist).Once()
	internalError := errors.New("internal error")
	artistRepository.On("Update", mock.IsType(mockArtist)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
}

func TestUpdateArtist_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := artist.NewUpdateArtist(artistRepository)

	request := requests.UpdateArtistRequest{
		ID:     uuid.New(),
		Name:   "New Artist",
		IsBand: false,
	}

	mockArtist := &model.Artist{
		ID:     request.ID,
		Name:   "Some Artist",
		IsBand: true,
	}

	artistRepository.On("Get", new(model.Artist), request.ID).Return(nil, mockArtist).Once()
	artistRepository.On("Update", mock.IsType(mockArtist)).
		Run(func(args mock.Arguments) {
			newArtist := args.Get(0).(*model.Artist)
			assert.Equal(t, request.Name, newArtist.Name)
			assert.Equal(t, request.IsBand, newArtist.IsBand)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	artistRepository.AssertExpectations(t)
}
