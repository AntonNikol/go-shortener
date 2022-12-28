package inmemory

import (
	"github.com/AntonNikol/go-shortener/internal/app/models"
	"github.com/AntonNikol/go-shortener/internal/app/storage/memory"
)

type Repository struct {
	db memory.Storage
}

func New(db *memory.Storage) *Repository {
	return &Repository{db: *db}
}

func (r *Repository) AddItem(item models.Item) (models.Item, error) {
	return r.db.Add(item)
}

func (r *Repository) GetItemByID(id string) (models.Item, error) {
	return r.db.Get(id)
}
