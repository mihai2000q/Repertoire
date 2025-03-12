package album

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/album"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateAlbum_WhenGetAlbumFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewUpdateAlbum(albumRepository, nil, nil)

	request := requests.UpdateAlbumRequest{
		ID:    uuid.New(),
		Title: "New Album",
	}

	internalError := errors.New("internal error")
	albumRepository.On("Get", new(model.Album), request.ID).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
}

func TestUpdateAlbum_WhenAlbumIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewUpdateAlbum(albumRepository, nil, nil)

	request := requests.UpdateAlbumRequest{
		ID:    uuid.New(),
		Title: "New Album",
	}

	albumRepository.On("Get", new(model.Album), request.ID).Return(nil).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "album not found", errCode.Error.Error())

	albumRepository.AssertExpectations(t)
}

func TestUpdateAlbum_WhenUpdateAlbumFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewUpdateAlbum(albumRepository, nil, nil)

	request := requests.UpdateAlbumRequest{
		ID:    uuid.New(),
		Title: "New Album",
	}

	mockAlbum := &model.Album{
		ID:    request.ID,
		Title: "Some Album",
	}

	albumRepository.On("Get", new(model.Album), request.ID).Return(nil, mockAlbum).Once()

	internalError := errors.New("internal error")
	albumRepository.On("Update", mock.IsType(mockAlbum)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
}

func TestUpdateAlbum_WhenGetSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := album.NewUpdateAlbum(albumRepository, songRepository, nil)

	request := requests.UpdateAlbumRequest{
		ID:       uuid.New(),
		Title:    "New Album",
		ArtistID: &[]uuid.UUID{uuid.New()}[0],
	}

	mockAlbum := &model.Album{
		ID:    request.ID,
		Title: "Some Album",
	}

	albumRepository.On("Get", new(model.Album), request.ID).Return(nil, mockAlbum).Once()
	albumRepository.On("Update", mock.IsType(mockAlbum)).Return(nil).Once()

	internalError := errors.New("internal error")
	songRepository.On("GetAllByAlbum", new([]model.Song), request.ID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestUpdateAlbum_WhenUpdateSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := album.NewUpdateAlbum(albumRepository, songRepository, nil)

	request := requests.UpdateAlbumRequest{
		ID:       uuid.New(),
		Title:    "New Album",
		ArtistID: &[]uuid.UUID{uuid.New()}[0],
	}

	mockAlbum := &model.Album{
		ID:    request.ID,
		Title: "Some Album",
	}

	albumRepository.On("Get", new(model.Album), request.ID).Return(nil, mockAlbum).Once()
	albumRepository.On("Update", mock.IsType(mockAlbum)).Return(nil).Once()

	songs := []model.Song{{ID: uuid.New()}}

	songRepository.On("GetAllByAlbum", new([]model.Song), request.ID).
		Return(nil, &songs).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("UpdateAll", mock.IsType(&[]model.Song{})).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestUpdateAlbum_WhenPublishFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := album.NewUpdateAlbum(albumRepository, nil, messagePublisherService)

	request := requests.UpdateAlbumRequest{
		ID:    uuid.New(),
		Title: "New Album",
	}

	mockAlbum := &model.Album{
		ID:    request.ID,
		Title: "Some Album",
	}

	albumRepository.On("Get", new(model.Album), request.ID).Return(nil, mockAlbum).Once()
	albumRepository.On("Update", mock.IsType(mockAlbum)).Return(nil).Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.AlbumsUpdatedTopic, []uuid.UUID{mockAlbum.ID}).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}

func TestUpdateAlbum_WhenArtistHasNotChanged_ShouldUpdateOnlyAlbumAndNotReturnAnyError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := album.NewUpdateAlbum(albumRepository, nil, messagePublisherService)

	request := requests.UpdateAlbumRequest{
		ID:          uuid.New(),
		Title:       "New Album",
		ReleaseDate: &[]time.Time{time.Now().UTC()}[0],
	}

	mockAlbum := &model.Album{
		ID:    request.ID,
		Title: "Some Album",
	}

	albumRepository.On("Get", new(model.Album), request.ID).
		Return(nil, mockAlbum).
		Once()

	albumRepository.On("Update", mock.IsType(mockAlbum)).
		Run(func(args mock.Arguments) {
			newAlbum := args.Get(0).(*model.Album)
			assertUpdatedAlbum(t, *newAlbum, request)
		}).
		Return(nil).
		Once()

	messagePublisherService.On("Publish", topics.AlbumsUpdatedTopic, []uuid.UUID{mockAlbum.ID}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	albumRepository.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}

func TestUpdateAlbum_WhenArtistHasChanged_ShouldUpdateAlbumAndSongsAndNotReturnAnyError(t *testing.T) {
	albumID := uuid.New()

	tests := []struct {
		name    string
		request requests.UpdateAlbumRequest
		album   model.Album
	}{
		{
			"Artist has changed",
			requests.UpdateAlbumRequest{
				ID:       albumID,
				Title:    "New Album",
				ArtistID: &[]uuid.UUID{uuid.New()}[0],
			},
			model.Album{
				ID:       albumID,
				Title:    "Some Album",
				ArtistID: &[]uuid.UUID{uuid.New()}[0],
			},
		},
		{
			"Artist has been added",
			requests.UpdateAlbumRequest{
				ID:       albumID,
				Title:    "New Album",
				ArtistID: &[]uuid.UUID{uuid.New()}[0],
			},
			model.Album{
				ID:    albumID,
				Title: "Some Album",
			},
		},
		{
			"Artist has been removed",
			requests.UpdateAlbumRequest{
				ID:    albumID,
				Title: "New Album",
			},
			model.Album{
				ID:       albumID,
				Title:    "Some Album",
				ArtistID: &[]uuid.UUID{uuid.New()}[0],
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			albumRepository := new(repository.AlbumRepositoryMock)
			songRepository := new(repository.SongRepositoryMock)
			messagePublisherService := new(service.MessagePublisherServiceMock)
			_uut := album.NewUpdateAlbum(albumRepository, songRepository, messagePublisherService)

			albumRepository.On("Get", new(model.Album), tt.request.ID).
				Return(nil, &tt.album).
				Once()
			albumRepository.On("Update", mock.IsType(&tt.album)).
				Run(func(args mock.Arguments) {
					newAlbum := args.Get(0).(*model.Album)
					assertUpdatedAlbum(t, *newAlbum, tt.request)
				}).
				Return(nil).
				Once()

			songs := []model.Song{{ID: uuid.New()}, {ID: uuid.New()}, {ID: uuid.New()}}

			songRepository.On("GetAllByAlbum", new([]model.Song), tt.request.ID).
				Return(nil, &songs).
				Once()

			songRepository.On("UpdateAll", mock.IsType(&[]model.Song{})).
				Run(func(args mock.Arguments) {
					newSongs := args.Get(0).(*[]model.Song)
					for _, song := range *newSongs {
						assert.Equal(t, tt.request.ArtistID, song.ArtistID)
					}
				}).
				Return(nil).
				Once()

			messagePublisherService.On("Publish", topics.AlbumsUpdatedTopic, []uuid.UUID{tt.request.ID}).
				Return(nil).
				Once()

			// when
			errCode := _uut.Handle(tt.request)

			// then
			assert.Nil(t, errCode)

			albumRepository.AssertExpectations(t)
			songRepository.AssertExpectations(t)
			messagePublisherService.AssertExpectations(t)
		})
	}
}

func assertUpdatedAlbum(
	t *testing.T,
	album model.Album,
	request requests.UpdateAlbumRequest,
) {
	assert.Equal(t, request.Title, album.Title)
	assert.Equal(t, request.ReleaseDate, album.ReleaseDate)
	assert.Equal(t, request.ArtistID, album.ArtistID)
}
