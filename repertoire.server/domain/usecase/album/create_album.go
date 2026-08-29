package album

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type CreateAlbum struct {
	jwtService              service.JwtService
	repository              repository.AlbumRepository
	messagePublisherService service.MessagePublisherService
}

func NewCreateAlbum(
	jwtService service.JwtService,
	repository repository.AlbumRepository,
	messagePublisherService service.MessagePublisherService,
) CreateAlbum {
	return CreateAlbum{
		jwtService:              jwtService,
		repository:              repository,
		messagePublisherService: messagePublisherService,
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
		UserID:      userID,
	}

	c.createArtist(&album, request, userID)

	err := c.repository.Create(&album)
	if err != nil {
		return uuid.Nil, wrapper.InternalServerError(err)
	}

	err = c.messagePublisherService.Publish(topics.AlbumCreatedTopic, album)
	if err != nil {
		return uuid.Nil, wrapper.InternalServerError(err)
	}

	return album.ID, nil
}

func (c CreateAlbum) createArtist(album *model.Album, request requests.CreateAlbumRequest, userID uuid.UUID) {
	if request.ArtistName == nil {
		return
	}
	album.Artist = &model.Artist{
		ID:     uuid.New(),
		Name:   *request.ArtistName,
		UserID: userID,
	}
	album.ArtistID = &album.Artist.ID
}
