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
	jwtService          service.JwtService
	repository          repository.ArtistRepository
	searchEngineService service.SearchEngineService
}

func NewCreateArtist(
	jwtService service.JwtService,
	repository repository.ArtistRepository,
	searchEngineService service.SearchEngineService,
) CreateArtist {
	return CreateArtist{
		jwtService:          jwtService,
		repository:          repository,
		searchEngineService: searchEngineService,
	}
}

func (c CreateArtist) Handle(request requests.CreateArtistRequest, token string) (uuid.UUID, *wrapper.ErrorCode) {
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
	err := c.repository.Create(&artist)
	if err != nil {
		return uuid.Nil, wrapper.InternalServerError(err)
	}

	errCode = c.searchEngineService.Add([]any{artist.ToSearch()})
	if errCode != nil {
		return uuid.Nil, errCode
	}

	return artist.ID, nil
}
