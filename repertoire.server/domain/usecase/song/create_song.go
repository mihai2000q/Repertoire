package song

import (
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/model"
	"repertoire/utils/wrapper"

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

	song := model.Song{
		ID:             uuid.New(),
		Title:          request.Title,
		Description:    request.Description,
		Bpm:            request.Bpm,
		SongsterrLink:  request.SongsterrLink,
		GuitarTuningID: request.GuitarTuningID,
		AlbumID:        request.AlbumID,
		ArtistID:       request.ArtistID,
		Album:          c.createAlbum(request),
		Artist:         c.createArtist(request),
		Sections:       c.createSections(request.Sections),
		UserID:         userID,
	}
	err := c.repository.Create(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	return nil
}

func (c CreateSong) createAlbum(request requests.CreateSongRequest) *model.Album {
	var album *model.Album
	if request.AlbumTitle != nil {
		album = &model.Album{
			ID:    uuid.New(),
			Title: *request.AlbumTitle,
		}
	}
	return album
}

func (c CreateSong) createArtist(request requests.CreateSongRequest) *model.Artist {
	var artist *model.Artist
	if request.ArtistName != nil {
		artist = &model.Artist{
			ID:   uuid.New(),
			Name: *request.ArtistName,
		}
	}
	return artist
}

func (c CreateSong) createSections(request []requests.CreateSongSectionRequest) []model.SongSection {
	var sections []model.SongSection
	for _, sectionRequest := range request {
		sections = append(sections, model.SongSection{
			ID:                uuid.New(),
			Name:              sectionRequest.Name,
			SongSectionTypeID: sectionRequest.TypeId,
		})
	}
	return sections
}
