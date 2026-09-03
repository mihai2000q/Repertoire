package song

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
)

type UpdateSongSettings struct {
	repository repository.SongRepository
}

func NewUpdateSongSettings(repository repository.SongRepository) UpdateSongSettings {
	return UpdateSongSettings{repository: repository}
}

func (u UpdateSongSettings) Handle(request requests.UpdateSongSettingsRequest) *httperror.ErrorCode {
	var settings model.SongSettings
	if err := u.repository.GetSettings(&settings, request.SettingsID); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(settings).IsZero() {
		return httperror.NotFoundError(errors.New("settings not found"))
	}

	settings.DefaultInstrumentID = request.DefaultInstrumentID
	settings.DefaultBandMemberID = request.DefaultBandMemberID

	if err := u.repository.UpdateSettings(&settings); err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
