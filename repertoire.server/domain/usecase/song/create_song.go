package song

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type CreateSong struct {
	jwtService          service.JwtService
	repository          repository.SongRepository
	albumRepository     repository.AlbumRepository
	artistRepository    repository.ArtistRepository
	searchEngineService service.SearchEngineService
}

func NewCreateSong(
	jwtService service.JwtService,
	repository repository.SongRepository,
	albumRepository repository.AlbumRepository,
	artistRepository repository.ArtistRepository,
	searchEngineService service.SearchEngineService,
) CreateSong {
	return CreateSong{
		jwtService:          jwtService,
		repository:          repository,
		albumRepository:     albumRepository,
		artistRepository:    artistRepository,
		searchEngineService: searchEngineService,
	}
}

func (c CreateSong) Handle(request requests.CreateSongRequest, token string) (uuid.UUID, *wrapper.ErrorCode) {
	userID, errCode := c.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return uuid.Nil, errCode
	}

	songID := uuid.New()
	song := model.Song{
		ID:             songID,
		Title:          request.Title,
		Description:    request.Description,
		Bpm:            request.Bpm,
		SongsterrLink:  request.SongsterrLink,
		YoutubeLink:    request.YoutubeLink,
		ReleaseDate:    request.ReleaseDate,
		Difficulty:     request.Difficulty,
		GuitarTuningID: request.GuitarTuningID,
		AlbumID:        request.AlbumID,
		ArtistID:       request.ArtistID,
		Sections:       c.createSections(request.Sections, songID),
		UserID:         userID,
	}

	c.createArtist(&song, request)
	c.createAlbum(&song, request)

	errCode = c.addToAlbum(&song, request)
	if errCode != nil {
		return uuid.Nil, errCode
	}

	err := c.repository.Create(&song)
	if err != nil {
		return uuid.Nil, wrapper.InternalServerError(err)
	}

	errCode = c.addToSearchEngine(song)
	if errCode != nil {
		return uuid.Nil, errCode
	}

	return song.ID, nil
}

func (c CreateSong) createArtist(song *model.Song, request requests.CreateSongRequest) {
	if request.ArtistName == nil {
		return
	}

	song.Artist = &model.Artist{
		ID:     uuid.New(),
		Name:   *request.ArtistName,
		UserID: song.UserID,
	}
}

func (c CreateSong) createSections(request []requests.CreateSectionRequest, songID uuid.UUID) []model.SongSection {
	var sections []model.SongSection
	for i, sectionRequest := range request {
		sections = append(sections, model.SongSection{
			ID:                uuid.New(),
			Name:              sectionRequest.Name,
			Confidence:        model.DefaultSongSectionConfidence,
			SongSectionTypeID: sectionRequest.TypeID,
			Order:             uint(i),
			SongID:            songID,
		})
	}
	return sections
}

func (c CreateSong) createAlbum(song *model.Song, request requests.CreateSongRequest) {
	if request.AlbumTitle == nil {
		return
	}

	song.Album = &model.Album{
		ID:          uuid.New(),
		Title:       *request.AlbumTitle,
		UserID:      song.UserID,
		ArtistID:    song.ArtistID,    // album inherits the artist from song
		ReleaseDate: song.ReleaseDate, // also the release date
	}
	song.AlbumTrackNo = &[]uint{1}[0]
}

func (c CreateSong) addToAlbum(song *model.Song, request requests.CreateSongRequest) *wrapper.ErrorCode {
	if request.AlbumID == nil {
		return nil
	}

	var album model.Album
	err := c.albumRepository.GetWithSongs(&album, *request.AlbumID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(album).IsZero() {
		return wrapper.NotFoundError(errors.New("album not found"))
	}

	song.AlbumID = request.AlbumID
	trackNo := uint(len(album.Songs)) + 1
	song.AlbumTrackNo = &trackNo
	// song inherits the artist from album
	song.ArtistID = album.ArtistID
	if song.ReleaseDate == nil {
		song.ReleaseDate = album.ReleaseDate
	}

	return nil
}

func (c CreateSong) addToSearchEngine(song model.Song) *wrapper.ErrorCode {
	var searches []any
	songSearch := song.ToSearch()

	if song.ArtistID != nil {
		var artist model.Artist
		err := c.artistRepository.Get(&artist, *song.ArtistID)
		if err != nil {
			return wrapper.InternalServerError(err)
		}
		songSearch.Artist = artist.ToSongSearch()
	}
	if song.AlbumID != nil {
		var album model.Album
		err := c.albumRepository.Get(&album, *song.AlbumID)
		if err != nil {
			return wrapper.InternalServerError(err)
		}
		songSearch.Album = album.ToSongSearch()
	}

	searches = append(searches, songSearch)

	if song.Artist != nil {
		searches = append(searches, song.Artist.ToSearch())
	}
	if song.Album != nil {
		searches = append(searches, song.Album.ToSearch())
	}

	c.searchEngineService.Add(searches)
	return nil
}
