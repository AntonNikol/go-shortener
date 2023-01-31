package repositories

import (
	"context"
	"errors"
	"github.com/AntonNikol/go-shortener/internal/app/models"
)

var (
	ErrNotFound = errors.New("not found")
)

type Repository interface {
	AddItem(item models.Item) (models.Item, error)
	GetItemByID(id string) (models.Item, error)
	GetItemsByUserID(userID string) ([]models.ItemResponse, error)
	Ping(ctx context.Context) error
}
