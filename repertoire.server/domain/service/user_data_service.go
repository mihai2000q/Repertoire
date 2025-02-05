package service

import (
	"github.com/google/uuid"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/udata/guitar/tuning"
	"repertoire/server/domain/usecase/udata/section/types"
	"repertoire/server/internal/wrapper"
)

type UserDataService interface {
	CreateGuitarTuning(request requests.CreateGuitarTuningRequest, token string) *wrapper.ErrorCode
	MoveGuitarTuning(request requests.MoveGuitarTuningRequest, token string) *wrapper.ErrorCode
	DeleteGuitarTuning(id uuid.UUID, token string) *wrapper.ErrorCode

	CreateSectionType(request requests.CreateSongSectionTypeRequest, token string) *wrapper.ErrorCode
	DeleteSectionType(id uuid.UUID, token string) *wrapper.ErrorCode
	MoveSectionType(request requests.MoveSongSectionTypeRequest, token string) *wrapper.ErrorCode
}

type userDataService struct {
	createGuitarTuning tuning.CreateGuitarTuning
	deleteGuitarTuning tuning.DeleteGuitarTuning
	moveGuitarTuning   tuning.MoveGuitarTuning

	createSongSectionType types.CreateSongSectionType
	deleteSongSectionType types.DeleteSongSectionType
	moveSongSectionType   types.MoveSongSectionType
}

func NewUserDataService(
	createGuitarTuning tuning.CreateGuitarTuning,
	deleteGuitarTuning tuning.DeleteGuitarTuning,
	moveGuitarTuning tuning.MoveGuitarTuning,

	createSongSectionType types.CreateSongSectionType,
	deleteSongSectionType types.DeleteSongSectionType,
	moveSongSectionType types.MoveSongSectionType,
) UserDataService {
	return &userDataService{
		createGuitarTuning: createGuitarTuning,
		deleteGuitarTuning: deleteGuitarTuning,
		moveGuitarTuning:   moveGuitarTuning,

		createSongSectionType: createSongSectionType,
		deleteSongSectionType: deleteSongSectionType,
		moveSongSectionType:   moveSongSectionType,
	}
}

// Guitar Tunings

func (u *userDataService) CreateGuitarTuning(request requests.CreateGuitarTuningRequest, token string) *wrapper.ErrorCode {
	return u.createGuitarTuning.Handle(request, token)
}

func (u *userDataService) DeleteGuitarTuning(id uuid.UUID, token string) *wrapper.ErrorCode {
	return u.deleteGuitarTuning.Handle(id, token)
}

func (u *userDataService) MoveGuitarTuning(request requests.MoveGuitarTuningRequest, token string) *wrapper.ErrorCode {
	return u.moveGuitarTuning.Handle(request, token)
}

// Song Section Types

func (u *userDataService) CreateSectionType(
	request requests.CreateSongSectionTypeRequest,
	token string,
) *wrapper.ErrorCode {
	return u.createSongSectionType.Handle(request, token)
}

func (u *userDataService) DeleteSectionType(id uuid.UUID, token string) *wrapper.ErrorCode {
	return u.deleteSongSectionType.Handle(id, token)
}

func (u *userDataService) MoveSectionType(request requests.MoveSongSectionTypeRequest, token string) *wrapper.ErrorCode {
	return u.moveSongSectionType.Handle(request, token)
}
