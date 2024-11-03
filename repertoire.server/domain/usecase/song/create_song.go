package song

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type CreateSong struct {
	jwtService service.JwtService
	repository repository.SongRepository
}

func NewCreateSong(jwtService service.JwtService, repository repository.SongRepository) CreateSong {
	return CreateSong{
		jwtService: jwtService,
		repository: repository,
	}
}

func (c CreateSong) Handle(request requests.CreateSongRequest, token string) *wrapper.ErrorCode {
	userID, errCode := c.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	songID := uuid.New()
	song := model.Song{
		ID:             songID,
		Title:          request.Title,
		Description:    request.Description,
		Bpm:            request.Bpm,
		SongsterrLink:  request.SongsterrLink,
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

	err := c.addToAlbum(&song, request.AlbumID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = c.repository.Create(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	return nil
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
		ID:     uuid.New(),
		Title:  *request.AlbumTitle,
		UserID: song.UserID,
	}
	song.AlbumTrackNo = &[]uint{1}[0]
}

func (c CreateSong) addToAlbum(song *model.Song, albumID *uuid.UUID) error {
	if albumID == nil {
		return nil
	}

	var count int64
	err := c.repository.CountByAlbum(&count, albumID)
	if err != nil {
		return err
	}

	song.AlbumID = albumID
	trackNo := uint(count) + 1
	song.AlbumTrackNo = &trackNo

	return nil
}
