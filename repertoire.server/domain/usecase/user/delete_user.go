package user

import (
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/internal/wrapper"
)

type DeleteUser struct {
	repository              repository.UserRepository
	jwtService              service.JwtService
	messagePublisherService service.MessagePublisherService
}

func NewDeleteUser(
	repository repository.UserRepository,
	jwtService service.JwtService,
	messagePublisherService service.MessagePublisherService,
) DeleteUser {
	return DeleteUser{
		repository:              repository,
		jwtService:              jwtService,
		messagePublisherService: messagePublisherService,
	}
}

func (d DeleteUser) Handle(token string) *wrapper.ErrorCode {
	id, errCode := d.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	err := d.repository.Delete(id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = d.messagePublisherService.Publish(topics.UserDeletedTopic, id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
