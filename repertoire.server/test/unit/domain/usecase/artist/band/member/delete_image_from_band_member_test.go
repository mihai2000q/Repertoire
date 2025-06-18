package member

import (
	"errors"
	"net/http"
	"repertoire/server/domain/usecase/artist/band/member"
	"repertoire/server/internal"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteImageFromBandMember_WhenGetBandMemberFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := member.NewDeleteImageFromBandMember(artistRepository, nil)

	id := uuid.New()

	// given - mocking
	internalError := errors.New("internal error")
	artistRepository.On("GetBandMember", new(model.BandMember), id).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
}

func TestDeleteImageFromBandMember_WhenMemberIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := member.NewDeleteImageFromBandMember(artistRepository, nil)

	id := uuid.New()

	// given - mocking
	artistRepository.On("GetBandMember", new(model.BandMember), id).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "band member not found", errCode.Error.Error())

	artistRepository.AssertExpectations(t)
}

func TestDeleteImageFromBandMember_WhenMemberHasNoImage_ShouldReturnConflictError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := member.NewDeleteImageFromBandMember(artistRepository, nil)

	id := uuid.New()

	// given - mocking
	mockBandMember := &model.BandMember{ID: id}
	artistRepository.On("GetBandMember", new(model.BandMember), id).
		Return(nil, mockBandMember).
		Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusConflict, errCode.Code)
	assert.Equal(t, "band member does not have an image", errCode.Error.Error())

	artistRepository.AssertExpectations(t)
}

func TestDeleteImageFromBandMember_WhenDeleteImageFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	storageService := new(service.StorageServiceMock)
	_uut := member.NewDeleteImageFromBandMember(artistRepository, storageService)

	id := uuid.New()

	// given - mocking
	mockBandMember := &model.BandMember{ID: id, ImageURL: &[]internal.FilePath{"This is some url"}[0]}
	artistRepository.On("GetBandMember", new(model.BandMember), id).
		Return(nil, mockBandMember).
		Once()

	internalError := wrapper.InternalServerError(errors.New("internal error"))
	storageService.On("DeleteFile", *mockBandMember.ImageURL).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode)

	artistRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestDeleteImageFromBandMember_WhenUpdateBandMemberFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	storageService := new(service.StorageServiceMock)
	_uut := member.NewDeleteImageFromBandMember(artistRepository, storageService)

	id := uuid.New()

	// given - mocking
	mockBandMember := &model.BandMember{ID: id, ImageURL: &[]internal.FilePath{"This is some url"}[0]}
	artistRepository.On("GetBandMember", new(model.BandMember), id).
		Return(nil, mockBandMember).
		Once()

	storageService.On("DeleteFile", *mockBandMember.ImageURL).Return(nil).Once()

	internalError := errors.New("internal error")
	artistRepository.On("UpdateBandMember", mock.IsType(mockBandMember)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestDeleteImageFromBandMember_WhenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	storageService := new(service.StorageServiceMock)
	_uut := member.NewDeleteImageFromBandMember(artistRepository, storageService)

	id := uuid.New()

	// given - mocking
	mockBandMember := &model.BandMember{ID: id, ImageURL: &[]internal.FilePath{"This is some url"}[0]}
	artistRepository.On("GetBandMember", new(model.BandMember), id).
		Return(nil, mockBandMember).
		Once()

	storageService.On("DeleteFile", *mockBandMember.ImageURL).
		Return(nil).
		Once()

	artistRepository.On("UpdateBandMember", mock.IsType(mockBandMember)).
		Run(func(args mock.Arguments) {
			newBandMember := args.Get(0).(*model.BandMember)
			assert.Nil(t, newBandMember.ImageURL)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.Nil(t, errCode)

	artistRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
}
