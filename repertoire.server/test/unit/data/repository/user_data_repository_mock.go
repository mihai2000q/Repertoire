package repository

import (
	"repertoire/server/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type UserDataRepositoryMock struct {
	mock.Mock
}

// Band Member Roles

func (u *UserDataRepositoryMock) GetBandMemberRoles(bandMemberRoles *[]model.BandMemberRole, userID uuid.UUID) error {
	args := u.Called(bandMemberRoles, userID)

	if len(args) > 1 {
		*bandMemberRoles = *args.Get(1).(*[]model.BandMemberRole)
	}

	return args.Error(0)
}

func (u *UserDataRepositoryMock) CountBandMemberRoles(count *int64, userID uuid.UUID) error {
	args := u.Called(count, userID)

	if len(args) > 1 {
		*count = *args.Get(1).(*int64)
	}

	return args.Error(0)
}

func (u *UserDataRepositoryMock) CreateBandMemberRole(bandMemberRole *model.BandMemberRole) error {
	args := u.Called(bandMemberRole)
	return args.Error(0)
}

func (u *UserDataRepositoryMock) UpdateAllBandMemberRoles(bandMemberRoles *[]model.BandMemberRole) error {
	args := u.Called(bandMemberRoles)
	return args.Error(0)
}

func (u *UserDataRepositoryMock) DeleteBandMemberRole(id uuid.UUID) error {
	args := u.Called(id)
	return args.Error(0)
}

// Guitar Tunings

func (u *UserDataRepositoryMock) GetGuitarTunings(tunings *[]model.GuitarTuning, userID uuid.UUID) error {
	args := u.Called(tunings, userID)

	if len(args) > 1 {
		*tunings = *args.Get(1).(*[]model.GuitarTuning)
	}

	return args.Error(0)
}

func (u *UserDataRepositoryMock) GetGuitarTuningsCount(count *int64, userID uuid.UUID) error {
	args := u.Called(count, userID)

	if len(args) > 1 {
		*count = *args.Get(1).(*int64)
	}

	return args.Error(0)
}

func (u *UserDataRepositoryMock) CreateGuitarTuning(tuning *model.GuitarTuning) error {
	args := u.Called(tuning)
	return args.Error(0)
}

func (u *UserDataRepositoryMock) UpdateAllGuitarTunings(tunings *[]model.GuitarTuning) error {
	args := u.Called(tunings)
	return args.Error(0)
}

func (u *UserDataRepositoryMock) DeleteGuitarTuning(id uuid.UUID) error {
	args := u.Called(id)
	return args.Error(0)
}

// Instruments

func (u *UserDataRepositoryMock) GetInstruments(instruments *[]model.Instrument, userID uuid.UUID) error {
	args := u.Called(instruments, userID)

	if len(args) > 1 {
		*instruments = *args.Get(1).(*[]model.Instrument)
	}

	return args.Error(0)
}

func (u *UserDataRepositoryMock) GetInstrumentsCount(count *int64, userID uuid.UUID) error {
	args := u.Called(count, userID)

	if len(args) > 1 {
		*count = *args.Get(1).(*int64)
	}

	return args.Error(0)
}

func (u *UserDataRepositoryMock) CreateInstrument(instrument *model.Instrument) error {
	args := u.Called(instrument)
	return args.Error(0)
}

func (u *UserDataRepositoryMock) UpdateAllInstruments(instruments *[]model.Instrument) error {
	args := u.Called(instruments)
	return args.Error(0)
}

func (u *UserDataRepositoryMock) DeleteInstrument(id uuid.UUID) error {
	args := u.Called(id)
	return args.Error(0)
}

// Section Types

func (u *UserDataRepositoryMock) GetSectionTypes(sectionTypes *[]model.SongSectionType, userID uuid.UUID) error {
	args := u.Called(sectionTypes, userID)

	if len(args) > 1 {
		*sectionTypes = *args.Get(1).(*[]model.SongSectionType)
	}

	return args.Error(0)
}

func (u *UserDataRepositoryMock) CountSectionTypes(count *int64, userID uuid.UUID) error {
	args := u.Called(count, userID)

	if len(args) > 1 {
		*count = *args.Get(1).(*int64)
	}

	return args.Error(0)
}

func (u *UserDataRepositoryMock) CreateSectionType(sectionType *model.SongSectionType) error {
	args := u.Called(sectionType)
	return args.Error(0)
}

func (u *UserDataRepositoryMock) UpdateAllSectionTypes(sectionTypes *[]model.SongSectionType) error {
	args := u.Called(sectionTypes)
	return args.Error(0)
}

func (u *UserDataRepositoryMock) DeleteSectionType(id uuid.UUID) error {
	args := u.Called(id)
	return args.Error(0)
}
