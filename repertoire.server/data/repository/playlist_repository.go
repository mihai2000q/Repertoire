package repository

import (
	"repertoire/data/database"
	"repertoire/models"

	"github.com/google/uuid"
)

type PlaylistRepository interface {
	Get(playlist *models.Playlist, id uuid.UUID) error
	GetAllByUser(playlists *[]models.Playlist, userId uuid.UUID, currentPage *int, pageSize *int) error
	GetAllByUserCount(count *int64, userId uuid.UUID) error
	Create(playlist *models.Playlist) error
	Update(playlist *models.Playlist) error
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

func (p playlistRepository) Get(playlist *models.Playlist, id uuid.UUID) error {
	return p.client.DB.Find(&playlist, models.Playlist{ID: id}).Error
}

func (p playlistRepository) GetAllByUser(
	playlists *[]models.Playlist,
	userId uuid.UUID,
	currentPage *int,
	pageSize *int,
) error {
	if currentPage == nil {
		currentPage = &[]int{1}[0]
	}
	if pageSize == nil {
		pageSize = &[]int{1}[0]
	}
	return p.client.DB.Model(&models.Playlist{}).
		Where(models.Playlist{UserID: userId}).
		Offset((*currentPage - 1) * *pageSize).
		Limit(*pageSize).
		Find(&playlists).
		Error
}

func (p playlistRepository) GetAllByUserCount(count *int64, userId uuid.UUID) error {
	return p.client.DB.Model(&models.Playlist{}).
		Where(models.Playlist{UserID: userId}).
		Count(count).
		Error
}

func (p playlistRepository) Create(playlist *models.Playlist) error {
	return p.client.DB.Create(&playlist).Error
}

func (p playlistRepository) Update(playlist *models.Playlist) error {
	return p.client.DB.Save(&playlist).Error
}

func (p playlistRepository) Delete(id uuid.UUID) error {
	return p.client.DB.Delete(&models.Playlist{}, id).Error
}
