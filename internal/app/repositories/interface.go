package repositories

import "github.com/AntonNikol/go-shortener/internal/app/models"

type RepositoryInterface interface {
	AddItem(item models.Item) (models.Item, error)
	GetItemById(id string) (models.Item, error)
}
