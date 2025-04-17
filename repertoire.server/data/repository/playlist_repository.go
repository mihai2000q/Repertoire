package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"repertoire/server/data/database"
	"repertoire/server/model"
)

type PlaylistRepository interface {
	Get(playlist *model.Playlist, id uuid.UUID) error
	GetPlaylistSongs(playlistSongs *[]model.PlaylistSong, id uuid.UUID) error
	GetWithAssociations(playlist *model.Playlist, id uuid.UUID) error
	GetAllByUser(
		playlists *[]model.EnhancedPlaylist,
		userID uuid.UUID,
		currentPage *int,
		pageSize *int,
		orderBy []string,
		searchBy []string,
	) error
	GetAllByUserCount(count *int64, userID uuid.UUID, searchBy []string) error
	CountSongs(count *int64, id uuid.UUID) error
	Create(playlist *model.Playlist) error
	AddSongs(playlistSongs *[]model.PlaylistSong) error
	Update(playlist *model.Playlist) error
	UpdateAllPlaylistSongs(playlistSongs *[]model.PlaylistSong) error
	Delete(id uuid.UUID) error
	RemoveSongs(playlistSongs *[]model.PlaylistSong) error
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
	return p.client.Find(&playlist, model.Playlist{ID: id}).Error
}

func (p playlistRepository) GetPlaylistSongs(playlistSongs *[]model.PlaylistSong, id uuid.UUID) error {
	return p.client.
		Order("song_track_no").
		Find(&playlistSongs, model.PlaylistSong{PlaylistID: id}).Error
}

func (p playlistRepository) GetWithAssociations(playlist *model.Playlist, id uuid.UUID) error {
	return p.client.
		Preload("PlaylistSongs", func(db *gorm.DB) *gorm.DB {
			return db.
				Preload("Song").
				Preload("Song.Artist").
				Preload("Song.Album").
				Order("song_track_no")
		}).
		Find(&playlist, model.Playlist{ID: id}).Error
}

func (p playlistRepository) GetAllByUser(
	playlists *[]model.EnhancedPlaylist,
	userID uuid.UUID,
	currentPage *int,
	pageSize *int,
	orderBy []string,
	searchBy []string,
) error {
	tx := p.client.Model(&model.Playlist{}).
		Select(
			"playlists.*",
			"COALESCE(es.songs_count, 0) AS songs_count",
			" es.songs_ids as songs_ids",
		).
		Joins("LEFT JOIN (?) AS es ON es.playlist_id = playlists.id", p.getSongsByPlaylistQuery(userID)).
		Where(model.Playlist{UserID: userID})
	tx = database.SearchBy(tx, searchBy)
	tx = database.OrderBy(tx, orderBy)
	tx = database.Paginate(tx, currentPage, pageSize)
	return tx.Find(&playlists).Error
}

func (p playlistRepository) GetAllByUserCount(count *int64, userID uuid.UUID, searchBy []string) error {
	tx := p.client.Model(&model.Playlist{}).
		Joins("LEFT JOIN (?) AS es ON es.playlist_id = playlists.id", p.getSongsByPlaylistQuery(userID)).
		Where(model.Playlist{UserID: userID})

	tx = database.SearchBy(tx, searchBy)
	return tx.Count(count).Error
}

func (p playlistRepository) CountSongs(count *int64, id uuid.UUID) error {
	return p.client.Model(&model.PlaylistSong{}).
		Where("playlist_id = ?", id).
		Count(count).
		Error
}

func (p playlistRepository) Create(playlist *model.Playlist) error {
	return p.client.Create(&playlist).Error
}

func (p playlistRepository) AddSongs(playlistSongs *[]model.PlaylistSong) error {
	return p.client.Create(&playlistSongs).Error
}

func (p playlistRepository) Update(playlist *model.Playlist) error {
	return p.client.Save(&playlist).Error
}

func (p playlistRepository) UpdateAllPlaylistSongs(playlistSongs *[]model.PlaylistSong) error {
	return p.client.Transaction(func(tx *gorm.DB) error {
		for _, playlistSong := range *playlistSongs {
			if err := tx.Save(&playlistSong).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (p playlistRepository) Delete(id uuid.UUID) error {
	return p.client.Delete(&model.Playlist{}, id).Error
}

func (p playlistRepository) RemoveSongs(playlistSongs *[]model.PlaylistSong) error {
	return p.client.Delete(&playlistSongs).Error
}

func (p playlistRepository) getSongsByPlaylistQuery(userID uuid.UUID) *gorm.DB {
	return p.client.Model(&model.PlaylistSong{}).
		Select("playlist_id, COUNT(*) as songs_count, STRING_AGG(song_id::text, ',') as songs_ids").
		Joins("JOIN playlists ON playlists.id = playlist_songs.playlist_id").
		Where("playlists.user_id = ?", userID).
		Group("playlist_id")
}
