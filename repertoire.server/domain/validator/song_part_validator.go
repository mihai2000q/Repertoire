package validator

import (
	"errors"
	"repertoire/server/data/repository"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type SongPartValidator interface {
	Validate(ids []uuid.UUID, songID uuid.UUID) *httperror.ErrorCode
}

type songPartValidator struct {
	songPartRepository repository.SongPartRepository
}

func NewSongPartValidator(songPartRepository repository.SongPartRepository) SongPartValidator {
	return songPartValidator{
		songPartRepository: songPartRepository,
	}
}

func (s songPartValidator) Validate(ids []uuid.UUID, songID uuid.UUID) *httperror.ErrorCode {
	if len(ids) == 0 {
		return nil
	}

	var parts []model.SongPart
	if err := s.songPartRepository.GetAllByIDs(&parts, ids); err != nil {
		return httperror.DatabaseError(err)
	}
	if len(parts) != len(ids) {
		return httperror.NotFoundError(errors.New("parts not found"))
	}

	for _, part := range parts {
		if part.SongID != songID {
			return httperror.ConflictError(errors.New("song part does not belong to the same song as the section"))
		}
	}

	return nil
}
