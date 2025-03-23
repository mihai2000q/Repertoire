package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"repertoire/server/data/database"
	"repertoire/server/model"
)

type UserDataRepository interface {
	GetBandMemberRoles(roles *[]model.BandMemberRole, userID uuid.UUID) error
	CountBandMemberRoles(count *int64, userID uuid.UUID) error
	CreateBandMemberRole(bandMember *model.BandMemberRole) error
	UpdateAllBandMemberRoles(bandMemberRoles *[]model.BandMemberRole) error
	DeleteBandMemberRole(id uuid.UUID) error

	GetGuitarTunings(tunings *[]model.GuitarTuning, userID uuid.UUID) error
	GetGuitarTuningsCount(count *int64, userID uuid.UUID) error
	CreateGuitarTuning(tuning *model.GuitarTuning) error
	UpdateAllGuitarTunings(tunings *[]model.GuitarTuning) error
	DeleteGuitarTuning(id uuid.UUID) error

	GetInstruments(instruments *[]model.Instrument, userID uuid.UUID) error
	GetInstrumentsCount(count *int64, userID uuid.UUID) error
	CreateInstrument(instrument *model.Instrument) error
	UpdateAllInstruments(instruments *[]model.Instrument) error
	DeleteInstrument(id uuid.UUID) error

	GetSectionTypes(types *[]model.SongSectionType, userID uuid.UUID) error
	CountSectionTypes(count *int64, userID uuid.UUID) error
	CreateSectionType(sectionType *model.SongSectionType) error
	UpdateAllSectionTypes(sectionTypes *[]model.SongSectionType) error
	DeleteSectionType(id uuid.UUID) error
}

type userDataRepository struct {
	client database.Client
}

func NewUserDataRepository(client database.Client) UserDataRepository {
	return userDataRepository{
		client: client,
	}
}

// Band Member - Roles

func (u userDataRepository) GetBandMemberRoles(bandMemberRoles *[]model.BandMemberRole, userID uuid.UUID) error {
	return u.client.Find(&bandMemberRoles, model.BandMemberRole{UserID: userID}).Error
}

func (u userDataRepository) CountBandMemberRoles(count *int64, userID uuid.UUID) error {
	return u.client.Model(&model.BandMemberRole{}).
		Where(model.BandMemberRole{UserID: userID}).
		Count(count).
		Error
}

func (u userDataRepository) CreateBandMemberRole(bandMemberRole *model.BandMemberRole) error {
	return u.client.Create(&bandMemberRole).Error
}

func (u userDataRepository) UpdateAllBandMemberRoles(bandMemberRoles *[]model.BandMemberRole) error {
	return u.client.Transaction(func(tx *gorm.DB) error {
		for _, sectionType := range *bandMemberRoles {
			if err := tx.Save(sectionType).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (u userDataRepository) DeleteBandMemberRole(id uuid.UUID) error {
	return u.client.Delete(&model.BandMemberRole{}, id).Error
}

// Guitar Tunings

func (u userDataRepository) GetGuitarTunings(tunings *[]model.GuitarTuning, userID uuid.UUID) error {
	return u.client.Model(&model.GuitarTuning{}).
		Where(model.GuitarTuning{UserID: userID}).
		Order("\"order\"").
		Find(&tunings).
		Error
}

func (u userDataRepository) GetGuitarTuningsCount(count *int64, userID uuid.UUID) error {
	return u.client.Model(&model.GuitarTuning{}).
		Where(model.GuitarTuning{UserID: userID}).
		Count(count).
		Error
}

func (u userDataRepository) CreateGuitarTuning(tuning *model.GuitarTuning) error {
	return u.client.Create(&tuning).Error
}

func (u userDataRepository) UpdateAllGuitarTunings(tunings *[]model.GuitarTuning) error {
	return u.client.Transaction(func(tx *gorm.DB) error {
		for _, tuning := range *tunings {
			if err := tx.Save(tuning).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (u userDataRepository) DeleteGuitarTuning(id uuid.UUID) error {
	return u.client.Delete(&model.GuitarTuning{}, id).Error
}

// Instruments

func (u userDataRepository) GetInstruments(instruments *[]model.Instrument, userID uuid.UUID) error {
	return u.client.Model(&model.Instrument{}).
		Where(model.Instrument{UserID: userID}).
		Order("\"order\"").
		Find(&instruments).
		Error
}

func (u userDataRepository) GetInstrumentsCount(count *int64, userID uuid.UUID) error {
	return u.client.Model(&model.Instrument{}).
		Where(model.Instrument{UserID: userID}).
		Count(count).
		Error
}

func (u userDataRepository) CreateInstrument(instrument *model.Instrument) error {
	return u.client.Create(&instrument).Error
}

func (u userDataRepository) UpdateAllInstruments(instruments *[]model.Instrument) error {
	return u.client.Transaction(func(tx *gorm.DB) error {
		for _, tuning := range *instruments {
			if err := tx.Save(tuning).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (u userDataRepository) DeleteInstrument(id uuid.UUID) error {
	return u.client.Delete(&model.Instrument{}, id).Error
}

// Section Types

func (u userDataRepository) GetSectionTypes(types *[]model.SongSectionType, userID uuid.UUID) error {
	return u.client.Model(&model.SongSectionType{}).
		Where(model.SongSectionType{UserID: userID}).
		Order("\"order\"").
		Find(&types).
		Error
}

func (u userDataRepository) CountSectionTypes(count *int64, userID uuid.UUID) error {
	return u.client.Model(&model.SongSectionType{}).
		Where(model.SongSectionType{UserID: userID}).
		Count(count).
		Error
}

func (u userDataRepository) CreateSectionType(sectionType *model.SongSectionType) error {
	return u.client.Create(&sectionType).Error
}

func (u userDataRepository) UpdateAllSectionTypes(sectionTypes *[]model.SongSectionType) error {
	return u.client.Transaction(func(tx *gorm.DB) error {
		for _, sectionType := range *sectionTypes {
			if err := tx.Save(sectionType).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (u userDataRepository) DeleteSectionType(id uuid.UUID) error {
	return u.client.Delete(&model.SongSectionType{}, id).Error
}

// Song Section History

func (u userDataRepository) GetSongSectionHistory(
	history *[]model.SongSectionHistory,
	sectionID uuid.UUID,
	property model.SongSectionProperty,
) error {
	return u.client.
		Order("created_at").
		Find(&history, model.SongSectionHistory{SongSectionID: sectionID, Property: property}).
		Error
}

func (u userDataRepository) CreateSongSectionHistory(history *model.SongSectionHistory) error {
	return u.client.Create(&history).Error
}
