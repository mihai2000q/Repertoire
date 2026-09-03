package artist

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type CreateArtist struct {
	jwtService              service.JwtService
	repository              repository.ArtistRepository
	messagePublisherService service.MessagePublisherService
}

func NewCreateArtist(
	jwtService service.JwtService,
	repository repository.ArtistRepository,
	messagePublisherService service.MessagePublisherService,
) CreateArtist {
	return CreateArtist{
		jwtService:              jwtService,
		repository:              repository,
		messagePublisherService: messagePublisherService,
	}
}

func (c CreateArtist) Handle(request requests.CreateArtistRequest, token string) (uuid.UUID, *httperror.ErrorCode) {
	userID, errCode := c.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return uuid.Nil, errCode
	}

	artist := model.Artist{
		ID:     uuid.New(),
		Name:   request.Name,
		IsBand: request.IsBand,
		UserID: userID,
	}
	if err := c.repository.Create(&artist); err != nil {
		return uuid.Nil, httperror.DatabaseError(err)
	}

	if err := c.messagePublisherService.Publish(topics.ArtistCreatedTopic, artist); err != nil {
		return uuid.Nil, httperror.MessagePublisherError(err)
	}

	return artist.ID, nil
}
