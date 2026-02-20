package arrangement

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type UpdateDefaultSongArrangement struct {
	songRepository repository.SongRepository
}

func NewUpdateDefaultSongArrangement(songRepository repository.SongRepository) UpdateDefaultSongArrangement {
	return UpdateDefaultSongArrangement{songRepository: songRepository}
}

func (g UpdateDefaultSongArrangement) Handle(request requests.UpdateDefaultSongArrangementRequest) *wrapper.ErrorCode {
	var song model.Song
	err := g.songRepository.Get(&song, request.SongID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	song.DefaultArrangementID = &request.ID

	err = g.songRepository.Update(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	return nil
}
