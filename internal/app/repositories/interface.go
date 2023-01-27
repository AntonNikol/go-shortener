package repositories

import "github.com/AntonNikol/go-shortener/internal/app/models"

type Repository interface {
	AddItem(item models.Item) (models.Item, error)
	GetItemByID(id string) (models.Item, error)
}

type ErrNotFound struct {
	msg string // description of error
}

func (m *ErrNotFound) Error() string {
	return "Not Found"
}
