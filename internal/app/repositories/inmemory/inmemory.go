package inmemory

import (
	"errors"
	"github.com/AntonNikol/go-shortener/internal/app/models"
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

	return models.Item{}, errors.New("not found")
}
