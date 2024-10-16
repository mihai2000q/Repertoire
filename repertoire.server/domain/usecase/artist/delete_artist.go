package artist

import (
	"repertoire/data/repository"
	"repertoire/utils/wrapper"

	"github.com/google/uuid"
)

type DeleteArtist struct {
	repository repository.ArtistRepository
}

func NewDeleteArtist(repository repository.ArtistRepository) DeleteArtist {
	return DeleteArtist{repository: repository}
}

func (d DeleteArtist) Handle(id uuid.UUID) *wrapper.ErrorCode {
	err := d.repository.Delete(id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	return nil
}
