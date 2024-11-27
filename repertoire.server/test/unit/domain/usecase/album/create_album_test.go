package album

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	album2 "repertoire/server/domain/usecase/album"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateAlbum_WhenGetUserIdFromJwtFails_ShouldReturnForbiddenError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := &album2.CreateAlbum{
		jwtService: jwtService,
	}
	request := requests.CreateAlbumRequest{
		Title: "Some Album",
	}
	token := "this is a token"

	forbiddenError := wrapper.ForbiddenError(errors.New("forbidden"))
	jwtService.On("GetUserIdFromJwt", token).Return(uuid.Nil, forbiddenError).Once()

	// when
	id, errCode := _uut.Handle(request, token)

	// then
	assert.Empty(t, id)
	assert.NotNil(t, errCode)
	assert.Equal(t, forbiddenError, errCode)

	jwtService.AssertExpectations(t)
}

func TestCreateAlbum_WhenGetAlbumFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	jwtService := new(service.JwtServiceMock)
	_uut := &album2.CreateAlbum{
		repository: albumRepository,
		jwtService: jwtService,
	}
	request := requests.CreateAlbumRequest{
		Title: "Some Album",
	}
	token := "this is a token"
	userID := uuid.New()

	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()
	internalError := errors.New("internal error")
	albumRepository.On("Create", mock.IsType(new(model.Album))).
		Return(internalError).
		Once()

	// when
	id, errCode := _uut.Handle(request, token)

	// then
	assert.Empty(t, id)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	albumRepository.AssertExpectations(t)
}

func TestCreateAlbum_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	tests := []struct {
		name    string
		request requests.CreateAlbumRequest
	}{
		{
			"With Existing Artist",
			requests.CreateAlbumRequest{
				Title:       "Some Album",
				ReleaseDate: &[]time.Time{time.Now()}[0],
				ArtistID:    &[]uuid.UUID{uuid.New()}[0],
			},
		},
		{
			"With New Artist",
			requests.CreateAlbumRequest{
				Title:      "Some Album",
				ArtistName: &[]string{"New Artist Name"}[0],
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			albumRepository := new(repository.AlbumRepositoryMock)
			jwtService := new(service.JwtServiceMock)
			_uut := &album2.CreateAlbum{
				repository: albumRepository,
				jwtService: jwtService,
			}
			request := requests.CreateAlbumRequest{
				Title: "Some Album",
			}
			token := "this is a token"
			userID := uuid.New()

			jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

			var albumID uuid.UUID
			albumRepository.On("Create", mock.IsType(new(model.Album))).
				Run(func(args mock.Arguments) {
					newAlbum := args.Get(0).(*model.Album)
					assertCreatedAlbum(t, *newAlbum, request, userID)
					albumID = newAlbum.ID
				}).
				Return(nil).
				Once()

			// when
			id, errCode := _uut.Handle(request, token)

			// then
			assert.Equal(t, albumID, id)
			assert.Nil(t, errCode)

			jwtService.AssertExpectations(t)
			albumRepository.AssertExpectations(t)
		})
	}
}

func assertCreatedAlbum(
	t *testing.T,
	album model.Album,
	request requests.CreateAlbumRequest,
	userID uuid.UUID,
) {
	assert.Equal(t, request.Title, album.Title)
	assert.Equal(t, request.ReleaseDate, album.ReleaseDate)
	assert.Nil(t, album.ImageURL)
	assert.Equal(t, request.ArtistID, album.ArtistID)
	assert.Equal(t, userID, album.UserID)
	if request.ArtistName != nil {
		assert.NotEmpty(t, album.Artist.ID)
		assert.Equal(t, request.ArtistName, album.Artist.Name)
		assert.Equal(t, userID, album.Artist.UserID)
	}
}
