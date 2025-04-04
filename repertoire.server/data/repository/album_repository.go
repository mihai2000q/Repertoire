package repository

import (
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
	GetAllByIDsWithSongs(albums *[]model.Album, ids []uuid.UUID) error
	GetAllByIDsWithSongsAndArtist(albums *[]model.Album, ids []uuid.UUID) error
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
	albums *[]model.Album,
	userID uuid.UUID,
	currentPage *int,
	pageSize *int,
	orderBy []string,
	searchBy []string,
) error {
	tx := a.client.Model(&model.Album{}).
		Joins("Artist").
		Where(model.Album{UserID: userID})
	tx = database.SearchBy(tx, searchBy)
	tx = database.OrderBy(tx, orderBy)
	tx = database.Paginate(tx, currentPage, pageSize)
	return tx.Find(&albums).Error
}

func (a albumRepository) GetAllByUserCount(count *int64, userID uuid.UUID, searchBy []string) error {
	tx := a.client.Model(&model.Album{}).
		Where(model.Album{UserID: userID})
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
