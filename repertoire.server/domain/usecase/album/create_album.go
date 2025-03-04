package album

import (
	"github.com/google/uuid"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type CreateAlbum struct {
	jwtService          service.JwtService
	repository          repository.AlbumRepository
	artistRepository    repository.ArtistRepository
	searchEngineService service.SearchEngineService
}

func NewCreateAlbum(
	jwtService service.JwtService,
	repository repository.AlbumRepository,
	artistRepository repository.ArtistRepository,
	searchEngineService service.SearchEngineService,
) CreateAlbum {
	return CreateAlbum{
		jwtService:          jwtService,
		repository:          repository,
		artistRepository:    artistRepository,
		searchEngineService: searchEngineService,
	}
}

func (c CreateAlbum) Handle(request requests.CreateAlbumRequest, token string) (uuid.UUID, *wrapper.ErrorCode) {
	userID, errCode := c.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return uuid.Nil, errCode
	}

	album := model.Album{
		ID:          uuid.New(),
		Title:       request.Title,
		ReleaseDate: request.ReleaseDate,
		ArtistID:    request.ArtistID,
		Artist:      c.createArtist(request, userID),
		UserID:      userID,
	}
	err := c.repository.Create(&album)
	if err != nil {
		return uuid.Nil, wrapper.InternalServerError(err)
	}

	errCode = c.addToSearchEngine(album)
	if errCode != nil {
		return uuid.Nil, errCode
	}

	return album.ID, nil
}

func (c CreateAlbum) createArtist(request requests.CreateAlbumRequest, userID uuid.UUID) *model.Artist {
	var artist *model.Artist
	if request.ArtistName != nil {
		artist = &model.Artist{
			ID:     uuid.New(),
			Name:   *request.ArtistName,
			UserID: userID,
		}
	}
	return artist
}

func (c CreateAlbum) addToSearchEngine(album model.Album) *wrapper.ErrorCode {
	var searches []any
	albumSearch := album.ToSearch()

	if album.ArtistID != nil {
		var artist model.Artist
		err := c.artistRepository.Get(&artist, *album.ArtistID)
		if err != nil {
			return wrapper.InternalServerError(err)
		}
		albumSearch.Artist = artist.ToAlbumSearch()
	}
	searches = append(searches, albumSearch)
	if album.Artist != nil {
		searches = append(searches, album.Artist.ToSearch())
	}

	c.searchEngineService.Add(searches)
	return nil
}
