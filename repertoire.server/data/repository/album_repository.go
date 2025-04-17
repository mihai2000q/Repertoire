package repository

import (
	"repertoire/server/data/database"
	"repertoire/server/model"
	"slices"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AlbumRepository interface {
	Get(album *model.Album, id uuid.UUID) error
	GetWithSongs(album *model.Album, id uuid.UUID) error
	GetWithSongsAndArtist(album *model.Album, id uuid.UUID) error
	GetWithAssociations(album *model.Album, id uuid.UUID, songsOrderBy []string) error
	GetAllByIDsWithSongs(albums *[]model.Album, ids []uuid.UUID) error
	GetAllByIDsWithSongsAndArtist(albums *[]model.Album, ids []uuid.UUID) error
	GetAllByUser(
		albums *[]model.EnhancedAlbum,
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
	UpdateAllWithSongs(albums *[]model.Album) error
	Delete(id uuid.UUID) error
	DeleteWithSongs(id uuid.UUID) error
	RemoveSongs(album *model.Album, song *[]model.Song) error
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
	return a.client.Find(&album, model.Album{ID: id}).Error
}

func (a albumRepository) GetWithSongs(album *model.Album, id uuid.UUID) error {
	return a.client.
		Preload("Songs", func(db *gorm.DB) *gorm.DB {
			return db.Order("songs.album_track_no")
		}).
		Find(&album, model.Album{ID: id}).
		Error
}

func (a albumRepository) GetWithSongsAndArtist(album *model.Album, id uuid.UUID) error {
	return a.client.
		Joins("Artist").
		Preload("Songs").
		Preload("Songs.Artist").
		Find(&album, model.Album{ID: id}).
		Error
}

func (a albumRepository) GetWithAssociations(album *model.Album, id uuid.UUID, songsOrderBy []string) error {
	return a.client.
		Preload("Songs", func(db *gorm.DB) *gorm.DB {
			return database.OrderBy(db, songsOrderBy)
		}).
		Joins("Artist").
		Find(&album, model.Album{ID: id}).
		Error
}

func (a albumRepository) GetAllByIDsWithSongs(albums *[]model.Album, ids []uuid.UUID) error {
	return a.client.Model(&model.Album{}).
		Preload("Songs", func(db *gorm.DB) *gorm.DB {
			return db.Order("songs.album_track_no")
		}).
		Find(&albums, ids).
		Error
}

func (a albumRepository) GetAllByIDsWithSongsAndArtist(albums *[]model.Album, ids []uuid.UUID) error {
	return a.client.Model(&model.Album{}).
		Joins("Artist").
		Preload("Songs").
		Preload("Songs.Artist").
		Find(&albums, ids).
		Error
}

func (a albumRepository) GetAllByUser(
	albums *[]model.EnhancedAlbum,
	userID uuid.UUID,
	currentPage *int,
	pageSize *int,
	orderBy []string,
	searchBy []string,
) error {
	tx := a.client.Model(&model.Album{}).
		Joins("Artist").
		Where(model.Album{UserID: userID}).
		Select("albums.*")

	containsFunc := func(s string) bool {
		return strings.Contains(s, "songs_count") ||
			strings.Contains(s, "rehearsals") ||
			strings.Contains(s, "confidence") ||
			strings.Contains(s, "progress") ||
			strings.Contains(s, "last_time_played")
	}

	if slices.ContainsFunc(searchBy, containsFunc) || slices.ContainsFunc(orderBy, containsFunc) {
		a.enhanceAlbumsWithSongsQuery(tx, userID)
	}

	tx = database.SearchBy(tx, searchBy)
	tx = database.OrderBy(tx, orderBy)
	tx = database.Paginate(tx, currentPage, pageSize)
	return tx.Find(&albums).Error
}

func (a albumRepository) GetAllByUserCount(count *int64, userID uuid.UUID, searchBy []string) error {
	tx := a.client.Model(&model.Album{}).
		Joins("Artist").
		Where(model.Album{UserID: userID})

	containsFunc := func(s string) bool {
		return strings.Contains(s, "songs_count") ||
			strings.Contains(s, "rehearsals") ||
			strings.Contains(s, "confidence") ||
			strings.Contains(s, "progress") ||
			strings.Contains(s, "last_time_played")
	}

	if slices.ContainsFunc(searchBy, containsFunc) {
		a.enhanceAlbumsWithSongsQuery(tx, userID)
	}

	tx = database.SearchBy(tx, searchBy)
	return tx.Count(count).Error
}

func (a albumRepository) Create(album *model.Album) error {
	return a.client.Create(&album).Error
}

func (a albumRepository) Update(album *model.Album) error {
	return a.client.Save(&album).Error
}

func (a albumRepository) UpdateWithAssociations(album *model.Album) error {
	return a.client.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Updates(&album).
		Error
}

func (a albumRepository) UpdateAllWithSongs(albums *[]model.Album) error {
	return a.client.Transaction(func(tx *gorm.DB) error {
		for _, album := range *albums {
			if err := tx.Save(&album).Error; err != nil {
				return err
			}
			for _, song := range album.Songs {
				if err := tx.Save(&song).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (a albumRepository) Delete(id uuid.UUID) error {
	return a.client.Delete(&model.Album{}, id).Error
}

func (a albumRepository) DeleteWithSongs(id uuid.UUID) error {
	return a.client.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("album_id = ?", id).Delete(&model.Song{}).Error
		if err != nil {
			return err
		}

		err = tx.Delete(&model.Album{}, id).Error
		if err != nil {
			return err
		}
		return nil
	})
}

func (a albumRepository) RemoveSongs(album *model.Album, songs *[]model.Song) error {
	return a.client.Model(&album).Association("Songs").Delete(&songs)
}

func (a albumRepository) enhanceAlbumsWithSongsQuery(tx *gorm.DB, userID uuid.UUID) {
	enhancedSongs := a.client.Model(&model.Song{}).
		Select("album_id",
			"COUNT(*) as songs_count",
			"ROUND(AVG(rehearsals)) as rehearsals",
			"ROUND(AVG(confidence)) as confidence",
			"CEIL(AVG(progress)) as progress",
			"MAX(last_time_played) as last_time_played",
		).
		Group("album_id").
		Where(model.Song{UserID: userID})

	tx.Joins("LEFT JOIN (?) AS es ON es.album_id = albums.id", enhancedSongs)
	tx.Select("albums.*",
		"COALESCE(es.songs_count, 0) AS songs_count",
		"COALESCE(es.rehearsals, 0) as rehearsals",
		"COALESCE(es.confidence, 0) as confidence",
		"COALESCE(es.progress, 0) as progress",
		"es.last_time_played as last_time_played",
	)
}
