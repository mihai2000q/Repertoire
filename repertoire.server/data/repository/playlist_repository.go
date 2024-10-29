package repository

import (
	"gorm.io/gorm/clause"
	"repertoire/data/database"
	"repertoire/model"

	"github.com/google/uuid"
)

type PlaylistRepository interface {
	Get(playlist *model.Playlist, id uuid.UUID) error
	GetWithAssociations(playlist *model.Playlist, id uuid.UUID) error
	GetAllByUser(
		playlists *[]model.Playlist,
		userID uuid.UUID,
		currentPage *int,
		pageSize *int,
		orderBy string,
	) error
	GetAllByUserCount(count *int64, userID uuid.UUID) error
	Create(playlist *model.Playlist) error
	Update(playlist *model.Playlist) error
	Delete(id uuid.UUID) error
}

type playlistRepository struct {
	client database.Client
}

func NewPlaylistRepository(client database.Client) PlaylistRepository {
	return playlistRepository{
		client: client,
	}
}

func (p playlistRepository) Get(playlist *model.Playlist, id uuid.UUID) error {
	return p.client.DB.Find(&playlist, model.Playlist{ID: id}).Error
}

func (p playlistRepository) GetWithAssociations(playlist *model.Playlist, id uuid.UUID) error {
	return p.client.DB.Preload(clause.Associations).Find(&playlist, model.Playlist{ID: id}).Error
}

func (p playlistRepository) GetAllByUser(
	playlists *[]model.Playlist,
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
	return p.client.DB.Model(&model.Playlist{}).
		Preload(clause.Associations).
		Where(model.Playlist{UserID: userID}).
		Order(orderBy).
		Offset(offset).
		Limit(*pageSize).
		Find(&playlists).
		Error
}

func (p playlistRepository) GetAllByUserCount(count *int64, userID uuid.UUID) error {
	return p.client.DB.Model(&model.Playlist{}).
		Where(model.Playlist{UserID: userID}).
		Count(count).
		Error
}

func (p playlistRepository) Create(playlist *model.Playlist) error {
	return p.client.DB.Create(&playlist).Error
}

func (p playlistRepository) Update(playlist *model.Playlist) error {
	return p.client.DB.Save(&playlist).Error
}

func (p playlistRepository) Delete(id uuid.UUID) error {
	return p.client.DB.Delete(&model.Playlist{}, id).Error
}
