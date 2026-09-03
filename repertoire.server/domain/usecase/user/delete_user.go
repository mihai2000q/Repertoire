package user

import (
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/message/topics"
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

func (d DeleteUser) Handle(token string) *httperror.ErrorCode {
	id, errCode := d.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	if err := d.repository.Delete(id); err != nil {
		return httperror.DatabaseError(err)
	}

	if err := d.messagePublisherService.Publish(topics.UserDeletedTopic, id); err != nil {
		return httperror.MessagePublisherError(err)
	}

	return nil
}
