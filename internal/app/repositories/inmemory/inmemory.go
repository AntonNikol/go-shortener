package inmemory

import (
	"github.com/AntonNikol/go-shortener/internal/app/models"
	"github.com/AntonNikol/go-shortener/internal/app/repositories"
	"log"
)

type Repository struct {
	items map[string]models.Item
}

func New() *Repository {
	return &Repository{
		items: map[string]models.Item{},
	}
}

func (r *Repository) AddItem(item models.Item) (models.Item, error) {
	// добавляем в мапу items
	r.items[item.ID] = item
	return item, nil
}

func (r *Repository) GetItemByID(id string) (models.Item, error) {
	log.Println("GetItemById file")

	// проверяем мапу на наличие там айтема по ключу
	if res, ok := r.items[id]; ok {
		log.Printf("Результат найден в мапе")
		return res, nil
	}

	return models.Item{}, repositories.ErrNotFound
}
