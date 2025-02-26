package search

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/service"
	"repertoire/server/internal/enums"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"strings"
)

type Get struct {
	jwtService         service.JwtService
	meiliSearchService service.MeiliSearchService
}

func NewGet(
	jwtService service.JwtService,
	meiliSearchService service.MeiliSearchService,
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

	searchResult, errCode := g.meiliSearchService.Get(
		request.Query,
		request.CurrentPage,
		request.PageSize,
		request.Type,
		userID,
	)

	var results []any
	for _, curr := range searchResult.Models {
		switch curr.(model.SearchBase).Type {
		case enums.Artist:
			var artist = curr.(model.ArtistSearch)
			artist.ID = strings.Replace(artist.ID, "artist-", "", 1)
			artist.ImageUrl = artist.ImageUrl.ToFullURL(artist.UpdatedAt)
			results = append(results, artist)

		case enums.Album:
			var album = curr.(model.AlbumSearch)
			album.ID = strings.Replace(album.ID, "album-", "", 1)
			album.ImageUrl = album.ImageUrl.ToFullURL(album.UpdatedAt)
			if album.Artist != nil {
				album.Artist.ImageUrl = album.Artist.ImageUrl.ToFullURL(album.Artist.UpdatedAt)
			}
			results = append(results, album)

		case enums.Song:
			var song = curr.(model.SongSearch)
			song.ID = strings.Replace(song.ID, "song-", "", 1)
			song.ImageUrl = song.ImageUrl.ToFullURL(song.UpdatedAt)
			if song.Artist != nil {
				song.Artist.ImageUrl = song.Artist.ImageUrl.ToFullURL(song.Artist.UpdatedAt)
			}
			if song.Album != nil {
				song.Album.ImageUrl = song.Album.ImageUrl.ToFullURL(song.Album.UpdatedAt)
			}
			results = append(results, song)

		case enums.Playlist:
			var playlist = curr.(model.PlaylistSearch)
			playlist.ID = strings.Replace(playlist.ID, "playlist-", "", 1)
			playlist.ImageUrl = playlist.ImageUrl.ToFullURL(playlist.UpdatedAt)
			results = append(results, playlist)
		}
	}

	return wrapper.WithTotalCount[any]{
		Models:     results,
		TotalCount: searchResult.TotalCount,
	}, errCode
}
