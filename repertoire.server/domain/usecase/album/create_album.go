package album

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/message/topics"
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

func (c CreateAlbum) Handle(request requests.CreateAlbumRequest, token string) (uuid.UUID, *httperror.ErrorCode) {
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

	if err := c.repository.Create(&album); err != nil {
		return uuid.Nil, httperror.DatabaseError(err)
	}

	if err := c.messagePublisherService.Publish(topics.AlbumCreatedTopic, album); err != nil {
		return uuid.Nil, httperror.MessagePublisherError(err)
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
