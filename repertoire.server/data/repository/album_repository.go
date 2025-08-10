package repository

import (
	"encoding/json"
	"repertoire/server/data/database"
	"repertoire/server/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AlbumRepository interface {
	Get(album *model.Album, id uuid.UUID) error
	GetWithSongs(album *model.Album, id uuid.UUID) error
	GetWithSongsAndArtist(album *model.Album, id uuid.UUID) error
	GetWithAssociations(album *model.Album, id uuid.UUID, songsOrderBy []string) error
	GetFiltersMetadata(metadata *model.AlbumFiltersMetadata, userID uuid.UUID, searchBy []string) error
	GetAllByIDs(albums *[]model.Album, ids []uuid.UUID) error
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
	Delete(ids []uuid.UUID) error
	DeleteWithSongs(ids []uuid.UUID) error
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

func (a albumRepository) GetFiltersMetadata(metadata *model.AlbumFiltersMetadata, userID uuid.UUID, searchBy []string) error {
	tx := a.client.
		Select(
			"JSON_AGG(DISTINCT artist_id) filter (WHERE artist_id IS NOT NULL) as artist_ids",
			"MIN(albums.release_date) AS min_release_date",
			"MAX(albums.release_date) AS max_release_date",
			"MIN(COALESCE(ss.songs_count, 0)) AS min_songs_count",
			"MAX(COALESCE(ss.songs_count, 0)) AS max_songs_count",
			"MIN(COALESCE(CEIL(ss.rehearsals), 0)) as min_rehearsals",
			"MAX(COALESCE(CEIL(ss.rehearsals), 0)) as max_rehearsals",
			"MIN(COALESCE(CEIL(ss.confidence), 0)) as min_confidence",
			"MAX(COALESCE(CEIL(ss.confidence), 0)) as max_confidence",
			"MIN(COALESCE(CEIL(ss.progress), 0)) as min_progress",
			"MAX(COALESCE(CEIL(ss.progress), 0)) as max_progress",
			"MIN(ss.last_time_played) as min_last_time_played",
			"MAX(ss.last_time_played) as max_last_time_played",
		).
		Table("albums").
		Joins("LEFT JOIN (?) AS ss ON ss.album_id = albums.id", a.getSongsSubQuery(userID)).
		Where("user_id = ?", userID)

	searchBy = database.AddCoalesceToCompoundFields(searchBy, compoundAlbumsFields)
	database.SearchBy(tx, searchBy)
	err := tx.Scan(&metadata).Error
	if err != nil {
		return err
	}
	if metadata.ArtistIDsAgg != "" {
		return json.Unmarshal([]byte(metadata.ArtistIDsAgg), &metadata.ArtistIDs)
	}
	return nil
}

func (a albumRepository) GetAllByIDs(albums *[]model.Album, ids []uuid.UUID) error {
	return a.client.Model(&model.Album{}).Find(&albums, ids).Error
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

var compoundAlbumsFields = []string{"songs_count", "rehearsals", "confidence", "progress"}

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

	a.addSongsSubQuery(tx, userID)

	searchBy = database.AddCoalesceToCompoundFields(searchBy, compoundAlbumsFields)
	orderBy = database.AddCoalesceToCompoundFields(orderBy, compoundAlbumsFields)

	database.SearchBy(tx, searchBy)
	database.OrderBy(tx, orderBy)
	database.Paginate(tx, currentPage, pageSize)
	return tx.Find(&albums).Error
}

func (a albumRepository) GetAllByUserCount(count *int64, userID uuid.UUID, searchBy []string) error {
	tx := a.client.Model(&model.Album{}).
		Joins("Artist").
		Where(model.Album{UserID: userID})

	a.addSongsSubQuery(tx, userID)

	searchBy = database.AddCoalesceToCompoundFields(searchBy, compoundAlbumsFields)

	database.SearchBy(tx, searchBy)
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

func (a albumRepository) Delete(ids []uuid.UUID) error {
	return a.client.Delete(&model.Album{}, ids).Error
}

func (a albumRepository) DeleteWithSongs(ids []uuid.UUID) error {
	return a.client.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("album_id IN (?)", ids).Delete(&model.Song{}).Error
		if err != nil {
			return err
		}

		err = tx.Delete(&model.Album{}, ids).Error
		if err != nil {
			return err
		}
		return nil
	})
}

func (a albumRepository) RemoveSongs(album *model.Album, songs *[]model.Song) error {
	return a.client.Model(&album).Association("Songs").Delete(&songs)
}

func (a albumRepository) addSongsSubQuery(tx *gorm.DB, userID uuid.UUID) {
	tx.Joins("LEFT JOIN (?) AS ss ON ss.album_id = albums.id", a.getSongsSubQuery(userID)).
		Select("albums.*",
			"COALESCE(ss.songs_count, 0) AS songs_count",
			"COALESCE(ss.rehearsals, 0) as rehearsals",
			"COALESCE(ss.confidence, 0) as confidence",
			"COALESCE(ss.progress, 0) as progress",
			"ss.last_time_played as last_time_played",
		)
}

func (a albumRepository) getSongsSubQuery(userID uuid.UUID) *gorm.DB {
	return a.client.Model(&model.Song{}).
		Select("album_id",
			"COUNT(*) as songs_count",
			"AVG(rehearsals) as rehearsals",
			"AVG(confidence) as confidence",
			"AVG(progress) as progress",
			"MAX(last_time_played) as last_time_played",
		).
		Group("album_id").
		Where(model.Song{UserID: userID})
}
