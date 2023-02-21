package repositories

import (
	"context"
	"errors"
	"github.com/AntonNikol/go-shortener/internal/app/models"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

type Repository interface {
	AddItem(ctx context.Context, item models.Item) (*models.Item, error)
	AddItemsList(ctx context.Context, items map[string]models.Item) (map[string]models.Item, error)
	GetItemByID(ctx context.Context, id string) (*models.Item, error)
	GetItemsByUserID(ctx context.Context, userID string) ([]models.ItemResponse, error)
	Ping(ctx context.Context) error
	Delete(ctx context.Context, list []string, userID string) error
}
