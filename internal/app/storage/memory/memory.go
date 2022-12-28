package memory

import (
	"errors"
	"github.com/AntonNikol/go-shortener/internal/app/models"
)

type Storage struct {
	items []models.Item
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Add(item models.Item) (models.Item, error) {
	s.items = append(s.items, item)
	return item, nil
}

func (s *Storage) Get(id string) (models.Item, error) {
	for _, item := range s.items {
		if item.ID == id {
			return item, nil
		}
	}

	return models.Item{}, errors.New("not found")
}
