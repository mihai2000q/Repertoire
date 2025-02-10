package album

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/album"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAlbum_WhenGetAlbumFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewGetAlbum(albumRepository)

	request := requests.GetAlbumRequest{
		ID:           uuid.New(),
		SongsOrderBy: []string{"ordering"},
	}

	internalError := errors.New("internal error")
	albumRepository.On("GetWithAssociations", new(model.Album), request.ID, request.SongsOrderBy).
		Return(internalError).
		Once()

	// when
	resultAlbum, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, resultAlbum)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
}

func TestGetAlbum_WhenAlbumIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewGetAlbum(albumRepository)

	request := requests.GetAlbumRequest{
		ID:           uuid.New(),
		SongsOrderBy: []string{"ordering"},
	}

	albumRepository.On("GetWithAssociations", new(model.Album), request.ID, request.SongsOrderBy).
		Return(nil).
		Once()

	// when
	resultAlbum, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, resultAlbum)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "album not found", errCode.Error.Error())

	albumRepository.AssertExpectations(t)
}

func TestGetAlbum_WhenSuccessful_ShouldReturnAlbum(t *testing.T) {
	tests := []struct {
		name    string
		request requests.GetAlbumRequest
	}{
		{
			"With Songs Order By",
			requests.GetAlbumRequest{
				ID:           uuid.New(),
				SongsOrderBy: []string{"ordering"},
			},
		},
		{
			"Without Songs Order By - Will have default value",
			requests.GetAlbumRequest{ID: uuid.New()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			albumRepository := new(repository.AlbumRepositoryMock)
			_uut := album.NewGetAlbum(albumRepository)

			expectedAlbum := &model.Album{
				ID:    tt.request.ID,
				Title: "Some Album",
			}

			if len(tt.request.SongsOrderBy) != 0 {
				albumRepository.
					On(
						"GetWithAssociations",
						new(model.Album),
						tt.request.ID,
						tt.request.SongsOrderBy,
					).
					Return(nil, expectedAlbum).
					Once()
			} else {
				albumRepository.
					On(
						"GetWithAssociations",
						new(model.Album),
						tt.request.ID,
						[]string{"album_track_no"},
					).
					Return(nil, expectedAlbum).
					Once()
			}

			// when
			resultAlbum, errCode := _uut.Handle(tt.request)

			// then
			assert.NotEmpty(t, resultAlbum)
			assert.Equal(t, expectedAlbum, &resultAlbum)
			assert.Nil(t, errCode)

			albumRepository.AssertExpectations(t)
		})
	}
}
