package service

import (
	"github.com/google/uuid"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/udata/band/member/role"
	"repertoire/server/domain/usecase/udata/guitar/tuning"
	"repertoire/server/domain/usecase/udata/instrument"
	"repertoire/server/domain/usecase/udata/section/types"
	"repertoire/server/internal/wrapper"
)

type UserDataService interface {
	CreateBandMemberRole(request requests.CreateBandMemberRoleRequest, token string) *wrapper.ErrorCode
	DeleteBandMemberRole(id uuid.UUID, token string) *wrapper.ErrorCode
	MoveBandMemberRole(request requests.MoveBandMemberRoleRequest, token string) *wrapper.ErrorCode

	CreateInstrument(request requests.CreateInstrumentRequest, token string) *wrapper.ErrorCode
	MoveInstrument(request requests.MoveInstrumentRequest, token string) *wrapper.ErrorCode
	DeleteInstrument(id uuid.UUID, token string) *wrapper.ErrorCode

	CreateGuitarTuning(request requests.CreateGuitarTuningRequest, token string) *wrapper.ErrorCode
	MoveGuitarTuning(request requests.MoveGuitarTuningRequest, token string) *wrapper.ErrorCode
	DeleteGuitarTuning(id uuid.UUID, token string) *wrapper.ErrorCode

	CreateSectionType(request requests.CreateSongSectionTypeRequest, token string) *wrapper.ErrorCode
	DeleteSectionType(id uuid.UUID, token string) *wrapper.ErrorCode
	MoveSectionType(request requests.MoveSongSectionTypeRequest, token string) *wrapper.ErrorCode
}

type userDataService struct {
	createBandMemberRole role.CreateBandMemberRole
	deleteBandMemberRole role.DeleteBandMemberRole
	moveBandMemberRole   role.MoveBandMemberRole

	createInstrument instrument.CreateInstrument
	deleteInstrument instrument.DeleteInstrument
	moveInstrument   instrument.MoveInstrument

	createGuitarTuning tuning.CreateGuitarTuning
	deleteGuitarTuning tuning.DeleteGuitarTuning
	moveGuitarTuning   tuning.MoveGuitarTuning

	createSongSectionType types.CreateSongSectionType
	deleteSongSectionType types.DeleteSongSectionType
	moveSongSectionType   types.MoveSongSectionType
}

func NewUserDataService(
	createBandMemberRole role.CreateBandMemberRole,
	deleteBandMemberRole role.DeleteBandMemberRole,
	moveBandMemberRole role.MoveBandMemberRole,

	createInstrument instrument.CreateInstrument,
	deleteInstrument instrument.DeleteInstrument,
	moveInstrument instrument.MoveInstrument,

	createGuitarTuning tuning.CreateGuitarTuning,
	deleteGuitarTuning tuning.DeleteGuitarTuning,
	moveGuitarTuning tuning.MoveGuitarTuning,

	createSongSectionType types.CreateSongSectionType,
	deleteSongSectionType types.DeleteSongSectionType,
	moveSongSectionType types.MoveSongSectionType,
) UserDataService {
	return &userDataService{
		createBandMemberRole: createBandMemberRole,
		deleteBandMemberRole: deleteBandMemberRole,
		moveBandMemberRole:   moveBandMemberRole,

		createInstrument: createInstrument,
		deleteInstrument: deleteInstrument,
		moveInstrument:   moveInstrument,

		createGuitarTuning: createGuitarTuning,
		deleteGuitarTuning: deleteGuitarTuning,
		moveGuitarTuning:   moveGuitarTuning,

		createSongSectionType: createSongSectionType,
		deleteSongSectionType: deleteSongSectionType,
		moveSongSectionType:   moveSongSectionType,
	}
}

// Band Member Roles

func (u *userDataService) CreateBandMemberRole(request requests.CreateBandMemberRoleRequest, token string) *wrapper.ErrorCode {
	return u.createBandMemberRole.Handle(request, token)
}

func (u *userDataService) DeleteBandMemberRole(id uuid.UUID, token string) *wrapper.ErrorCode {
	return u.deleteBandMemberRole.Handle(id, token)
}

func (u *userDataService) MoveBandMemberRole(request requests.MoveBandMemberRoleRequest, token string) *wrapper.ErrorCode {
	return u.moveBandMemberRole.Handle(request, token)
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

// Instruments

func (u *userDataService) CreateInstrument(request requests.CreateInstrumentRequest, token string) *wrapper.ErrorCode {
	return u.createInstrument.Handle(request, token)
}

func (u *userDataService) DeleteInstrument(id uuid.UUID, token string) *wrapper.ErrorCode {
	return u.deleteInstrument.Handle(id, token)
}

func (u *userDataService) MoveInstrument(request requests.MoveInstrumentRequest, token string) *wrapper.ErrorCode {
	return u.moveInstrument.Handle(request, token)
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
