package search

import (
	"encoding/json"
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/search"
	"repertoire/server/internal"
	"repertoire/server/internal/enums"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/service"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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
		Return(wrapper.WithTotalCount[map[string]any]{}, errorCode).
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

	modelsResult := []map[string]any{
		{
			"id":     "artist-" + uuid.New().String(),
			"type":   enums.Artist,
			"name":   "Artist 1",
			"userID": userID,
		},
		{
			"id":       "artist-" + uuid.New().String(),
			"type":     enums.Artist,
			"name":     "Artist 2",
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
			"title":    "Album 2",
			"imageUrl": "something.png",
			"artist": model.AlbumArtistSearch{
				ID:   uuid.New(),
				Name: "Album Artist",
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
			"title":    "Song 2",
			"imageUrl": "something.png",
			"artist": model.SongArtistSearch{
				ID:       uuid.New(),
				Name:     "Song Artist",
				ImageUrl: &[]internal.FilePath{"something.png"}[0],
			},
			"album": model.SongAlbumSearch{
				ID:    uuid.New(),
				Title: "Song Album",
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

	searchResult := wrapper.WithTotalCount[map[string]any]{
		Models:     modelsResult,
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
		var currBase model.SearchBase
		jr, _ := json.Marshal(result.Models[i])
		_ = json.Unmarshal(jr, &currBase)
		expectedMap := searchResult.Models[i]

		assert.Equal(t, expectedMap["type"], currBase.Type)
		assert.Equal(t, expectedMap["userID"], currBase.UserID)

		switch currBase.Type {
		case enums.Artist:
			curr := result.Models[i].(model.ArtistSearch)
			assert.Equal(t, strings.Replace((expectedMap["id"]).(string), "artist-", "", 1), currBase.ID)
			assert.Equal(t, expectedMap["name"], curr.Name)
			if expectedMap["imageUrl"] != nil {
				filePath := internal.FilePath(expectedMap["imageUrl"].(string))
				assert.Equal(t, filePath.StripURL(), curr.ImageUrl)
			} else {
				assert.Nil(t, curr.ImageUrl)
			}
		case enums.Album:
			curr := result.Models[i].(model.AlbumSearch)
			assert.Equal(t, strings.Replace((expectedMap["id"]).(string), "album-", "", 1), currBase.ID)
			assert.Equal(t, expectedMap["title"], curr.Title)
			if expectedMap["imageUrl"] != nil {
				filePath := internal.FilePath(expectedMap["imageUrl"].(string))
				assert.Equal(t, filePath.StripURL(), curr.ImageUrl)
			} else {
				assert.Nil(t, curr.ImageUrl)
			}
			if expectedMap["artist"] == nil {
				assert.Nil(t, curr.Artist)
			} else {
				expectedArtist := expectedMap["artist"].(model.AlbumArtistSearch)
				assert.Equal(t, expectedArtist.ID, curr.Artist.ID)
				assert.Equal(t, expectedArtist.Name, curr.Artist.Name)
				assert.Equal(t, expectedArtist.ImageUrl.StripURL(), curr.Artist.ImageUrl)
			}
		case enums.Song:
			curr := result.Models[i].(model.SongSearch)
			assert.Equal(t, strings.Replace((expectedMap["id"]).(string), "song-", "", 1), currBase.ID)
			assert.Equal(t, expectedMap["title"], curr.Title)
			if expectedMap["imageUrl"] != nil {
				filePath := internal.FilePath(expectedMap["imageUrl"].(string))
				assert.Equal(t, filePath.StripURL(), curr.ImageUrl)
			} else {
				assert.Nil(t, curr.ImageUrl)
			}
			if expectedMap["artist"] == nil {
				assert.Nil(t, curr.Artist)
			} else {
				expectedArtist := expectedMap["artist"].(model.SongArtistSearch)
				assert.Equal(t, expectedArtist.ID, curr.Artist.ID)
				assert.Equal(t, expectedArtist.Name, curr.Artist.Name)
				assert.Equal(t, expectedArtist.ImageUrl.StripURL(), curr.Artist.ImageUrl)
			}
			if expectedMap["album"] == nil {
				assert.Nil(t, curr.Album)
			} else {
				expectedAlbum := expectedMap["album"].(model.SongAlbumSearch)
				assert.Equal(t, expectedAlbum.ID, curr.Album.ID)
				assert.Equal(t, expectedAlbum.Title, curr.Album.Title)
				assert.Equal(t, expectedAlbum.ImageUrl.StripURL(), curr.Album.ImageUrl)
			}
		case enums.Playlist:
			curr := result.Models[i].(model.PlaylistSearch)
			assert.Equal(t, strings.Replace((expectedMap["id"]).(string), "playlist-", "", 1), currBase.ID)
			assert.Equal(t, expectedMap["title"], curr.Title)
			if expectedMap["imageUrl"] != nil {
				filePath := internal.FilePath(expectedMap["imageUrl"].(string))
				assert.Equal(t, filePath.StripURL(), curr.ImageUrl)
			} else {
				assert.Nil(t, curr.ImageUrl)
			}
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

	modelsResult := []map[string]any{
		{
			"id":     "artist-" + uuid.New().String(),
			"type":   enums.Artist,
			"name":   "Artist 1",
			"userID": userID,
		},
		{
			"id":       "artist-" + uuid.New().String(),
			"type":     enums.Artist,
			"name":     "Artist 1",
			"imageUrl": "something.png",
			"userID":   userID,
		},
	}

	idsFilter := "id IN ["
	for _, id := range request.IDs {
		idsFilter = idsFilter + string(*request.Type) + "-" + id + ", "
	}
	idsFilter = strings.TrimSuffix(idsFilter, ", ") + "]"

	notIDsFilter := "id NOT IN ["
	for _, id := range request.NotIDs {
		notIDsFilter = notIDsFilter + string(*request.Type) + "-" + id + ", "
	}
	notIDsFilter = strings.TrimSuffix(notIDsFilter, ", ") + "]"

	filter := append(request.Filter, idsFilter, notIDsFilter)
	searchResult := wrapper.WithTotalCount[map[string]any]{
		Models:     modelsResult,
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
		var currBase model.SearchBase
		jr, _ := json.Marshal(result.Models[i])
		_ = json.Unmarshal(jr, &currBase)
		expectedMap := searchResult.Models[i]

		assert.Equal(t, expectedMap["type"], currBase.Type)
		assert.Equal(t, expectedMap["userID"], currBase.UserID)

		curr := result.Models[i].(model.ArtistSearch)
		assert.Equal(t, strings.Replace((expectedMap["id"]).(string), "artist-", "", 1), currBase.ID)
		assert.Equal(t, expectedMap["name"], curr.Name)
		if expectedMap["imageUrl"] != nil {
			filePath := internal.FilePath(expectedMap["imageUrl"].(string))
			assert.Equal(t, filePath.StripURL(), curr.ImageUrl)
		} else {
			assert.Nil(t, curr.ImageUrl)
		}
	}

	jwtService.AssertExpectations(t)
	searchEngineService.AssertExpectations(t)
}
