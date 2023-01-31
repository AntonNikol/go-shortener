package inmemory

import (
	"context"
	"errors"
	"github.com/AntonNikol/go-shortener/internal/app/models"
	"github.com/AntonNikol/go-shortener/internal/app/repositories"
	"log"
)

type Repository struct {
	items map[string]models.Item
}

func (r *Repository) Ping(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
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
	log.Println("GetItemById memory")

	// проверяем мапу на наличие там айтема по ключу
	if res, ok := r.items[id]; ok {
		log.Printf("Результат найден в мапе")
		return res, nil
	}

	return models.Item{}, repositories.ErrNotFound
}

func (r *Repository) GetItemsByUserID(userID string) ([]models.ItemResponse, error) {
	log.Println("GetItemsByUserID memory")

	res := make([]models.ItemResponse, 0)
	// проверяем мапу на наличие там айтема с userID
	for _, v := range r.items {
		if v.UserID == userID {
			res = append(res, models.ItemResponse{ShortURL: v.ShortURL, FullURL: v.FullURL})
		}
	}
	if len(res) == 0 {
		return res, errors.New("items not found")
	}

	return res, nil
}
