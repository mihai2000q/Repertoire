package album

import (
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/models"
	"repertoire/utils"

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

func (c CreateAlbum) Handle(request requests.CreateAlbumRequest, token string) *utils.ErrorCode {
	userId, errCode := c.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	album := models.Album{
		ID:     uuid.New(),
		Title:  request.Title,
		UserID: userId,
	}
	err := c.repository.Create(&album)
	if err != nil {
		return utils.InternalServerError(err)
	}
	return nil
}
