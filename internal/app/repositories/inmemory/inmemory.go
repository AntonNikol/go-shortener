package inmemory

import (
	"context"
	"errors"
	"github.com/AntonNikol/go-shortener/internal/app/models"
	"github.com/AntonNikol/go-shortener/internal/app/repositories"
	"github.com/AntonNikol/go-shortener/pkg/generator"
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
			UserID:  i.UserID,
		}
		newItems[k] = newItem
		r.items[id] = newItem
	}
	// добавляем newItems в items
	return newItems, nil
}

func (r *Repository) Delete(ctx context.Context, list []string, userID string) error {
	for _, v := range list {
		log.Printf("inmemory Delete id: %s", v)
		tmp, ok := r.items[v]
		log.Printf("inmemory Delete tmp:%v, ok %v", v, ok)
		log.Printf("inmemory Delete tmp:%+v, userID %v", tmp, userID)

		if ok && tmp.UserID == userID {
			//if ok {
			log.Printf("inmemory Delete условие if: %v", v)
			tmp.IsDeleted = true
			r.items[v] = tmp
		}
	}

	return nil
}
