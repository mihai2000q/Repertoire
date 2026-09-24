package service

import (
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/userdata/band/member/role"
	"repertoire/server/domain/usecase/userdata/guitartuning"
	"repertoire/server/domain/usecase/userdata/instrument"
	"repertoire/server/domain/usecase/userdata/section/types"
	"repertoire/server/internal/httperror"

	"github.com/google/uuid"
)

type UserDataService interface {
	CreateBandMemberRole(request requests.CreateBandMemberRoleRequest, token string) *httperror.ErrorCode
	DeleteBandMemberRole(id uuid.UUID, token string) *httperror.ErrorCode
	MoveBandMemberRole(request requests.MoveBandMemberRoleRequest, token string) *httperror.ErrorCode

	CreateInstrument(request requests.CreateInstrumentRequest, token string) *httperror.ErrorCode
	MoveInstrument(request requests.MoveInstrumentRequest, token string) *httperror.ErrorCode
	DeleteInstrument(id uuid.UUID, token string) *httperror.ErrorCode

	CreateGuitarTuning(request requests.CreateGuitarTuningRequest, token string) *httperror.ErrorCode
	MoveGuitarTuning(request requests.MoveGuitarTuningRequest, token string) *httperror.ErrorCode
	DeleteGuitarTuning(id uuid.UUID, token string) *httperror.ErrorCode

	CreateSectionType(request requests.CreateSongSectionTypeRequest, token string) *httperror.ErrorCode
	DeleteSectionType(id uuid.UUID, token string) *httperror.ErrorCode
	MoveSectionType(request requests.MoveSongSectionTypeRequest, token string) *httperror.ErrorCode
}

type userDataService struct {
	createBandMemberRole role.CreateBandMemberRole
	deleteBandMemberRole role.DeleteBandMemberRole
	moveBandMemberRole   role.MoveBandMemberRole

	createInstrument instrument.CreateInstrument
	deleteInstrument instrument.DeleteInstrument
	moveInstrument   instrument.MoveInstrument

	createGuitarTuning guitartuning.CreateGuitarTuning
	deleteGuitarTuning guitartuning.DeleteGuitarTuning
	moveGuitarTuning   guitartuning.MoveGuitarTuning

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

	createGuitarTuning guitartuning.CreateGuitarTuning,
	deleteGuitarTuning guitartuning.DeleteGuitarTuning,
	moveGuitarTuning guitartuning.MoveGuitarTuning,

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

func (u *userDataService) CreateBandMemberRole(request requests.CreateBandMemberRoleRequest, token string) *httperror.ErrorCode {
	return u.createBandMemberRole.Handle(request, token)
}

func (u *userDataService) DeleteBandMemberRole(id uuid.UUID, token string) *httperror.ErrorCode {
	return u.deleteBandMemberRole.Handle(id, token)
}

func (u *userDataService) MoveBandMemberRole(request requests.MoveBandMemberRoleRequest, token string) *httperror.ErrorCode {
	return u.moveBandMemberRole.Handle(request, token)
}

// Guitar Tunings

func (u *userDataService) CreateGuitarTuning(request requests.CreateGuitarTuningRequest, token string) *httperror.ErrorCode {
	return u.createGuitarTuning.Handle(request, token)
}

func (u *userDataService) DeleteGuitarTuning(id uuid.UUID, token string) *httperror.ErrorCode {
	return u.deleteGuitarTuning.Handle(id, token)
}

func (u *userDataService) MoveGuitarTuning(request requests.MoveGuitarTuningRequest, token string) *httperror.ErrorCode {
	return u.moveGuitarTuning.Handle(request, token)
}

// Instruments

func (u *userDataService) CreateInstrument(request requests.CreateInstrumentRequest, token string) *httperror.ErrorCode {
	return u.createInstrument.Handle(request, token)
}

func (u *userDataService) DeleteInstrument(id uuid.UUID, token string) *httperror.ErrorCode {
	return u.deleteInstrument.Handle(id, token)
}

func (u *userDataService) MoveInstrument(request requests.MoveInstrumentRequest, token string) *httperror.ErrorCode {
	return u.moveInstrument.Handle(request, token)
}

// Song Section Types

func (u *userDataService) CreateSectionType(
	request requests.CreateSongSectionTypeRequest,
	token string,
) *httperror.ErrorCode {
	return u.createSongSectionType.Handle(request, token)
}

func (u *userDataService) DeleteSectionType(id uuid.UUID, token string) *httperror.ErrorCode {
	return u.deleteSongSectionType.Handle(id, token)
}

func (u *userDataService) MoveSectionType(request requests.MoveSongSectionTypeRequest, token string) *httperror.ErrorCode {
	return u.moveSongSectionType.Handle(request, token)
}
