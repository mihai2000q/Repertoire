package song

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"repertoire/server/domain/usecase/song"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"repertoire/server/test/unit/domain/provider"
	"slices"
	"testing"
)

func TestDeleteSong_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewDeleteSong(songRepository, nil, nil)

	id := uuid.New()

	internalError := errors.New("internal error")
	songRepository.On("Get", mock.IsType(new(model.Song)), id).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestDeleteSong_WhenGetSongIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewDeleteSong(songRepository, nil, nil)

	id := uuid.New()

	songRepository.On("Get", mock.IsType(new(model.Song)), id).Return(nil).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestDeleteSong_WhenGetAllByAlbumAndTrackNoFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewDeleteSong(songRepository, nil, nil)

	id := uuid.New()

	mockSong := &model.Song{
		ID:           id,
		AlbumID:      &[]uuid.UUID{uuid.New()}[0],
		AlbumTrackNo: &[]uint{1}[0],
	}
	songRepository.On("Get", mock.IsType(mockSong), id).
		Return(nil, mockSong).
		Once()

	internalError := errors.New("internal error")
	songRepository.
		On(
			"GetAllByAlbumAndTrackNo",
			mock.IsType(&[]model.Song{}),
			*mockSong.AlbumID,
			*mockSong.AlbumTrackNo,
		).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestDeleteSong_WhenUpdateAllFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewDeleteSong(songRepository, nil, nil)

	id := uuid.New()

	mockSong := &model.Song{
		ID:           id,
		AlbumID:      &[]uuid.UUID{uuid.New()}[0],
		AlbumTrackNo: &[]uint{1}[0],
	}
	songRepository.On("Get", mock.IsType(mockSong), id).
		Return(nil, mockSong).
		Once()

	mockAlbumSongs := &[]model.Song{
		{AlbumTrackNo: &[]uint{1}[0]},
	}
	songRepository.
		On(
			"GetAllByAlbumAndTrackNo",
			mock.IsType(&[]model.Song{}),
			*mockSong.AlbumID,
			*mockSong.AlbumTrackNo,
		).
		Return(nil, mockAlbumSongs).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("UpdateAll", mock.IsType(mockAlbumSongs)).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestDeleteSong_WhenDeleteDirectoryFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	storageService := new(service.StorageServiceMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	_uut := song.NewDeleteSong(songRepository, storageService, storageFilePathProvider)

	id := uuid.New()

	mockSong := &model.Song{ID: id}
	songRepository.On("Get", mock.IsType(mockSong), id).Return(nil, mockSong).Once()

	storageFilePathProvider.On("HasSongFiles", *mockSong).Return(true).Once()

	directoryPath := "some directory path"
	storageFilePathProvider.On("GetSongDirectoryPath", *mockSong).Return(directoryPath).Once()

	internalError := errors.New("internal error")
	storageService.On("DeleteDirectory", directoryPath).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
}

func TestDeleteSong_WhenDeleteSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	_uut := song.NewDeleteSong(songRepository, nil, storageFilePathProvider)

	id := uuid.New()

	mockSong := &model.Song{ID: id}
	songRepository.On("Get", mock.IsType(mockSong), id).Return(nil, mockSong).Once()

	storageFilePathProvider.On("HasSongFiles", *mockSong).Return(false).Once()

	internalError := errors.New("internal error")
	songRepository.On("Delete", id).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
}

func TestDeleteSong_WhenSuccessful_ShouldDeleteSong(t *testing.T) {
	tests := []struct {
		name       string
		song       model.Song
		albumSongs []model.Song
		hasFiles   bool
	}{
		{
			"Normal delete, without album or files",
			model.Song{ID: uuid.New()},
			[]model.Song{},
			false,
		},
		{
			"With Files",
			model.Song{ID: uuid.New()},
			[]model.Song{},
			true,
		},
		{
			"With Album",
			model.Song{
				ID:           uuid.New(),
				AlbumID:      &[]uuid.UUID{uuid.New()}[0],
				AlbumTrackNo: &[]uint{2}[0],
			},
			[]model.Song{
				{AlbumTrackNo: &[]uint{3}[0]},
				{AlbumTrackNo: &[]uint{4}[0]},
				{AlbumTrackNo: &[]uint{5}[0]},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			songRepository := new(repository.SongRepositoryMock)
			storageService := new(service.StorageServiceMock)
			storageFilePathProvider := new(provider.StorageFilePathProviderMock)
			_uut := song.NewDeleteSong(songRepository, storageService, storageFilePathProvider)

			id := tt.song.ID

			songRepository.On("Get", mock.IsType(&tt.song), id).
				Return(nil, &tt.song).
				Once()
			songRepository.On("Delete", id).Return(nil).Once()

			if tt.song.AlbumID != nil {
				mockAlbumSongs := slices.Clone(tt.albumSongs)
				songRepository.
					On(
						"GetAllByAlbumAndTrackNo",
						mock.IsType(&tt.albumSongs),
						*tt.song.AlbumID,
						*tt.song.AlbumTrackNo,
					).
					Return(nil, &mockAlbumSongs).
					Once()

				songRepository.On("UpdateAll", mock.IsType(&tt.albumSongs)).
					Run(func(args mock.Arguments) {
						newSongs := args.Get(0).(*[]model.Song)
						for i := range *newSongs {
							assert.Equal(t, *tt.albumSongs[i].AlbumTrackNo-1, *(*newSongs)[i].AlbumTrackNo)
						}
					}).
					Return(nil).
					Once()
			}

			storageFilePathProvider.On("HasSongFiles", tt.song).Return(tt.hasFiles).Once()

			if tt.hasFiles {
				directoryPath := "some directory path"
				storageFilePathProvider.On("GetSongDirectoryPath", tt.song).
					Return(directoryPath).
					Once()

				storageService.On("DeleteDirectory", directoryPath).
					Return(nil).
					Once()
			}

			// when
			errCode := _uut.Handle(id)

			// then
			assert.Nil(t, errCode)

			songRepository.AssertExpectations(t)
			storageService.AssertExpectations(t)
			storageFilePathProvider.AssertExpectations(t)
		})
	}
}
