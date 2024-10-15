package artist

import (
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/models"
	"repertoire/utils"

	"github.com/google/uuid"
)

type CreateArtist struct {
	jwtService service.JwtService
	repository repository.ArtistRepository
}

func NewCreateArtist(jwtService service.JwtService, repository repository.ArtistRepository) CreateArtist {
	return CreateArtist{
		jwtService: jwtService,
		repository: repository,
	}
}

func (c CreateArtist) Handle(request requests.CreateArtistRequest, token string) *utils.ErrorCode {
	userId, errCode := c.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	artist := models.Artist{
		ID:     uuid.New(),
		Name:   request.Name,
		UserID: userId,
	}
	err := c.repository.Create(&artist)
	if err != nil {
		return utils.InternalServerError(err)
	}
	return nil
}
