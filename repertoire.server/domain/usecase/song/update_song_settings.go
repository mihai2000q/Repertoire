package song

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type UpdateSongSettings struct {
	repository repository.SongRepository
}

func NewUpdateSongSettings(repository repository.SongRepository) UpdateSongSettings {
	return UpdateSongSettings{repository: repository}
}

func (u UpdateSongSettings) Handle(request requests.UpdateSongSettingsRequest) *wrapper.ErrorCode {
	var settings model.SongSettings
	err := u.repository.GetSettings(&settings, request.SettingsID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(settings).IsZero() {
		return wrapper.NotFoundError(errors.New("settings not found"))
	}

	settings.DefaultInstrumentID = request.DefaultInstrumentID
	settings.DefaultBandMemberID = request.DefaultBandMemberID

	err = u.repository.UpdateSettings(&settings)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
