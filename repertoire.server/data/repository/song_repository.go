package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"repertoire/server/data/database"
	"repertoire/server/model"
)

type SongRepository interface {
	Get(song *model.Song, id uuid.UUID) error
	GetWithPlaylistsAndSongs(song *model.Song, id uuid.UUID) error
	GetWithSections(song *model.Song, id uuid.UUID) error
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
	GetAllByAlbumAndTrackNo(songs *[]model.Song, albumID uuid.UUID, trackNo uint) error
	GetAllByIDs(songs *[]model.Song, ids []uuid.UUID) error
	GetAllByIDsWithSongs(songs *[]model.Song, ids []uuid.UUID) error
	IsBandMemberAssociatedWithSong(songID uuid.UUID, bandMemberID uuid.UUID) (bool, error)
	Create(song *model.Song) error
	Update(song *model.Song) error
	UpdateAll(songs *[]model.Song) error
	UpdateWithAssociations(song *model.Song) error
	UpdateAllWithAssociations(songs *[]model.Song) error
	Delete(id uuid.UUID) error

	GetGuitarTunings(tunings *[]model.GuitarTuning, userID uuid.UUID) error
	GetGuitarTuningsCount(count *int64, userID uuid.UUID) error
	CreateGuitarTuning(tuning *model.GuitarTuning) error
	UpdateAllGuitarTunings(tunings *[]model.GuitarTuning) error
	DeleteGuitarTuning(id uuid.UUID) error

	GetSection(section *model.SongSection, id uuid.UUID) error
	CountSectionsBySong(count *int64, songID uuid.UUID) error
	CreateSection(section *model.SongSection) error
	UpdateSection(section *model.SongSection) error
	DeleteSection(id uuid.UUID) error

	GetSectionTypes(types *[]model.SongSectionType, userID uuid.UUID) error
	CountSectionTypes(count *int64, userID uuid.UUID) error
	CreateSectionType(sectionType *model.SongSectionType) error
	UpdateAllSectionTypes(sectionTypes *[]model.SongSectionType) error
	DeleteSectionType(id uuid.UUID) error

	GetSongSectionHistory(
		history *[]model.SongSectionHistory,
		sectionID uuid.UUID,
		property model.SongSectionProperty,
	) error
	CreateSongSectionHistory(history *model.SongSectionHistory) error
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

func (s songRepository) GetWithPlaylistsAndSongs(song *model.Song, id uuid.UUID) error {
	return s.client.DB.
		Preload("Playlists").
		Preload("Playlists.PlaylistSongs").
		Find(&song, model.Song{ID: id}).Error
}

func (s songRepository) GetWithSections(song *model.Song, id uuid.UUID) error {
	return s.client.DB.
		Preload("Sections", func(db *gorm.DB) *gorm.DB {
			return db.Order("song_sections.order")
		}).
		Find(&song, model.Song{ID: id}).
		Error
}

func (s songRepository) GetWithAssociations(song *model.Song, id uuid.UUID) error {
	return s.client.DB.
		Preload("Sections", func(db *gorm.DB) *gorm.DB {
			return db.Order("song_sections.order")
		}).
		Preload("Sections.SongSectionType").
		Preload("Playlists").
		Joins("GuitarTuning").
		Joins("Artist").
		Joins("Album").
		Preload("Artist.BandMembers").
		Preload("Artist.BandMembers.Roles").
		Find(&song, model.Song{ID: id}).
		Error
}

func (s songRepository) GetAllByIDs(songs *[]model.Song, ids []uuid.UUID) error {
	return s.client.DB.Model(&model.Song{}).Find(&songs, ids).Error
}

func (s songRepository) GetAllByAlbumAndTrackNo(songs *[]model.Song, albumID uuid.UUID, trackNo uint) error {
	return s.client.DB.Model(&model.Song{}).
		Where("album_id = ? AND album_track_no > ?", albumID, trackNo).
		Order("album_track_no").
		Find(&songs).
		Error
}

func (s songRepository) GetAllByIDsWithSongs(songs *[]model.Song, ids []uuid.UUID) error {
	return s.client.DB.Model(&model.Song{}).
		Preload("Album").
		Preload("Album.Songs").
		Find(&songs, ids).
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
		Preload("Sections").
		Preload("Sections.SongSectionType").
		Preload("Playlists").
		Joins("GuitarTuning").
		Joins("Artist").
		Joins("Album").
		Joins("LEFT JOIN playlist_songs ON playlist_songs.song_id = songs.id"). // TODO: Based on the search by add programmatically
		Where(model.Song{UserID: userID})
	tx = database.SearchBy(tx, searchBy)
	tx = database.OrderBy(tx, orderBy)
	tx = database.Paginate(tx, currentPage, pageSize)
	return tx.Find(&songs).Error
}

func (s songRepository) GetAllByUserCount(count *int64, userID uuid.UUID, searchBy []string) error {
	tx := s.client.DB.Model(&model.Song{}).
		Joins("LEFT JOIN playlist_songs ON playlist_songs.song_id = songs.id").
		Where(model.Song{UserID: userID})
	tx = database.SearchBy(tx, searchBy)
	return tx.Count(count).Error
}

func (s songRepository) IsBandMemberAssociatedWithSong(songID uuid.UUID, bandMemberID uuid.UUID) (bool, error) {
	var count int64
	err := s.client.DB.
		Model(&model.Song{}).
		Joins("JOIN artists ON artists.id = songs.artist_id").
		Joins("JOIN band_members ON artists.id = band_members.artist_id").
		Where("songs.id = ?", songID).
		Where("band_members.id = ?", bandMemberID).
		Count(&count).
		Error
	return count != 0, err
}

func (s songRepository) Create(song *model.Song) error {
	return s.client.DB.Create(&song).Error
}

func (s songRepository) Update(song *model.Song) error {
	return s.client.DB.Save(&song).Error
}

func (s songRepository) UpdateAll(songs *[]model.Song) error {
	return s.client.DB.Transaction(func(tx *gorm.DB) error {
		for _, song := range *songs {
			if err := tx.Save(&song).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s songRepository) UpdateWithAssociations(song *model.Song) error {
	return s.client.DB.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Updates(&song).
		Error
}

func (s songRepository) UpdateAllWithAssociations(songs *[]model.Song) error {
	return s.client.DB.Transaction(func(tx *gorm.DB) error {
		for _, song := range *songs {
			err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&song).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s songRepository) Delete(id uuid.UUID) error {
	return s.client.DB.Delete(&model.Song{}, id).Error
}

// Guitar Tunings

func (s songRepository) GetGuitarTunings(tunings *[]model.GuitarTuning, userID uuid.UUID) error {
	return s.client.DB.Model(&model.GuitarTuning{}).
		Where(model.GuitarTuning{UserID: userID}).
		Order("\"order\"").
		Find(&tunings).
		Error
}

func (s songRepository) GetGuitarTuningsCount(count *int64, userID uuid.UUID) error {
	return s.client.DB.Model(&model.GuitarTuning{}).
		Where(model.GuitarTuning{UserID: userID}).
		Count(count).
		Error
}

func (s songRepository) CreateGuitarTuning(tuning *model.GuitarTuning) error {
	return s.client.DB.Create(&tuning).Error
}

func (s songRepository) UpdateAllGuitarTunings(tunings *[]model.GuitarTuning) error {
	return s.client.DB.Transaction(func(tx *gorm.DB) error {
		for _, tuning := range *tunings {
			if err := tx.Save(tuning).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s songRepository) DeleteGuitarTuning(id uuid.UUID) error {
	return s.client.DB.Delete(&model.GuitarTuning{}, id).Error
}

// Sections

func (s songRepository) GetSection(section *model.SongSection, id uuid.UUID) error {
	return s.client.DB.Find(&section, model.SongSection{ID: id}).Error
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

// Section Types

func (s songRepository) GetSectionTypes(types *[]model.SongSectionType, userID uuid.UUID) error {
	return s.client.DB.Model(&model.SongSectionType{}).
		Where(model.SongSectionType{UserID: userID}).
		Order("\"order\"").
		Find(&types).
		Error
}

func (s songRepository) CountSectionTypes(count *int64, userID uuid.UUID) error {
	return s.client.DB.Model(&model.SongSectionType{}).
		Where(model.SongSectionType{UserID: userID}).
		Count(count).
		Error
}

func (s songRepository) CreateSectionType(sectionType *model.SongSectionType) error {
	return s.client.DB.Create(&sectionType).Error
}

func (s songRepository) UpdateAllSectionTypes(sectionTypes *[]model.SongSectionType) error {
	return s.client.DB.Transaction(func(tx *gorm.DB) error {
		for _, sectionType := range *sectionTypes {
			if err := tx.Save(sectionType).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s songRepository) DeleteSectionType(id uuid.UUID) error {
	return s.client.DB.Delete(&model.SongSectionType{}, id).Error
}

// Song Section History

func (s songRepository) GetSongSectionHistory(
	history *[]model.SongSectionHistory,
	sectionID uuid.UUID,
	property model.SongSectionProperty,
) error {
	return s.client.DB.
		Order("created_at").
		Find(&history, model.SongSectionHistory{SongSectionID: sectionID, Property: property}).
		Error
}

func (s songRepository) CreateSongSectionHistory(history *model.SongSectionHistory) error {
	return s.client.DB.Create(&history).Error
}
