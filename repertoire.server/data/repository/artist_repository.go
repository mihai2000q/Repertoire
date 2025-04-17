package repository

import (
	"gorm.io/gorm"
	"repertoire/server/data/database"
	"repertoire/server/model"
	"slices"
	"strings"

	"github.com/google/uuid"
)

type ArtistRepository interface {
	Get(artist *model.Artist, id uuid.UUID) error
	GetWithAssociations(artist *model.Artist, id uuid.UUID) error
	GetWithBandMembers(artist *model.Artist, id uuid.UUID) error
	GetWithAlbums(artist *model.Artist, id uuid.UUID) error
	GetWithSongs(artist *model.Artist, id uuid.UUID) error
	GetWithAlbumsAndSongs(artist *model.Artist, id uuid.UUID) error
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
	Delete(id uuid.UUID) error
	DeleteAlbums(id uuid.UUID) error
	DeleteSongs(id uuid.UUID) error

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

func (a artistRepository) GetAllByIDsWithSongs(artists *[]model.Artist, ids []uuid.UUID) error {
	return a.client.Model(&model.Artist{}).
		Preload("Songs").
		Find(&artists, ids).
		Error
}

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

	songContains := func(s string) bool {
		return strings.Contains(s, "songs_count") ||
			strings.Contains(s, "rehearsals") ||
			strings.Contains(s, "confidence") ||
			strings.Contains(s, "progress") ||
			strings.Contains(s, "last_time_played ")
	}

	dbSelect := []string{"artists.*"}
	for _, s := range searchBy {
		if strings.Contains(s, "band_members_count") {
			dbSelect = append(dbSelect, a.enhanceArtistsWithBandMembersQuery(tx, userID))
		}
		if strings.Contains(s, "albums_count") {
			dbSelect = append(dbSelect, a.enhanceArtistsWithAlbumsQuery(tx, userID))
		}
		if songContains(s) {
			dbSelect = append(dbSelect, a.enhanceArtistsWithSongsQuery(tx, userID)...)
		}
	}

	for _, o := range orderBy {
		if strings.Contains(o, "band_members_count") &&
			!slices.ContainsFunc(dbSelect, func(s string) bool { return strings.Contains(s, "band_members_count") }) {
			dbSelect = append(dbSelect, a.enhanceArtistsWithBandMembersQuery(tx, userID))
		}
		if strings.Contains(o, "albums_count") &&
			!slices.ContainsFunc(dbSelect, func(s string) bool { return strings.Contains(s, "albums_count") }) {
			dbSelect = append(dbSelect, a.enhanceArtistsWithAlbumsQuery(tx, userID))
		}
		if songContains(o) && !slices.ContainsFunc(dbSelect, func(s string) bool {
			return songContains(s)
		}) {
			dbSelect = append(dbSelect, a.enhanceArtistsWithSongsQuery(tx, userID)...)
		}
	}
	tx.Select(dbSelect)

	tx = database.SearchBy(tx, searchBy)
	tx = database.OrderBy(tx, orderBy)
	tx = database.Paginate(tx, currentPage, pageSize)
	return tx.Find(&artists).Error
}

func (a artistRepository) GetAllByUserCount(count *int64, userID uuid.UUID, searchBy []string) error {
	tx := a.client.Model(&model.Artist{}).
		Where(model.Artist{UserID: userID})

	for _, s := range searchBy {
		if strings.Contains(s, "band_members_count") {
			a.enhanceArtistsWithBandMembersQuery(tx, userID)
		}
		if strings.Contains(s, "albums_count") {
			a.enhanceArtistsWithAlbumsQuery(tx, userID)
		}
		if strings.Contains(s, "songs_count") ||
			strings.Contains(s, "rehearsals") ||
			strings.Contains(s, "confidence") ||
			strings.Contains(s, "progress") {
			a.enhanceArtistsWithSongsQuery(tx, userID)
		}
	}

	tx = database.SearchBy(tx, searchBy)
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

func (a artistRepository) Delete(id uuid.UUID) error {
	return a.client.Delete(&model.Artist{}, id).Error
}

func (a artistRepository) DeleteAlbums(id uuid.UUID) error {
	return a.client.Where("artist_id = ?", id).Delete(&model.Album{}).Error
}

func (a artistRepository) DeleteSongs(id uuid.UUID) error {
	return a.client.Where("artist_id = ?", id).Delete(&model.Song{}).Error
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

func (a artistRepository) enhanceArtistsWithBandMembersQuery(tx *gorm.DB, userID uuid.UUID) string {
	enhancedBandMembers := a.client.Model(&model.BandMember{}).
		Select("artist_id, COUNT(*) as band_members_count").
		Joins("JOIN artists ON artists.id = band_members.artist_id").
		Where("artists.user_id = ?", userID).
		Group("artist_id")

	tx.Joins("LEFT JOIN (?) AS ebm ON ebm.artist_id = artists.id", enhancedBandMembers)
	return "COALESCE(ebm.band_members_count, 0) AS band_members_count"
}

func (a artistRepository) enhanceArtistsWithAlbumsQuery(tx *gorm.DB, userID uuid.UUID) string {
	enhancedAlbums := a.client.Model(&model.Album{}).
		Select("artist_id, COUNT(*) as albums_count").
		Group("artist_id").
		Where(model.Album{UserID: userID})

	tx.Joins("LEFT JOIN (?) AS ea ON ea.artist_id = artists.id", enhancedAlbums)
	return "COALESCE(ea.albums_count, 0) AS albums_count"
}

func (a artistRepository) enhanceArtistsWithSongsQuery(tx *gorm.DB, userID uuid.UUID) []string {
	enhancedSongs := a.client.Model(&model.Song{}).
		Select("artist_id",
			"COUNT(*) as songs_count",
			"ROUND(AVG(rehearsals)) as rehearsals",
			"ROUND(AVG(confidence)) as confidence",
			"CEIL(AVG(progress)) as progress",
			"MAX(last_time_played) as last_time_played",
		).
		Group("artist_id").
		Where(model.Song{UserID: userID})

	tx.Joins("LEFT JOIN (?) AS es ON es.artist_id = artists.id", enhancedSongs)
	return []string{
		"COALESCE(es.songs_count, 0) AS songs_count",
		"COALESCE(es.rehearsals, 0) as rehearsals",
		"COALESCE(es.confidence, 0) as confidence",
		"COALESCE(es.progress, 0) as progress",
		"es.last_time_played as last_time_played",
	}
}
