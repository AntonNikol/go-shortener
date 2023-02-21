package inmemory

import (
	"context"
	"errors"
	"github.com/AntonNikol/go-shortener/internal/app/models"
	"github.com/AntonNikol/go-shortener/internal/app/repositories"
	"github.com/AntonNikol/go-shortener/pkg/generator"
	"sync"
)

type Repository struct {
	items map[string]models.Item
}

func New() *Repository {
	return &Repository{
		items: map[string]models.Item{},
	}
}

func (r *Repository) Ping(ctx context.Context) error {
	return nil
}

func (r *Repository) AddItem(ctx context.Context, item models.Item) (*models.Item, error) {
	// добавляем в мапу items
	id, _ := generator.GenerateRandomID(3)
	item.ID = id
	r.items[item.ID] = item
	return &item, nil
}

func (r *Repository) GetItemByID(ctx context.Context, id string) (*models.Item, error) {
	// проверяем мапу на наличие там айтема по ключу
	if res, ok := r.items[id]; ok {
		return &res, nil
	}

	return nil, repositories.ErrNotFound
}

func (r *Repository) GetItemsByUserID(ctx context.Context, userID string) ([]models.ItemResponse, error) {
	res := make([]models.ItemResponse, 0)
	// проверяем мапу на наличие там айтема с userID
	for _, v := range r.items {
		if v.UserID == userID {
			res = append(res, models.ItemResponse{FullURL: v.FullURL, ID: v.ID})
		}
	}
	if len(res) == 0 {
		return res, errors.New("items not found")
	}

	return res, nil
}

func (r *Repository) AddItemsList(ctx context.Context, items map[string]models.Item) (map[string]models.Item, error) {
	newItems := map[string]models.Item{}
	// добавляем в мапу items

	for k, i := range items {
		id, _ := generator.GenerateRandomID(3)
		newItem := models.Item{
			ID:      id,
			FullURL: i.FullURL,
		}
		newItems[k] = newItem
		r.items[id] = newItem
	}
	// добавляем newItems в items
	return newItems, nil
}

func (r *Repository) Delete(ctx context.Context, list []string, userID string) error {
	var wg sync.WaitGroup

	for _, v := range list {
		wg.Add(1)

		go func(key string) {
			tmp, ok := r.items[key]
			if ok && tmp.UserID == userID {
				tmp.IsDeleted = true
				r.items[key] = tmp
			}
			wg.Done()
		}(v)
	}

	wg.Wait()
	return nil
}
