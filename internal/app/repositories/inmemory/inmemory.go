package inmemory

import (
	"context"
	"errors"
	"fmt"
	"github.com/AntonNikol/go-shortener/internal/app/models"
	"github.com/AntonNikol/go-shortener/internal/app/repositories"
	"log"
	"math/rand"
	"strconv"
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

func (r *Repository) AddItem(item models.Item) (models.Item, error) {
	// добавляем в мапу items
	r.items[item.ID] = item
	return item, nil
}

func (r *Repository) GetItemByID(id string) (models.Item, error) {
	log.Println("GetItemById memory")

	// проверяем мапу на наличие там айтема по ключу
	if res, ok := r.items[id]; ok {
		log.Printf("Результат найден в мапе")
		return res, nil
	}

	return models.Item{}, repositories.ErrNotFound
}

func (r *Repository) GetItemsByUserID(userID string) ([]models.ItemResponse, error) {
	log.Println("GetItemsByUserID memory")

	res := make([]models.ItemResponse, 0)
	// проверяем мапу на наличие там айтема с userID
	for _, v := range r.items {
		if v.UserID == userID {
			res = append(res, models.ItemResponse{ShortURL: v.ShortURL, FullURL: v.FullURL})
		}
	}
	if len(res) == 0 {
		return res, errors.New("items not found")
	}

	return res, nil
}

func (r *Repository) AddItemsList(items map[string]models.Item) (map[string]models.Item, error) {
	//return items, nil
	//TODO: implements me
	panic(errors.New("inmerory AddItemsList"))
}

// Получение рандомного id
func (r *Repository) generateUniqueItemID(id string) (string, error) {
	randomInt := rand.Intn(999999)
	randomString := strconv.Itoa(randomInt)

	log.Printf("generateUniqueItemID Получение рандомного id: %s", id)
	exist, err := r.checkItemExist(randomString)
	if err != nil {
		return "", fmt.Errorf("unable to check item exist item by id: %w", err)
	}

	log.Printf("generateUniqueItemID exists id: %v", exist)

	if randomString != id && !exist {
		return randomString, nil
	}

	return r.generateUniqueItemID(randomString)
}

// Проверка есть ли в файле item с таким id
func (r *Repository) checkItemExist(id string) (bool, error) {
	_, err := r.GetItemByID(id)

	// проверяем что ошибка не пустая и она не нот фаунд
	if err != nil && !errors.Is(err, repositories.ErrNotFound) {
		return false, fmt.Errorf("unable to get item by id: %w", err)
	}
	return !errors.Is(err, repositories.ErrNotFound), nil
}
