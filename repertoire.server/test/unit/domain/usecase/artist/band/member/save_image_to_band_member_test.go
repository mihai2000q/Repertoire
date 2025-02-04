package member

import (
	"errors"
	"mime/multipart"
	"net/http"
	"repertoire/server/domain/usecase/artist/band/member"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"repertoire/server/test/unit/domain/provider"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSaveImageToBandMember_WhenGetBandMemberFails_ShouldReturnNotFoundError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := member.NewSaveImageToBandMember(artistRepository, nil, nil)

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	internalError := errors.New("internal error")
	artistRepository.On("GetBandMemberWithArtist", new(model.BandMember), id).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(file, id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
}

func TestSaveImageToBandMember_WhenMEmberIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := member.NewSaveImageToBandMember(artistRepository, nil, nil)

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	artistRepository.On("GetBandMemberWithArtist", new(model.BandMember), id).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(file, id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "band member not found", errCode.Error.Error())

	artistRepository.AssertExpectations(t)
}

func TestSaveImageToBandMember_WhenStorageUploadFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	storageService := new(service.StorageServiceMock)
	_uut := member.NewSaveImageToBandMember(artistRepository, storageFilePathProvider, storageService)

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	mockBandMember := &model.BandMember{ID: id, ImageURL: nil}
	artistRepository.On("GetBandMemberWithArtist", new(model.BandMember), id).
		Return(nil, mockBandMember).
		Once()

	imagePath := "artists file path"
	storageFilePathProvider.On("GetBandMemberImagePath", file, *mockBandMember).
		Return(imagePath).
		Once()

	internalError := errors.New("internal error")
	storageService.On("Upload", file, imagePath).Return(internalError).Once()

	// when
	errCode := _uut.Handle(file, id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestSaveImageToBandMember_WhenUpdateBandMemberFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	storageService := new(service.StorageServiceMock)
	_uut := member.NewSaveImageToBandMember(artistRepository, storageFilePathProvider, storageService)

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	mockBandMember := &model.BandMember{ID: id, ImageURL: nil}
	artistRepository.On("GetBandMemberWithArtist", new(model.BandMember), id).
		Return(nil, mockBandMember).
		Once()

	imagePath := "artists file path"
	storageFilePathProvider.On("GetBandMemberImagePath", file, *mockBandMember).
		Return(imagePath).
		Once()

	storageService.On("Upload", file, imagePath).Return(nil).Once()

	internalError := errors.New("internal error")
	artistRepository.On("UpdateBandMember", mock.IsType(new(model.BandMember))).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(file, id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestSaveImageToBandMember_WhenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	storageService := new(service.StorageServiceMock)
	_uut := member.NewSaveImageToBandMember(artistRepository, storageFilePathProvider, storageService)

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	mockBandMember := &model.BandMember{ID: id, ImageURL: nil}
	artistRepository.On("GetBandMemberWithArtist", new(model.BandMember), id).
		Return(nil, mockBandMember).
		Once()

	imagePath := "artists file path"
	storageFilePathProvider.On("GetBandMemberImagePath", file, *mockBandMember).
		Return(imagePath).
		Once()

	storageService.On("Upload", file, imagePath).Return(nil).Once()

	artistRepository.On("UpdateBandMember", mock.IsType(new(model.BandMember))).
		Run(func(args mock.Arguments) {
			newBandMember := args.Get(0).(*model.BandMember)
			assert.Equal(t, imagePath, string(*newBandMember.ImageURL))
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(file, id)

	// then
	assert.Nil(t, errCode)

	artistRepository.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
	storageService.AssertExpectations(t)
}
