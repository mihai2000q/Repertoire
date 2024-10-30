package artist

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

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

func (c CreateArtist) Handle(request requests.CreateArtistRequest, token string) *wrapper.ErrorCode {
	userID, errCode := c.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	artist := model.Artist{
		ID:     uuid.New(),
		Name:   request.Name,
		UserID: userID,
	}
	err := c.repository.Create(&artist)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	return nil
}
