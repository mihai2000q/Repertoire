package repository

import (
	"repertoire/server/data/database"
	"repertoire/server/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/google/uuid"
)

type AlbumRepository interface {
	Get(album *model.Album, id uuid.UUID) error
	GetWithSongs(album *model.Album, id uuid.UUID) error
	GetWithAssociations(album *model.Album, id uuid.UUID) error
	GetAllByUser(
		albums *[]model.Album,
		userID uuid.UUID,
		currentPage *int,
		pageSize *int,
		orderBy []string,
		searchBy []string,
	) error
	GetAllByUserCount(count *int64, userID uuid.UUID, searchBy []string) error
	Create(album *model.Album) error
	Update(album *model.Album) error
	UpdateWithAssociations(album *model.Album) error
	Delete(id uuid.UUID) error
	RemoveSong(album *model.Album, song *model.Song) error
}

type albumRepository struct {
	client database.Client
}

func NewAlbumRepository(client database.Client) AlbumRepository {
	return albumRepository{
		client: client,
	}
}

func (a albumRepository) Get(album *model.Album, id uuid.UUID) error {
	return a.client.DB.Find(&album, model.Album{ID: id}).Error
}

func (a albumRepository) GetWithSongs(album *model.Album, id uuid.UUID) error {
	return a.client.DB.
		Preload("Songs", func(db *gorm.DB) *gorm.DB {
			return db.Order("songs.track_no")
		}).
		Find(&album, model.Album{ID: id}).
		Error
}

func (a albumRepository) GetWithAssociations(album *model.Album, id uuid.UUID) error {
	return a.client.DB.
		Preload("Songs", func(db *gorm.DB) *gorm.DB {
			return db.Order("songs.track_no")
		}).
		Preload(clause.Associations).
		Find(&album, model.Album{ID: id}).
		Error
}

func (a albumRepository) GetAllByUser(
	albums *[]model.Album,
	userID uuid.UUID,
	currentPage *int,
	pageSize *int,
	orderBy []string,
	searchBy []string,
) error {
	tx := a.client.DB.Model(&model.Album{}).
		Preload(clause.Associations).
		Where(model.Album{UserID: userID})
	tx = database.SearchBy(tx, searchBy)
	tx = database.OrderBy(tx, orderBy)
	tx = database.Paginate(tx, currentPage, pageSize)
	return tx.Find(&albums).Error
}

func (a albumRepository) GetAllByUserCount(count *int64, userID uuid.UUID, searchBy []string) error {
	tx := a.client.DB.Model(&model.Album{}).
		Where(model.Album{UserID: userID})
	tx = database.SearchBy(tx, searchBy)
	return tx.Count(count).Error
}

func (a albumRepository) Create(album *model.Album) error {
	return a.client.DB.Create(&album).Error
}

func (a albumRepository) Update(album *model.Album) error {
	return a.client.DB.Save(&album).Error
}

func (a albumRepository) UpdateWithAssociations(album *model.Album) error {
	return a.client.DB.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Updates(&album).
		Error
}

func (a albumRepository) Delete(id uuid.UUID) error {
	return a.client.DB.Delete(&model.Album{}, id).Error
}

func (a albumRepository) RemoveSong(album *model.Album, song *model.Song) error {
	return a.client.DB.Model(&album).Association("Songs").Delete(&song)
}
