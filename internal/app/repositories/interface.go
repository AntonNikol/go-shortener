package repositories

import "github.com/AntonNikol/go-shortener/internal/app/models"

type Repository interface {
	AddItem(item models.Item) (models.Item, error)
	GetItemByID(id string) (models.Item, error)
}
