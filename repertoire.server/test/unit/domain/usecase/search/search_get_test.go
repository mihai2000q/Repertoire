package search

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/search"
	"repertoire/server/internal"
	"repertoire/server/internal/enums"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/service"
	"strings"
	"testing"
	"time"
)

func TestSearchGet_WhenJwtGetUserIDFails_ShouldReturnErrorCode(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := search.NewGet(jwtService, nil)

	request := requests.SearchGetRequest{
		Query: "test",
	}
	token := "some token"

	errorCode := &wrapper.ErrorCode{Error: errors.New("internalError"), Code: 400}
	jwtService.On("GetUserIdFromJwt", token).Return(uuid.Nil, errorCode).Once()

	// when
	result, errCode := _uut.Handle(request, token)

	// then
	assert.Empty(t, result)
	assert.NotNil(t, errCode)
	assert.Equal(t, errorCode, errCode)

	jwtService.AssertExpectations(t)
}

func TestSearchGet_WhenSearchEngineGetFails_ShouldReturnErrorCode(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	searchEngineService := new(service.SearchEngineServiceMock)
	_uut := search.NewGet(jwtService, searchEngineService)

	request := requests.SearchGetRequest{
		Query: "test",
	}
	token := "some token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	errorCode := &wrapper.ErrorCode{Error: errors.New("internalError"), Code: 400}
	searchEngineService.
		On(
			"Search",
			request.Query,
			request.CurrentPage,
			request.PageSize,
			request.Type,
			userID,
			request.Filter,
			request.Order,
		).
		Return(wrapper.WithTotalCount[any]{}, errorCode).
		Once()

	// when
	result, errCode := _uut.Handle(request, token)

	// then
	assert.Empty(t, result)
	assert.NotNil(t, errCode)
	assert.Equal(t, errorCode, errCode)

	jwtService.AssertExpectations(t)
	searchEngineService.AssertExpectations(t)
}

func TestSearchGet_WhenSuccessful_ShouldReturnSearchResult(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	searchEngineService := new(service.SearchEngineServiceMock)
	_uut := search.NewGet(jwtService, searchEngineService)

	request := requests.SearchGetRequest{
		Query:       "test",
		CurrentPage: &[]int{1}[0],
		PageSize:    &[]int{20}[0],
	}
	token := "some token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	modelsResult := []map[string]interface{}{
		{
			"id":     "artist-" + uuid.New().String(),
			"type":   enums.Artist,
			"title":  "Artist 1",
			"userID": userID,
		},
		{
			"id":       "artist-" + uuid.New().String(),
			"type":     enums.Artist,
			"title":    "Artist 1",
			"imageUrl": "something.png",
			"userID":   userID,
		},
		{
			"id":     "album-" + uuid.New().String(),
			"type":   enums.Album,
			"title":  "Album 1",
			"userID": userID,
		},
		{
			"id":       "album-" + uuid.New().String(),
			"type":     enums.Album,
			"title":    "Album 1",
			"imageUrl": "something.png",
			"artist": model.AlbumArtistSearch{
				ID:        uuid.New(),
				Name:      "Album Artist",
				UpdatedAt: time.Now(),
			},
			"userID": userID,
		},
		{
			"id":     "song-" + uuid.New().String(),
			"type":   enums.Song,
			"title":  "Song 1",
			"userID": userID,
		},
		{
			"id":       "song-" + uuid.New().String(),
			"type":     enums.Song,
			"title":    "Song 1",
			"imageUrl": "something.png",
			"artist": model.SongArtistSearch{
				ID:        uuid.New(),
				Name:      "Song Artist",
				UpdatedAt: time.Now(),
			},
			"album": model.SongAlbumSearch{
				ID:        uuid.New(),
				Title:     "Song Album",
				UpdatedAt: time.Now(),
			},
			"userID": userID,
		},
		{
			"id":     "playlist-" + uuid.New().String(),
			"type":   enums.Playlist,
			"title":  "Playlist 1",
			"userID": userID,
		},
	}

	var finalResult []any
	for _, m := range modelsResult {
		finalResult = append(finalResult, m)
	}

	searchResult := wrapper.WithTotalCount[any]{
		Models:     finalResult,
		TotalCount: int64(len(modelsResult)),
	}
	searchEngineService.
		On(
			"Search",
			request.Query,
			request.CurrentPage,
			request.PageSize,
			request.Type,
			userID,
			request.Filter,
			request.Order,
		).
		Return(searchResult, nil).
		Once()

	// when
	result, errCode := _uut.Handle(request, token)

	// then
	assert.NotEmpty(t, result)
	assert.Nil(t, errCode)

	assert.Equal(t, searchResult.TotalCount, result.TotalCount)
	for i := range result.Models {
		currBase := result.Models[i].(model.SearchBase)
		expectedMap := searchResult.Models[i].(map[string]interface{})

		assert.Equal(t, searchResult.Models[i].(map[string]interface{})["type"], currBase.Type)
		assert.Equal(t, searchResult.Models[i].(map[string]interface{})["userId"], currBase.UserID)

		switch currBase.Type {
		case enums.Artist:
			curr := result.Models[i].(model.ArtistSearch)
			assert.Equal(t, strings.Replace("artist-", (expectedMap["id"]).(string), "", 1), currBase.ID)
			assert.Equal(t, expectedMap["name"], curr.Name)
			assert.Equal(t, expectedMap["imageUrl"].(*internal.FilePath).StripURL(), curr.ImageUrl)
			assert.Equal(t, expectedMap["updatedAt"], curr.UpdatedAt)
		case enums.Album:
			curr := result.Models[i].(model.AlbumSearch)
			assert.Equal(t, strings.Replace("album-", (expectedMap["id"]).(string), "", 1), currBase.ID)
			assert.Equal(t, expectedMap["title"], curr.Title)
			assert.Equal(t, expectedMap["imageUrl"].(*internal.FilePath).StripURL(), curr.ImageUrl)
			assert.Equal(t, expectedMap["updatedAt"], curr.UpdatedAt)
			if expectedMap["artist"] == nil {
				assert.Nil(t, curr.Artist)
			} else {
				expectedArtist := expectedMap["artist"].(model.AlbumArtistSearch)
				assert.Equal(t, expectedArtist.ID, curr.Artist.ID)
				assert.Equal(t, expectedArtist.Name, curr.Artist.Name)
				assert.Equal(t, expectedArtist.ImageUrl.StripURL(), curr.Artist.ImageUrl)
				assert.Equal(t, expectedArtist.UpdatedAt, curr.Artist.UpdatedAt)
			}
		case enums.Song:
			curr := result.Models[i].(model.SongSearch)
			assert.Equal(t, strings.Replace("song-", (expectedMap["id"]).(string), "", 1), currBase.ID)
			assert.Equal(t, expectedMap["title"], curr.Title)
			assert.Equal(t, expectedMap["imageUrl"].(*internal.FilePath).StripURL(), curr.ImageUrl)
			assert.Equal(t, expectedMap["updatedAt"], curr.UpdatedAt)
			if expectedMap["artist"] == nil {
				assert.Nil(t, curr.Artist)
			} else {
				expectedArtist := expectedMap["artist"].(model.SongArtistSearch)
				assert.Equal(t, expectedArtist.ID, curr.Artist.ID)
				assert.Equal(t, expectedArtist.Name, curr.Artist.Name)
				assert.Equal(t, expectedArtist.ImageUrl.StripURL(), curr.Artist.ImageUrl)
				assert.Equal(t, expectedArtist.UpdatedAt, curr.Artist.UpdatedAt)
			}
			if expectedMap["album"] == nil {
				assert.Nil(t, curr.Album)
			} else {
				expectedAlbum := expectedMap["album"].(model.SongAlbumSearch)
				assert.Equal(t, expectedAlbum.ID, curr.Album.ID)
				assert.Equal(t, expectedAlbum.Title, curr.Album.UpdatedAt)
				assert.Equal(t, expectedAlbum.ImageUrl.StripURL(), curr.Album.ImageUrl)
				assert.Equal(t, expectedAlbum.UpdatedAt, curr.Album.UpdatedAt)
			}
		case enums.Playlist:
			curr := result.Models[i].(model.PlaylistSearch)
			assert.Equal(t, strings.Replace("playlist-", (expectedMap["id"]).(string), "", 1), currBase.ID)
			assert.Equal(t, expectedMap["title"], curr.Title)
			assert.Equal(t, expectedMap["imageUrl"].(*internal.FilePath).StripURL(), curr.ImageUrl)
			assert.Equal(t, expectedMap["updatedAt"], curr.UpdatedAt)
		}
	}

	jwtService.AssertExpectations(t)
	searchEngineService.AssertExpectations(t)
}

func TestSearchGet_WhenArtistsWithIDsAndNotIDs_ShouldReturnSearchResult(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	searchEngineService := new(service.SearchEngineServiceMock)
	_uut := search.NewGet(jwtService, searchEngineService)

	request := requests.SearchGetRequest{
		Query:       "test",
		CurrentPage: &[]int{1}[0],
		PageSize:    &[]int{20}[0],
		Type:        &[]enums.SearchType{enums.Artist}[0],
		Filter:      []string{"artist.id is null"},
		IDs:         []string{uuid.New().String(), uuid.New().String(), uuid.New().String()},
		NotIDs:      []string{uuid.New().String(), uuid.New().String(), uuid.New().String()},
	}
	token := "some token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	modelsResult := []map[string]interface{}{
		{
			"id":     "artist-" + uuid.New().String(),
			"type":   enums.Artist,
			"title":  "Artist 1",
			"userID": userID,
		},
		{
			"id":       "artist-" + uuid.New().String(),
			"type":     enums.Artist,
			"title":    "Artist 1",
			"imageUrl": "something.png",
			"userID":   userID,
		},
	}

	var finalResult []any
	for _, m := range modelsResult {
		finalResult = append(finalResult, m)
	}

	idsFilter := "id IN ["
	for _, id := range request.IDs {
		idsFilter = idsFilter + string(*request.Type) + "-" + id + ", "
	}
	idsFilter = strings.TrimRight(idsFilter, ", ") + "]"

	notIDsFilter := "id NOT IN ["
	for _, id := range request.IDs {
		notIDsFilter = notIDsFilter + string(*request.Type) + "-" + id + ", "
	}
	notIDsFilter = strings.TrimRight(notIDsFilter, ", ") + "]"

	filter := append(request.Filter, idsFilter, notIDsFilter)
	searchResult := wrapper.WithTotalCount[any]{
		Models:     finalResult,
		TotalCount: int64(len(modelsResult)),
	}
	searchEngineService.
		On(
			"Search",
			request.Query,
			request.CurrentPage,
			request.PageSize,
			request.Type,
			userID,
			filter,
			request.Order,
		).
		Return(searchResult, nil).
		Once()

	// when
	result, errCode := _uut.Handle(request, token)

	// then
	assert.NotEmpty(t, result)
	assert.Nil(t, errCode)

	assert.Equal(t, searchResult.TotalCount, result.TotalCount)
	for i := range result.Models {
		currBase := result.Models[i].(model.SearchBase)
		expectedMap := searchResult.Models[i].(map[string]interface{})

		assert.Equal(t, searchResult.Models[i].(map[string]interface{})["type"], currBase.Type)
		assert.Equal(t, searchResult.Models[i].(map[string]interface{})["userId"], currBase.UserID)

		curr := result.Models[i].(model.ArtistSearch)
		assert.Equal(t, strings.Replace("artist-", (expectedMap["id"]).(string), "", 1), currBase.ID)
		assert.Equal(t, expectedMap["name"], curr.Name)
		assert.Equal(t, expectedMap["imageUrl"].(*internal.FilePath).StripURL(), curr.ImageUrl)
		assert.Equal(t, expectedMap["updatedAt"], curr.UpdatedAt)
	}

	jwtService.AssertExpectations(t)
	searchEngineService.AssertExpectations(t)
}
