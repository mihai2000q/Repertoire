package album

import (
	"repertoire/api/request"
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/model"
	"repertoire/utils/wrapper"

	"github.com/google/uuid"
)

type CreateAlbum struct {
	jwtService service.JwtService
	repository repository.AlbumRepository
}

func NewCreateAlbum(jwtService service.JwtService, repository repository.AlbumRepository) CreateAlbum {
	return CreateAlbum{
		jwtService: jwtService,
		repository: repository,
	}
}

func (c CreateAlbum) Handle(request request.CreateAlbumRequest, token string) *wrapper.ErrorCode {
	userId, errCode := c.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	album := model.Album{
		ID:     uuid.New(),
		Title:  request.Title,
		UserID: userId,
	}
	err := c.repository.Create(&album)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	return nil
}
