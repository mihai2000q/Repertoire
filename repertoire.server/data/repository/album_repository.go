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
	GetWithAssociations(album *model.Album, id uuid.UUID) error
	GetAllByUser(
		albums *[]model.Album,
		userID uuid.UUID,
		currentPage *int,
		pageSize *int,
		orderBy string,
	) error
	GetAllByUserCount(count *int64, userID uuid.UUID) error
	Create(album *model.Album) error
	Update(album *model.Album) error
	Delete(id uuid.UUID) error
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
	orderBy string,
) error {
	offset := -1
	if pageSize == nil {
		pageSize = &[]int{-1}[0]
	} else {
		offset = (*currentPage - 1) * *pageSize
	}
	return a.client.DB.Model(&model.Album{}).
		Preload(clause.Associations).
		Where(model.Album{UserID: userID}).
		Order(orderBy).
		Offset(offset).
		Limit(*pageSize).
		Find(&albums).
		Error
}

func (a albumRepository) GetAllByUserCount(count *int64, userID uuid.UUID) error {
	return a.client.DB.Model(&model.Album{}).
		Where(model.Album{UserID: userID}).
		Count(count).
		Error
}

func (a albumRepository) Create(album *model.Album) error {
	return a.client.DB.Create(&album).Error
}

func (a albumRepository) Update(album *model.Album) error {
	return a.client.DB.Save(&album).Error
}

func (a albumRepository) Delete(id uuid.UUID) error {
	return a.client.DB.Delete(&model.Album{}, id).Error
}
