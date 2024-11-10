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
	jwtService      service.JwtService
	repository      repository.SongRepository
	albumRepository repository.AlbumRepository
}

func NewCreateSong(
	jwtService service.JwtService,
	repository repository.SongRepository,
	albumRepository repository.AlbumRepository,
) CreateSong {
	return CreateSong{
		jwtService:      jwtService,
		repository:      repository,
		albumRepository: albumRepository,
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
		Artist:         c.createArtist(request, userID),
		Sections:       c.createSections(request.Sections, songID),
		UserID:         userID,
	}

	c.createAlbum(&song, request)

	errCode = c.addToAlbum(&song, request)
	if errCode != nil {
		return uuid.Nil, errCode
	}

	err := c.repository.Create(&song)
	if err != nil {
		return uuid.Nil, wrapper.InternalServerError(err)
	}
	return song.ID, nil
}

func (c CreateSong) createArtist(request requests.CreateSongRequest, userID uuid.UUID) *model.Artist {
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

func (c CreateSong) createSections(request []requests.CreateSectionRequest, songID uuid.UUID) []model.SongSection {
	var sections []model.SongSection
	for i, sectionRequest := range request {
		sections = append(sections, model.SongSection{
			ID:                uuid.New(),
			Name:              sectionRequest.Name,
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
		ID:       uuid.New(),
		Title:    *request.AlbumTitle,
		UserID:   song.UserID,
		ArtistID: song.ArtistID, // album inherits that artist from song
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

	return nil
}
