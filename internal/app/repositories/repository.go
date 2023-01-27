package repositories

import (
	"errors"
	"github.com/AntonNikol/go-shortener/internal/app/models"
)

var (
	ErrNotFound = errors.New("not found")
)

type Repository interface {
	AddItem(item models.Item) (models.Item, error)
	GetItemByID(id string) (models.Item, error)
}
