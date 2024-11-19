package repository

import (
	"repertoire/server/data/database"
	"repertoire/server/model"

	"gorm.io/gorm/clause"

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
		orderBy []string,
		searchBy []string,
	) error
	GetAllByUserCount(count *int64, userID uuid.UUID, searchBy []string) error
	CountSongs(count *int64, id uuid.UUID) error
	Create(playlist *model.Playlist) error
	AddSong(playlist *model.Playlist, song *model.Song) error
	Update(playlist *model.Playlist) error
	Delete(id uuid.UUID) error
	RemoveSong(playlist *model.Playlist, song *model.Song) error
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
	orderBy []string,
	searchBy []string,
) error {
	tx := p.client.DB.Model(&model.Playlist{}).
		Preload("Songs").
		Where(model.Playlist{UserID: userID})
	tx = database.SearchBy(tx, searchBy)
	tx = database.OrderBy(tx, orderBy)
	tx = database.Paginate(tx, currentPage, pageSize)
	return tx.Find(&playlists).Error
}

func (p playlistRepository) GetAllByUserCount(count *int64, userID uuid.UUID, searchBy []string) error {
	tx := p.client.DB.Model(&model.Playlist{}).
		Where(model.Playlist{UserID: userID})
	tx = database.SearchBy(tx, searchBy)
	return tx.Count(count).Error
}

func (p playlistRepository) CountSongs(count *int64, id uuid.UUID) error {
	return p.client.DB.Model(&model.Song{}).
		Preload("Playlists").
		Where("playlist_song.playlist_id = ?", id).
		Count(count).
		Error
}

func (p playlistRepository) Create(playlist *model.Playlist) error {
	return p.client.DB.Create(&playlist).Error
}

func (p playlistRepository) AddSong(playlist *model.Playlist, song *model.Song) error {
	return p.client.DB.Model(&playlist).Association("Songs").Append(&song)
}

func (p playlistRepository) Update(playlist *model.Playlist) error {
	return p.client.DB.Save(&playlist).Error
}

func (p playlistRepository) Delete(id uuid.UUID) error {
	return p.client.DB.Delete(&model.Playlist{}, id).Error
}

func (p playlistRepository) RemoveSong(playlist *model.Playlist, song *model.Song) error {
	return p.client.DB.Model(&playlist).Association("Songs").Delete(&song)
}
