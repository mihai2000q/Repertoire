package repository

import (
	"repertoire/server/data/database"
	"repertoire/server/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SongSectionRepository interface {
	Get(section *model.SongSection, id uuid.UUID) error
	GetWithSectionParts(section *model.SongSection, id uuid.UUID) error
	GetAllByIDs(sections *[]model.SongSection, ids []uuid.UUID) error
	GetAllByIDsWithSectionParts(sections *[]model.SongSection, ids []uuid.UUID) error
	GetAllByPartWithSectionParts(sections *[]model.SongSection, partID uuid.UUID) error
	GetAllByPartIDsWithSectionParts(sections *[]model.SongSection, partIDs []uuid.UUID) error
	CountAllBySong(count *int64, songID uuid.UUID) error
	Create(section *model.SongSection) error
	Update(section *model.SongSection) error
	UpdateWithAssociations(section *model.SongSection) error
	UpdateAllWithAssociations(sections *[]model.SongSection) error
	Delete(ids []uuid.UUID) error

	CreateAllSectionParts(sectionParts *[]model.SongSectionPart) error
	UpdateAllSectionParts(sectionParts *[]model.SongSectionPart) error
	DeleteSectionParts(sectionParts *[]model.SongSectionPart) error

	GetTypes(types *[]model.SongSectionType, userID uuid.UUID) error
}

type songSectionRepository struct {
	client database.Client
}

func NewSongSectionRepository(client database.Client) SongSectionRepository {
	return songSectionRepository{
		client: client,
	}
}

func (s songSectionRepository) Get(section *model.SongSection, id uuid.UUID) error {
	return s.client.Find(section, model.SongSection{ID: id}).Error
}
func (s songSectionRepository) GetWithSectionParts(section *model.SongSection, id uuid.UUID) error {
	return s.client.
		Preload("SectionParts", func(db *gorm.DB) *gorm.DB {
			return db.Order("\"order\"")
		}).
		Find(section, model.SongSection{ID: id}).
		Error
}

func (s songSectionRepository) GetAllByIDs(sections *[]model.SongSection, ids []uuid.UUID) error {
	return s.client.Find(sections, ids).Error
}

func (s songSectionRepository) GetAllByIDsWithSectionParts(sections *[]model.SongSection, ids []uuid.UUID) error {
	return s.client.
		Preload("SectionParts", func(db *gorm.DB) *gorm.DB {
			return db.Order("\"order\"")
		}).
		Find(sections, ids).
		Error
}

func (s songSectionRepository) GetAllByPartWithSectionParts(sections *[]model.SongSection, partID uuid.UUID) error {
	return s.client.
		Joins("LEFT JOIN song_section_parts ON song_sections.id = song_section_parts.section_id").
		Where("song_section_parts.part_id = ?", partID).
		Preload("SectionParts", func(db *gorm.DB) *gorm.DB {
			return db.Order("song_section_parts.order")
		}).
		Find(sections).
		Error
}

func (s songSectionRepository) GetAllByPartIDsWithSectionParts(sections *[]model.SongSection, partIDs []uuid.UUID) error {
	return s.client.
		Joins("LEFT JOIN song_section_parts ON song_sections.id = song_section_parts.section_id").
		Where("song_section_parts.part_id IN ?", partIDs).
		Preload("SectionParts", func(db *gorm.DB) *gorm.DB {
			return db.Order("song_section_parts.order")
		}).
		Find(sections).
		Error
}

func (s songSectionRepository) CountAllBySong(count *int64, songID uuid.UUID) error {
	return s.client.Model(&model.SongSection{}).
		Where(model.SongSection{SongID: songID}).
		Count(count).
		Error
}

func (s songSectionRepository) Create(section *model.SongSection) error {
	return s.client.Create(section).Error
}

func (s songSectionRepository) Update(section *model.SongSection) error {
	return s.client.Save(section).Error
}

func (s songSectionRepository) UpdateWithAssociations(section *model.SongSection) error {
	return s.client.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Save(section).
		Error
}

func (s songSectionRepository) UpdateAllWithAssociations(sections *[]model.SongSection) error {
	return s.client.Transaction(func(tx *gorm.DB) error {
		for _, song := range *sections {
			err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(song).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s songSectionRepository) Delete(ids []uuid.UUID) error {
	return s.client.Delete(&model.SongSection{}, ids).Error
}

// Section Parts

func (s songSectionRepository) CreateAllSectionParts(sectionParts *[]model.SongSectionPart) error {
	return s.client.Transaction(func(tx *gorm.DB) error {
		for _, sectionPart := range *sectionParts {
			if err := tx.Create(sectionPart).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s songSectionRepository) UpdateAllSectionParts(sectionParts *[]model.SongSectionPart) error {
	return s.client.Transaction(func(tx *gorm.DB) error {
		for _, sectionPart := range *sectionParts {
			if err := tx.Save(sectionPart).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s songSectionRepository) DeleteSectionParts(sectionParts *[]model.SongSectionPart) error {
	return s.client.Delete(sectionParts).Error
}

// Types

func (s songSectionRepository) GetTypes(types *[]model.SongSectionType, userID uuid.UUID) error {
	return s.client.Model(&model.SongSectionType{}).
		Where(model.SongSectionType{UserID: userID}).
		Order("\"order\"").
		Find(types).
		Error
}
