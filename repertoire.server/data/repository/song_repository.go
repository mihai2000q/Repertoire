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
		songs *[]model.EnhancedSong,
		userID uuid.UUID,
		currentPage *int,
		pageSize *int,
		orderBy []string,
		searchBy []string,
	) error
	GetAllByUserCount(count *int64, userID uuid.UUID, searchBy []string) error
	GetAllByAlbum(songs *[]model.Song, albumID uuid.UUID) error
	GetAllByAlbumAndTrackNo(songs *[]model.Song, albumID uuid.UUID, trackNo uint) error
	GetAllByIDs(songs *[]model.Song, ids []uuid.UUID) error
	GetAllByIDsWithSongs(songs *[]model.Song, ids []uuid.UUID) error
	GetAllByIDsWithArtistAndAlbum(songs *[]model.Song, ids []uuid.UUID) error
	CountByAlbum(count *int64, albumID uuid.UUID) error
	IsBandMemberAssociatedWithSong(songID uuid.UUID, bandMemberID uuid.UUID) (bool, error)
	Create(song *model.Song) error
	Update(song *model.Song) error
	UpdateAll(songs *[]model.Song) error
	UpdateWithAssociations(song *model.Song) error
	UpdateAllWithAssociations(songs *[]model.Song) error
	Delete(id uuid.UUID) error

	GetSettings(settings *model.SongSettings, settingsID uuid.UUID) error
	UpdateSettings(settings *model.SongSettings) error

	GetGuitarTunings(tunings *[]model.GuitarTuning, userID uuid.UUID) error
	GetInstruments(instruments *[]model.Instrument, userID uuid.UUID) error
	GetSectionTypes(types *[]model.SongSectionType, userID uuid.UUID) error

	GetSection(section *model.SongSection, id uuid.UUID) error
	CountSectionsBySong(count *int64, songID uuid.UUID) error
	CreateSection(section *model.SongSection) error
	UpdateSection(section *model.SongSection) error
	DeleteSection(id uuid.UUID) error

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
	return s.client.Find(&song, model.Song{ID: id}).Error
}

func (s songRepository) GetWithPlaylistsAndSongs(song *model.Song, id uuid.UUID) error {
	return s.client.
		Preload("Playlists").
		Preload("Playlists.PlaylistSongs").
		Find(&song, model.Song{ID: id}).
		Error
}

func (s songRepository) GetWithSections(song *model.Song, id uuid.UUID) error {
	return s.client.
		Preload("Sections", func(db *gorm.DB) *gorm.DB {
			return db.Order("song_sections.order")
		}).
		Find(&song, model.Song{ID: id}).
		Error
}

func (s songRepository) GetWithAssociations(song *model.Song, id uuid.UUID) error {
	return s.client.
		Joins("Settings").
		Joins("Settings.DefaultBandMember").
		Joins("Settings.DefaultInstrument").
		Joins("GuitarTuning").
		Joins("Artist").
		Joins("Album").
		Preload("Artist.BandMembers").
		Preload("Artist.BandMembers.Roles").
		Preload("Sections", func(db *gorm.DB) *gorm.DB {
			return db.Order("song_sections.order")
		}).
		Preload("Sections.SongSectionType").
		Preload("Sections.Instrument").
		Preload("Sections.BandMember").
		Preload("Sections.BandMember.Roles").
		Find(&song, model.Song{ID: id}).
		Error
}

func (s songRepository) GetAllByUser(
	songs *[]model.EnhancedSong,
	userID uuid.UUID,
	currentPage *int,
	pageSize *int,
	orderBy []string,
	searchBy []string,
) error {
	tx := s.client.Model(&model.Song{}).
		Preload("Sections").
		Preload("Sections.SongSectionType").
		Preload("Sections.Instrument").
		Joins("GuitarTuning").
		Joins("Artist").
		Joins("Album").
		Joins("LEFT JOIN playlist_songs ON playlist_songs.song_id = songs.id"). // TODO: Based on the search by add programmatically
		Where(model.Song{UserID: userID})

	s.enhanceSongsWithSongSectionsQuery(tx, userID)

	tx = database.SearchBy(tx, searchBy)
	tx = database.OrderBy(tx, orderBy)
	tx = database.Paginate(tx, currentPage, pageSize)
	return tx.Find(&songs).Error
}

func (s songRepository) GetAllByUserCount(count *int64, userID uuid.UUID, searchBy []string) error {
	tx := s.client.Model(&model.Song{}).
		Distinct("songs.id").
		Joins("LEFT JOIN playlist_songs ON playlist_songs.song_id = songs.id"). // TODO REMOVE
		Where(model.Song{UserID: userID})

	s.enhanceSongsWithSongSectionsQuery(tx, userID)

	tx = database.SearchBy(tx, searchBy)
	return tx.Count(count).Error
}

func (s songRepository) GetAllByIDs(songs *[]model.Song, ids []uuid.UUID) error {
	return s.client.Model(&model.Song{}).Find(&songs, ids).Error
}

func (s songRepository) GetAllByAlbum(songs *[]model.Song, albumID uuid.UUID) error {
	return s.client.Model(&model.Song{}).
		Find(&songs, model.Song{AlbumID: &albumID}).
		Error
}

func (s songRepository) GetAllByAlbumAndTrackNo(songs *[]model.Song, albumID uuid.UUID, trackNo uint) error {
	return s.client.Model(&model.Song{}).
		Where("album_id = ? AND album_track_no > ?", albumID, trackNo).
		Order("album_track_no").
		Find(&songs).
		Error
}

func (s songRepository) GetAllByIDsWithSongs(songs *[]model.Song, ids []uuid.UUID) error {
	return s.client.Model(&model.Song{}).
		Preload("Album").
		Preload("Album.Songs").
		Find(&songs, ids).
		Error
}

func (s songRepository) GetAllByIDsWithArtistAndAlbum(songs *[]model.Song, ids []uuid.UUID) error {
	return s.client.
		Joins("Artist").
		Joins("Album").
		Find(&songs, ids).
		Error
}

func (s songRepository) CountByAlbum(count *int64, albumID uuid.UUID) error {
	return s.client.Model(&model.Song{}).
		Where("album_id = ?", albumID).
		Count(count).
		Error
}

// TODO: Isn't this authorization basically? (therefore, it might need to be deleted)

func (s songRepository) IsBandMemberAssociatedWithSong(songID uuid.UUID, bandMemberID uuid.UUID) (bool, error) {
	var count int64
	err := s.client.
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
	return s.client.Create(&song).Error
}

func (s songRepository) Update(song *model.Song) error {
	return s.client.Save(&song).Error
}

func (s songRepository) UpdateAll(songs *[]model.Song) error {
	return s.client.Transaction(func(tx *gorm.DB) error {
		for _, song := range *songs {
			if err := tx.Save(&song).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s songRepository) UpdateWithAssociations(song *model.Song) error {
	return s.client.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Updates(&song).
		Error
}

func (s songRepository) UpdateAllWithAssociations(songs *[]model.Song) error {
	return s.client.Transaction(func(tx *gorm.DB) error {
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
	return s.client.Delete(&model.Song{}, id).Error
}

func (s songRepository) enhanceSongsWithSongSectionsQuery(tx *gorm.DB, userID uuid.UUID) {
	enhancedSections := s.client.Model(&model.SongSection{}).
		Select("song_id",
			"COUNT(*) as sections_count",
			"COUNT(*) filter (where song_section_types.name = 'Solo') as solos",
			"COUNT(*) filter (where song_section_types.name = 'Riff') as riffs",
		).
		Joins("LEFT JOIN song_section_types ON song_section_types.id = song_sections.song_section_type_id").
		Joins("JOIN songs ON songs.id = song_sections.song_id").
		Where("songs.user_id = ?", userID).
		Group("song_id")

	tx.Joins("LEFT JOIN (?) AS es ON es.song_id = songs.id", enhancedSections).
		Select(
			"songs.*",
			"COALESCE(es.sections_count, 0) as sections_count",
			"COALESCE(es.solos, 0) as solos",
			"COALESCE(es.riffs, 0) as riffs",
		)
}

// Settings

func (s songRepository) GetSettings(settings *model.SongSettings, settingsID uuid.UUID) error {
	return s.client.Find(&settings, settingsID).Error
}

func (s songRepository) UpdateSettings(settings *model.SongSettings) error {
	return s.client.Save(&settings).Error
}

// Guitar Tunings

func (s songRepository) GetGuitarTunings(tunings *[]model.GuitarTuning, userID uuid.UUID) error {
	return s.client.Model(&model.GuitarTuning{}).
		Where(model.GuitarTuning{UserID: userID}).
		Order("\"order\"").
		Find(&tunings).
		Error
}

// Instruments

func (s songRepository) GetInstruments(instruments *[]model.Instrument, userID uuid.UUID) error {
	return s.client.Model(&model.Instrument{}).
		Where(model.Instrument{UserID: userID}).
		Order("\"order\"").
		Find(&instruments).
		Error
}

// Section Types

func (s songRepository) GetSectionTypes(types *[]model.SongSectionType, userID uuid.UUID) error {
	return s.client.Model(&model.SongSectionType{}).
		Where(model.SongSectionType{UserID: userID}).
		Order("\"order\"").
		Find(&types).
		Error
}

// Sections

func (s songRepository) GetSection(section *model.SongSection, id uuid.UUID) error {
	return s.client.Find(&section, model.SongSection{ID: id}).Error
}

func (s songRepository) CountSectionsBySong(count *int64, songID uuid.UUID) error {
	return s.client.Model(&model.SongSection{}).
		Where(model.SongSection{SongID: songID}).
		Count(count).
		Error
}

func (s songRepository) CreateSection(section *model.SongSection) error {
	return s.client.Create(&section).Error
}

func (s songRepository) UpdateSection(section *model.SongSection) error {
	return s.client.Save(&section).Error
}

func (s songRepository) DeleteSection(id uuid.UUID) error {
	return s.client.Delete(&model.SongSection{}, id).Error
}

// Song Section History

func (s songRepository) GetSongSectionHistory(
	history *[]model.SongSectionHistory,
	sectionID uuid.UUID,
	property model.SongSectionProperty,
) error {
	return s.client.
		Order("created_at").
		Find(&history, model.SongSectionHistory{SongSectionID: sectionID, Property: property}).
		Error
}

func (s songRepository) CreateSongSectionHistory(history *model.SongSectionHistory) error {
	return s.client.Create(&history).Error
}
