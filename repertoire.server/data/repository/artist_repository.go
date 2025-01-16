package repository

import (
	"repertoire/server/data/database"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type ArtistRepository interface {
	Get(artist *model.Artist, id uuid.UUID) error
	// Deprecated: Use normal Get instead
	GetWithAssociations(artist *model.Artist, id uuid.UUID) error
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
	Delete(id uuid.UUID) error
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
	return a.client.DB.Find(&artist, model.Artist{ID: id}).Error
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

func (a artistRepository) Delete(id uuid.UUID) error {
	return a.client.DB.Delete(&model.Artist{}, id).Error
}
