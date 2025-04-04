package playlist

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type CreatePlaylist struct {
	jwtService              service.JwtService
	repository              repository.PlaylistRepository
	messagePublisherService service.MessagePublisherService
}

func NewCreatePlaylist(
	jwtService service.JwtService,
	repository repository.PlaylistRepository,
	messagePublisherService service.MessagePublisherService,
) CreatePlaylist {
	return CreatePlaylist{
		jwtService:              jwtService,
		repository:              repository,
		messagePublisherService: messagePublisherService,
	}
}

func (c CreatePlaylist) Handle(request requests.CreatePlaylistRequest, token string) (uuid.UUID, *wrapper.ErrorCode) {
	userID, errCode := c.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return uuid.Nil, errCode
	}

	playlist := model.Playlist{
		ID:          uuid.New(),
		Title:       request.Title,
		Description: request.Description,
		UserID:      userID,
	}
	err := c.repository.Create(&playlist)
	if err != nil {
		return uuid.Nil, wrapper.InternalServerError(err)
	}

	err = c.messagePublisherService.Publish(topics.PlaylistCreatedTopic, playlist)
	if err != nil {
		return uuid.Nil, wrapper.InternalServerError(err)
	}

	return playlist.ID, nil
}
