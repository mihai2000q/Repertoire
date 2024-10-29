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
		Album:          c.createAlbum(request, userID),
		Artist:         c.createArtist(request, userID),
		Sections:       c.createSections(request.Sections, songID),
		UserID:         userID,
	}
	err := c.repository.Create(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	return nil
}

func (c CreateSong) createAlbum(request requests.CreateSongRequest, userID uuid.UUID) *model.Album {
	var album *model.Album
	if request.AlbumTitle != nil {
		album = &model.Album{
			ID:     uuid.New(),
			Title:  *request.AlbumTitle,
			UserID: userID,
		}
	}
	return album
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
