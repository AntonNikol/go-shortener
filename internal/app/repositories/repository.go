package repositories

import "github.com/AntonNikol/go-shortener/internal/app/models"

type Repository struct {
	repo RepositoryInterface
}

func NewRepository(repo RepositoryInterface) *Repository {
	return &Repository{repo: repo}
}

func (r Repository) AddItem(item models.Item) (models.Item, error) {
	return r.repo.AddItem(item)
}

func (r Repository) GetItemByID(id string) (models.Item, error) {
	res, err := r.repo.GetItemByID(id)
	return res, err
}
