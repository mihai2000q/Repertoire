package search

import (
	"github.com/goccy/go-json"
	"repertoire/server/api/requests"
	"repertoire/server/data/service"
	"repertoire/server/internal/enums"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"strings"
)

type Get struct {
	jwtService         service.JwtService
	meiliSearchService service.SearchEngineService
}

func NewGet(
	jwtService service.JwtService,
	meiliSearchService service.SearchEngineService,
) Get {
	return Get{
		jwtService:         jwtService,
		meiliSearchService: meiliSearchService,
	}
}

func (g Get) Handle(
	request requests.SearchGetRequest,
	token string,
) (wrapper.WithTotalCount[any], *wrapper.ErrorCode) {
	userID, errCode := g.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return wrapper.WithTotalCount[any]{}, errCode
	}

	if len(request.IDs) > 0 {
		filter := "id IN ["
		for _, id := range request.IDs {
			filter = filter + string(*request.Type) + "-" + id + ", "
		}
		filter = strings.TrimRight(filter, ", ") + "]"
		request.Filter = append(request.Filter, filter)
	}

	if len(request.IDs) > 0 {
		filter := "id NOT IN ["
		for _, id := range request.IDs {
			filter = filter + string(*request.Type) + "-" + id + ", "
		}
		filter = strings.TrimRight(filter, ", ") + "]"
		request.Filter = append(request.Filter, filter)
	}

	searchResult, errCode := g.meiliSearchService.Search(
		request.Query,
		request.CurrentPage,
		request.PageSize,
		request.Type,
		userID,
		request.Filter,
		request.Order,
	)

	if errCode != nil {
		return wrapper.WithTotalCount[any]{}, errCode
	}

	var results []any
	for _, curr := range searchResult.Models {
		switch curr.(map[string]interface{})["type"] {
		case string(enums.Artist):
			var artist model.ArtistSearch
			jsonRes, _ := json.Marshal(curr)
			_ = json.Unmarshal(jsonRes, &artist)

			artist.ID = strings.Replace(artist.ID, "artist-", "", 1)
			artist.ImageUrl = artist.ImageUrl.ToFullURL(artist.UpdatedAt)

			results = append(results, artist)

		case string(enums.Album):
			var album model.AlbumSearch
			jsonRes, _ := json.Marshal(curr)
			_ = json.Unmarshal(jsonRes, &album)

			album.ID = strings.Replace(album.ID, "album-", "", 1)
			album.ImageUrl = album.ImageUrl.ToFullURL(album.UpdatedAt)
			if album.Artist != nil {
				album.Artist.ImageUrl = album.Artist.ImageUrl.ToFullURL(album.Artist.UpdatedAt)
			}

			results = append(results, album)

		case string(enums.Song):
			var song model.SongSearch
			jsonRes, _ := json.Marshal(curr)
			_ = json.Unmarshal(jsonRes, &song)

			song.ID = strings.Replace(song.ID, "song-", "", 1)
			song.ImageUrl = song.ImageUrl.ToFullURL(song.UpdatedAt)
			if song.Artist != nil {
				song.Artist.ImageUrl = song.Artist.ImageUrl.ToFullURL(song.Artist.UpdatedAt)
			}
			if song.Album != nil {
				song.Album.ImageUrl = song.Album.ImageUrl.ToFullURL(song.Album.UpdatedAt)
			}

			results = append(results, song)

		case string(enums.Playlist):
			var playlist model.PlaylistSearch
			jsonRes, _ := json.Marshal(curr)
			_ = json.Unmarshal(jsonRes, &playlist)

			playlist.ID = strings.Replace(playlist.ID, "playlist-", "", 1)
			playlist.ImageUrl = playlist.ImageUrl.ToFullURL(playlist.UpdatedAt)

			results = append(results, playlist)
		}
	}

	if len(results) == 0 {
		return wrapper.WithTotalCount[any]{
			Models:     []any{}, // otherwise it would be nil by default
			TotalCount: searchResult.TotalCount,
		}, nil
	}

	return wrapper.WithTotalCount[any]{
		Models:     results,
		TotalCount: searchResult.TotalCount,
	}, nil
}
