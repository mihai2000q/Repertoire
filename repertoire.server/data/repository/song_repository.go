package repository

import (
	"repertoire/server/data/database"
	"repertoire/server/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/google/uuid"
)

type SongRepository interface {
	Get(song *model.Song, id uuid.UUID) error
	GetWithSections(song *model.Song, id uuid.UUID) error
	GetWithSongs(song *model.Song, id uuid.UUID) error
	GetWithAssociations(song *model.Song, id uuid.UUID) error
	GetAllByUser(
		songs *[]model.Song,
		userID uuid.UUID,
		currentPage *int,
		pageSize *int,
		orderBy []string,
		searchBy []string,
	) error
	GetAllByUserCount(count *int64, userID uuid.UUID, searchBy []string) error
	GetGuitarTunings(tunings *[]model.GuitarTuning, userID uuid.UUID) error
	Create(song *model.Song) error
	Update(song *model.Song) error
	UpdateWithAssociations(song *model.Song) error
	Delete(id uuid.UUID) error

	GetSection(section *model.SongSection, id uuid.UUID) error
	GetSectionTypes(types *[]model.SongSectionType, userID uuid.UUID) error
	CountSectionsBySong(count *int64, songID uuid.UUID) error
	CreateSection(section *model.SongSection) error
	UpdateSection(section *model.SongSection) error
	DeleteSection(id uuid.UUID) error
}

type songRepository struct {
	client database.Client
}

func NewSongRepository(client database.Client) SongRepository {
	return songRepository{
		client: client,
	}
}

func (s songRepository) Get(song *model.Song, id uuid.UUID) error {
	return s.client.DB.Find(&song, model.Song{ID: id}).Error
}

func (s songRepository) GetWithSections(song *model.Song, id uuid.UUID) error {
	return s.client.DB.
		Preload("Sections", func(db *gorm.DB) *gorm.DB {
			return db.Order("song_sections.order")
		}).
		Find(&song, model.Song{ID: id}).
		Error
}

func (s songRepository) GetWithSongs(song *model.Song, id uuid.UUID) error {
	return s.client.DB.
		Preload("Album").
		Preload("Album.Songs").
		Find(&song, model.Song{ID: id}).
		Error
}

func (s songRepository) GetWithAssociations(song *model.Song, id uuid.UUID) error {
	return s.client.DB.
		Preload("Sections", func(db *gorm.DB) *gorm.DB {
			return db.Order("song_sections.order")
		}).
		Preload("Sections.SongSectionType").
		Preload(clause.Associations).
		Find(&song, model.Song{ID: id}).
		Error
}

func (s songRepository) GetAllByUser(
	songs *[]model.Song,
	userID uuid.UUID,
	currentPage *int,
	pageSize *int,
	orderBy []string,
	searchBy []string,
) error {
	tx := s.client.DB.Model(&model.Song{}).
		Preload(clause.Associations).
		Where(model.Song{UserID: userID})
	tx = database.SearchBy(tx, searchBy)
	tx = database.OrderBy(tx, orderBy)
	tx = database.Paginate(tx, currentPage, pageSize)
	return tx.Find(&songs).Error
}

func (s songRepository) GetAllByUserCount(count *int64, userID uuid.UUID, searchBy []string) error {
	tx := s.client.DB.Model(&model.Song{}).
		Where(model.Song{UserID: userID})
	tx = database.SearchBy(tx, searchBy)
	return tx.Count(count).Error
}

func (s songRepository) GetGuitarTunings(tunings *[]model.GuitarTuning, userID uuid.UUID) error {
	return s.client.DB.Model(&model.GuitarTuning{}).
		Where(model.GuitarTuning{UserID: userID}).
		Order("\"order\"").
		Find(&tunings).
		Error
}

func (s songRepository) Create(song *model.Song) error {
	return s.client.DB.Create(&song).Error
}

func (s songRepository) Update(song *model.Song) error {
	return s.client.DB.Save(&song).Error
}

func (s songRepository) UpdateWithAssociations(song *model.Song) error {
	return s.client.DB.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Updates(&song).
		Error
}

func (s songRepository) Delete(id uuid.UUID) error {
	return s.client.DB.Delete(&model.Song{}, id).Error
}

// Sections

func (s songRepository) GetSection(section *model.SongSection, id uuid.UUID) error {
	return s.client.DB.Find(&section, model.SongSection{ID: id}).Error
}

func (s songRepository) GetSectionTypes(types *[]model.SongSectionType, userID uuid.UUID) error {
	return s.client.DB.Model(&model.SongSectionType{}).
		Where(model.SongSectionType{UserID: userID}).
		Order("\"order\"").
		Find(&types).
		Error
}

func (s songRepository) CountSectionsBySong(count *int64, songID uuid.UUID) error {
	return s.client.DB.Model(&model.SongSection{}).
		Where(model.SongSection{SongID: songID}).
		Count(count).
		Error
}

func (s songRepository) CreateSection(section *model.SongSection) error {
	return s.client.DB.Create(&section).Error
}

func (s songRepository) UpdateSection(section *model.SongSection) error {
	return s.client.DB.Save(&section).Error
}

func (s songRepository) DeleteSection(id uuid.UUID) error {
	return s.client.DB.Delete(&model.SongSection{}, id).Error
}
