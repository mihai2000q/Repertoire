package album

import (
	"repertoire/api/requests"
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

func (c CreateAlbum) Handle(request requests.CreateAlbumRequest, token string) *wrapper.ErrorCode {
	userID, errCode := c.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	album := model.Album{
		ID:     uuid.New(),
		Title:  request.Title,
		UserID: userID,
	}
	err := c.repository.Create(&album)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	return nil
}
