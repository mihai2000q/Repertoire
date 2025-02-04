package repository

import (
	"gorm.io/gorm"
	"repertoire/server/data/database"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type ArtistRepository interface {
	Get(artist *model.Artist, id uuid.UUID) error
	GetWithAssociations(artist *model.Artist, id uuid.UUID) error
	GetWithBandMembers(artist *model.Artist, id uuid.UUID) error
	GetAllByIDsWithSongs(artists *[]model.Artist, ids []uuid.UUID) error
	GetAllByUser(
		artists *[]model.Artist,
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
	CountBandMemberRoles(count *int64, userID uuid.UUID) error
	CreateBandMemberRole(bandMember *model.BandMemberRole) error
	UpdateBandMemberRole(bandMember *model.BandMemberRole) error
	UpdateAllBandMemberRoles(bandMemberRoles *[]model.BandMemberRole) error
	DeleteBandMemberRole(id uuid.UUID) error
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
	return a.client.DB.Find(&artist, model.Artist{ID: id}).Error
}

func (a artistRepository) GetWithAssociations(artist *model.Artist, id uuid.UUID) error {
	return a.client.DB.
		Preload("BandMembers").
		Preload("BandMembers.Roles").
		Find(&artist, model.Artist{ID: id}).
		Error
}

func (a artistRepository) GetWithBandMembers(artist *model.Artist, id uuid.UUID) error {
	return a.client.DB.
		Preload("BandMembers").
		Find(&artist, model.Artist{ID: id}).
		Error
}

func (a artistRepository) GetAllByIDsWithSongs(artists *[]model.Artist, ids []uuid.UUID) error {
	return a.client.DB.Model(&model.Artist{}).
		Preload("Songs").
		Find(&artists, ids).
		Error
}

func (a artistRepository) GetAllByUser(
	artists *[]model.Artist,
	userID uuid.UUID,
	currentPage *int,
	pageSize *int,
	orderBy []string,
	searchBy []string,
) error {
	tx := a.client.DB.Model(&model.Artist{}).
		Where(model.Artist{UserID: userID})
	tx = database.SearchBy(tx, searchBy)
	tx = database.OrderBy(tx, orderBy)
	tx = database.Paginate(tx, currentPage, pageSize)
	return tx.Find(&artists).Error
}

func (a artistRepository) GetAllByUserCount(count *int64, userID uuid.UUID, searchBy []string) error {
	tx := a.client.DB.Model(&model.Artist{}).
		Where(model.Artist{UserID: userID})
	tx = database.SearchBy(tx, searchBy)
	return tx.Count(count).Error
}

func (a artistRepository) Create(artist *model.Artist) error {
	return a.client.DB.Create(&artist).Error
}

func (a artistRepository) Update(artist *model.Artist) error {
	return a.client.DB.Save(&artist).Error
}

func (a artistRepository) UpdateWithAssociations(artist *model.Artist) error {
	return a.client.DB.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Updates(&artist).
		Error
}

func (a artistRepository) Delete(id uuid.UUID) error {
	return a.client.DB.Delete(&model.Artist{}, id).Error
}

func (a artistRepository) DeleteAlbums(id uuid.UUID) error {
	return a.client.DB.Where("artist_id = ?", id).Delete(&model.Album{}).Error
}

func (a artistRepository) DeleteSongs(id uuid.UUID) error {
	return a.client.DB.Where("artist_id = ?", id).Delete(&model.Song{}).Error
}

// Band Member

func (a artistRepository) GetBandMember(bandMember *model.BandMember, id uuid.UUID) error {
	return a.client.DB.Find(&bandMember, id).Error
}

func (a artistRepository) GetBandMemberWithArtist(bandMember *model.BandMember, id uuid.UUID) error {
	return a.client.DB.Preload("Artist").Find(&bandMember, id).Error
}

func (a artistRepository) CreateBandMember(bandMember *model.BandMember) error {
	return a.client.DB.Create(&bandMember).Error
}

func (a artistRepository) UpdateBandMember(bandMember *model.BandMember) error {
	return a.client.DB.Save(&bandMember).Error
}

func (a artistRepository) ReplaceRolesFromBandMember(roles []model.BandMemberRole, bandMember *model.BandMember) error {
	return a.client.DB.Model(&bandMember).Association("Roles").Replace(roles)
}

func (a artistRepository) DeleteBandMember(id uuid.UUID) error {
	return a.client.DB.Delete(&model.BandMember{}, id).Error
}

// Band Member - Roles

func (a artistRepository) GetBandMemberRoles(bandMemberRoles *[]model.BandMemberRole, userID uuid.UUID) error {
	return a.client.DB.Find(&bandMemberRoles, model.BandMemberRole{UserID: userID}).Error
}

func (a artistRepository) GetBandMemberRolesByIDs(bandMemberRoles *[]model.BandMemberRole, ids []uuid.UUID) error {
	return a.client.DB.Find(&bandMemberRoles, ids).Error
}

func (a artistRepository) CountBandMemberRoles(count *int64, userID uuid.UUID) error {
	return a.client.DB.Model(&model.BandMemberRole{}).
		Where(model.BandMemberRole{UserID: userID}).
		Count(count).
		Error
}

func (a artistRepository) CreateBandMemberRole(bandMemberRole *model.BandMemberRole) error {
	return a.client.DB.Create(&bandMemberRole).Error
}

func (a artistRepository) UpdateBandMemberRole(bandMemberRole *model.BandMemberRole) error {
	return a.client.DB.Save(&bandMemberRole).Error
}

func (a artistRepository) UpdateAllBandMemberRoles(bandMemberRoles *[]model.BandMemberRole) error {
	return a.client.DB.Transaction(func(tx *gorm.DB) error {
		for _, sectionType := range *bandMemberRoles {
			if err := tx.Save(sectionType).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (a artistRepository) DeleteBandMemberRole(id uuid.UUID) error {
	return a.client.DB.Delete(&model.BandMemberRole{}, id).Error
}
