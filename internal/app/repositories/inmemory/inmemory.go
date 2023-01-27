package inmemory

import (
	"github.com/AntonNikol/go-shortener/internal/app/models"
	"github.com/AntonNikol/go-shortener/internal/app/repositories"
)

type Repository struct {
	items []models.Item
}

func New() *Repository {
	return &Repository{}
}

func (r *Repository) AddItem(item models.Item) (models.Item, error) {
	r.items = append(r.items, item)
	return item, nil
}

func (r *Repository) GetItemByID(id string) (models.Item, error) {
	for _, item := range r.items {
		if item.ID == id {
			return item, nil
		}
	}

	return models.Item{}, repositories.ErrNotFound
}
