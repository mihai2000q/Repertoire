package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"repertoire/server/data/database"
	"repertoire/server/model"
)

type ArtistRepository interface {
	Get(artist *model.Artist, id uuid.UUID) error
	GetWithAssociations(artist *model.Artist, id uuid.UUID) error
	GetWithBandMembers(artist *model.Artist, id uuid.UUID) error
	GetWithAlbums(artist *model.Artist, id uuid.UUID) error
	GetWithSongs(artist *model.Artist, id uuid.UUID) error
	GetWithAlbumsAndSongs(artist *model.Artist, id uuid.UUID) error
	GetFiltersMetadata(metadata *model.ArtistFiltersMetadata, userID uuid.UUID, searchBy []string) error
	GetAllByIDsWithSongs(artists *[]model.Artist, ids []uuid.UUID) error
	GetAllByUser(
		artists *[]model.EnhancedArtist,
		userID uuid.UUID,
		currentPage *int,
		pageSize *int,
		orderBy []string,
		searchBy []string,
	) error
	GetAllByUserCount(count *int64, userID uuid.UUID, searchBy []string) error
	Create(artist *model.Artist) error
	Update(artist *model.Artist) error
	UpdateWithAssociations(artist *model.Artist) error
	Delete(ids []uuid.UUID) error
	DeleteAlbums(ids []uuid.UUID) error
	DeleteSongs(ids []uuid.UUID) error

	GetBandMember(bandMember *model.BandMember, id uuid.UUID) error
	GetBandMemberWithArtist(bandMember *model.BandMember, id uuid.UUID) error
	CreateBandMember(bandMember *model.BandMember) error
	UpdateBandMember(bandMember *model.BandMember) error
	ReplaceRolesFromBandMember(roles []model.BandMemberRole, bandMember *model.BandMember) error
	DeleteBandMember(id uuid.UUID) error

	GetBandMemberRoles(roles *[]model.BandMemberRole, userID uuid.UUID) error
	GetBandMemberRolesByIDs(roles *[]model.BandMemberRole, ids []uuid.UUID) error
}

type artistRepository struct {
	client database.Client
}

func NewArtistRepository(client database.Client) ArtistRepository {
	return artistRepository{
		client: client,
	}
}

func (a artistRepository) Get(artist *model.Artist, id uuid.UUID) error {
	return a.client.Find(&artist, model.Artist{ID: id}).Error
}

func (a artistRepository) GetWithAssociations(artist *model.Artist, id uuid.UUID) error {
	return a.client.
		Preload("BandMembers", func(db *gorm.DB) *gorm.DB {
			return db.Order("band_members.order")
		}).
		Preload("BandMembers.Roles").
		Find(&artist, model.Artist{ID: id}).
		Error
}

func (a artistRepository) GetWithBandMembers(artist *model.Artist, id uuid.UUID) error {
	return a.client.
		Preload("BandMembers", func(db *gorm.DB) *gorm.DB {
			return db.Order("band_members.order")
		}).
		Find(&artist, model.Artist{ID: id}).
		Error
}

func (a artistRepository) GetWithAlbums(artist *model.Artist, id uuid.UUID) error {
	return a.client.
		Preload("Albums").
		Find(&artist, model.Artist{ID: id}).
		Error
}

func (a artistRepository) GetWithSongs(artist *model.Artist, id uuid.UUID) error {
	return a.client.
		Preload("Songs").
		Preload("Songs.Album").
		Find(&artist, model.Artist{ID: id}).
		Error
}

func (a artistRepository) GetWithAlbumsAndSongs(artist *model.Artist, id uuid.UUID) error {
	return a.client.
		Preload("Albums").
		Preload("Songs").
		Preload("Songs.Album").
		Find(&artist, model.Artist{ID: id}).
		Error
}

func (a artistRepository) GetFiltersMetadata(metadata *model.ArtistFiltersMetadata, userID uuid.UUID, searchBy []string) error {
	tx := a.client.
		Select(
			"MIN(COALESCE(bms.band_members_count, 0)) AS min_band_members_count",
			"MAX(COALESCE(bms.band_members_count, 0)) AS max_band_members_count",
			"MIN(COALESCE(aas.albums_count, 0)) AS min_albums_count",
			"MAX(COALESCE(aas.albums_count, 0)) AS max_albums_count",
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
		Table("artists").
		Joins("LEFT JOIN (?) AS bms ON bms.artist_id = artists.id", a.getBandMembersSubQuery(userID)).
		Joins("LEFT JOIN (?) AS aas ON aas.artist_id = artists.id", a.getAlbumsSubQuery(userID)).
		Joins("LEFT JOIN (?) AS ss ON ss.artist_id = artists.id", a.getSongsSubQuery(userID)).
		Where("user_id = ?", userID)

	searchBy = database.AddCoalesceToCompoundFields(searchBy, compoundArtistsFields)
	database.SearchBy(tx, searchBy)
	return tx.Scan(&metadata).Error
}

func (a artistRepository) GetAllByIDsWithSongs(artists *[]model.Artist, ids []uuid.UUID) error {
	return a.client.Model(&model.Artist{}).
		Preload("Songs").
		Find(&artists, ids).
		Error
}

var compoundArtistsFields = []string{"band_members_count", "albums_count", "songs_count",
	"rehearsals", "confidence", "progress"}

func (a artistRepository) GetAllByUser(
	artists *[]model.EnhancedArtist,
	userID uuid.UUID,
	currentPage *int,
	pageSize *int,
	orderBy []string,
	searchBy []string,
) error {
	tx := a.client.Model(&model.Artist{}).
		Where(model.Artist{UserID: userID})

	dbSelect := []string{"artists.*"}
	dbSelect = append(dbSelect, a.addBandMembersSubQuery(tx, userID))
	dbSelect = append(dbSelect, a.addAlbumsSubQuery(tx, userID))
	dbSelect = append(dbSelect, a.addSongsSubQuery(tx, userID)...)
	tx.Select(dbSelect)

	searchBy = database.AddCoalesceToCompoundFields(searchBy, compoundArtistsFields)
	orderBy = database.AddCoalesceToCompoundFields(orderBy, compoundArtistsFields)

	database.SearchBy(tx, searchBy)
	database.OrderBy(tx, orderBy)
	database.Paginate(tx, currentPage, pageSize)
	return tx.Find(&artists).Error
}

func (a artistRepository) GetAllByUserCount(count *int64, userID uuid.UUID, searchBy []string) error {
	tx := a.client.Model(&model.Artist{}).
		Where(model.Artist{UserID: userID})

	dbSelect := []string{"artists.*"}
	dbSelect = append(dbSelect, a.addBandMembersSubQuery(tx, userID))
	dbSelect = append(dbSelect, a.addAlbumsSubQuery(tx, userID))
	dbSelect = append(dbSelect, a.addSongsSubQuery(tx, userID)...)
	tx.Select(dbSelect)

	searchBy = database.AddCoalesceToCompoundFields(searchBy, compoundArtistsFields)

	database.SearchBy(tx, searchBy)
	return tx.Count(count).Error
}

func (a artistRepository) Create(artist *model.Artist) error {
	return a.client.Create(&artist).Error
}

func (a artistRepository) Update(artist *model.Artist) error {
	return a.client.Save(&artist).Error
}

func (a artistRepository) UpdateWithAssociations(artist *model.Artist) error {
	return a.client.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Updates(&artist).
		Error
}

func (a artistRepository) Delete(ids []uuid.UUID) error {
	return a.client.Delete(&model.Artist{}, ids).Error
}

func (a artistRepository) DeleteAlbums(ids []uuid.UUID) error {
	return a.client.Where("artist_id IN (?)", ids).Delete(&model.Album{}).Error
}

func (a artistRepository) DeleteSongs(ids []uuid.UUID) error {
	return a.client.Where("artist_id IN (?)", ids).Delete(&model.Song{}).Error
}

// Band Member

func (a artistRepository) GetBandMember(bandMember *model.BandMember, id uuid.UUID) error {
	return a.client.Find(&bandMember, id).Error
}

func (a artistRepository) GetBandMemberWithArtist(bandMember *model.BandMember, id uuid.UUID) error {
	return a.client.Preload("Artist").Find(&bandMember, id).Error
}

func (a artistRepository) CreateBandMember(bandMember *model.BandMember) error {
	return a.client.Create(&bandMember).Error
}

func (a artistRepository) UpdateBandMember(bandMember *model.BandMember) error {
	return a.client.Save(&bandMember).Error
}

func (a artistRepository) ReplaceRolesFromBandMember(roles []model.BandMemberRole, bandMember *model.BandMember) error {
	return a.client.Model(&bandMember).Association("Roles").Replace(roles)
}

func (a artistRepository) DeleteBandMember(id uuid.UUID) error {
	return a.client.Delete(&model.BandMember{}, id).Error
}

// Band Member - Roles

func (a artistRepository) GetBandMemberRoles(bandMemberRoles *[]model.BandMemberRole, userID uuid.UUID) error {
	return a.client.Find(&bandMemberRoles, model.BandMemberRole{UserID: userID}).Error
}

func (a artistRepository) GetBandMemberRolesByIDs(bandMemberRoles *[]model.BandMemberRole, ids []uuid.UUID) error {
	return a.client.Find(&bandMemberRoles, ids).Error
}

func (a artistRepository) addBandMembersSubQuery(tx *gorm.DB, userID uuid.UUID) string {
	tx.Joins("LEFT JOIN (?) AS bms ON bms.artist_id = artists.id", a.getBandMembersSubQuery(userID))
	return "COALESCE(bms.band_members_count, 0) AS band_members_count"
}

func (a artistRepository) addAlbumsSubQuery(tx *gorm.DB, userID uuid.UUID) string {
	tx.Joins("LEFT JOIN (?) AS asq ON asq.artist_id = artists.id", a.getAlbumsSubQuery(userID))
	return "COALESCE(asq.albums_count, 0) AS albums_count"
}

func (a artistRepository) addSongsSubQuery(tx *gorm.DB, userID uuid.UUID) []string {
	tx.Joins("LEFT JOIN (?) AS ss ON ss.artist_id = artists.id", a.getSongsSubQuery(userID))
	return []string{
		"COALESCE(ss.songs_count, 0) AS songs_count",
		"COALESCE(ss.rehearsals, 0) as rehearsals",
		"COALESCE(ss.confidence, 0) as confidence",
		"COALESCE(ss.progress, 0) as progress",
		"ss.last_time_played as last_time_played",
	}
}

func (a artistRepository) getBandMembersSubQuery(userID uuid.UUID) *gorm.DB {
	return a.client.Model(&model.BandMember{}).
		Select("artist_id, COUNT(*) as band_members_count").
		Joins("JOIN artists ON artists.id = band_members.artist_id").
		Where("artists.user_id = ?", userID).
		Group("artist_id")
}

func (a artistRepository) getAlbumsSubQuery(userID uuid.UUID) *gorm.DB {
	return a.client.Model(&model.Album{}).
		Select("artist_id, COUNT(*) AS albums_count").
		Group("artist_id").
		Where(model.Album{UserID: userID})
}

func (a artistRepository) getSongsSubQuery(userID uuid.UUID) *gorm.DB {
	return a.client.Model(&model.Song{}).
		Select("artist_id",
			"COUNT(*) as songs_count",
			"AVG(rehearsals) as rehearsals",
			"AVG(confidence) as confidence",
			"AVG(progress) as progress",
			"MAX(last_time_played) as last_time_played",
		).
		Group("artist_id").
		Where(model.Song{UserID: userID})
}
