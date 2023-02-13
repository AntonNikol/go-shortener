package file

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AntonNikol/go-shortener/internal/app/models"
	"github.com/AntonNikol/go-shortener/internal/app/repositories"
	"github.com/AntonNikol/go-shortener/pkg/generator"
	"log"
	"os"
)

// встраивание реализации inmemory в file
// При создании структуры вкладываем в нее интерфейс,
// далее вызывая в хенделере метод repository.GetItemByID
// будет вызываться метод repositories.inmemory.GetItemByID
// Если реализаций интерфейса будет больше, то в структуру необходимо класть поле repository: inmemory
// и в file.go вызывать метод r.repository.GetItemByID

type Repository struct {
	//repositories.Repository
	items    map[string]models.Item
	file     *os.File
	filename string
}

func New(filename string) *Repository {
	// Открываем файл
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}

	//internal := inmemory.New()
	items := map[string]models.Item{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		item := models.Item{}

		if err := json.Unmarshal(scanner.Bytes(), &item); err != nil {
			panic(err)
		}
		//internal.AddItem(context.Background(), item)
		items[item.ID] = item

		log.Printf("Построчное чтение, item : %s", item)
	}

	return &Repository{
		filename: filename,
		file:     file,
		//Repository: internal,
		items: items,
	}
}

func (r *Repository) AddItem(ctx context.Context, item models.Item) (models.Item, error) {
	id, err := generator.GenerateRandomID(3)
	if err != nil {
		return models.Item{}, err
	}
	item.ID = id

	data, err := json.Marshal(item)
	if err != nil {
		return models.Item{}, fmt.Errorf("unable serialise item %w", err)
	}

	//// Тут добавляем item в Repository который положили в структуру (в этом случае inmemory),
	//// а в файл просто пишем новые строки чтобы при инициализации репозитория после перезапуска сервера
	//// заполнить мапу данными
	//_, err = r.AddItem(ctx, item)
	//if err != nil {
	//	return models.Item{}, fmt.Errorf("unable to add item to internal repository: %w", err)
	//}
	//  кладем данные в память
	r.items[id] = item

	// пишем в буфер
	writer := bufio.NewWriter(r.file)
	_, err = writer.Write(data)
	if err != nil {
		return models.Item{}, fmt.Errorf("unable to write file: %w", err)
	}

	// добавляем перенос строки
	_, err = writer.Write([]byte("\n"))
	if err != nil {
		return models.Item{}, fmt.Errorf("unable to write file: %w", err)
	}

	log.Printf("Запись в файл произведена")

	// записываем буфер в файл
	writer.Flush()
	return item, nil
}

func (r *Repository) AddItemsList(ctx context.Context, items map[string]models.Item) (map[string]models.Item, error) {

	writer := bufio.NewWriter(r.file)

	newItems := map[string]models.Item{}
	// добавляем в мапу items

	for k, i := range items {
		id, _ := generator.GenerateRandomID(3)
		newItem := models.Item{
			ID:      id,
			FullURL: i.FullURL,
		}
		newItems[k] = newItem
		log.Printf("file AddItemsList добавляем item в новую мапу newItem %v", newItem)
		r.items[id] = newItem

		data, err := json.Marshal(newItem)
		if err != nil {
			return nil, fmt.Errorf("unable serialise item %w", err)
		}

		_, err = writer.Write(data)
		if err != nil {
			return nil, fmt.Errorf("unable to write file: %w", err)
		}

		// добавляем перенос строки
		_, err = writer.Write([]byte("\n"))
		if err != nil {
			return nil, fmt.Errorf("unable to write file: %w", err)
		}

		writer.Flush()
	}
	// добавляем newItems в items
	return newItems, nil
}

func (r *Repository) GetItemByID(ctx context.Context, id string) (models.Item, error) {
	log.Println("GetItemById file")

	// проверяем мапу на наличие там айтема по ключу
	if res, ok := r.items[id]; ok {
		log.Printf("Результат найден в мапе")
		return res, nil
	}

	return models.Item{}, repositories.ErrNotFound
}

func (r *Repository) GetItemsByUserID(ctx context.Context, userID string) ([]models.ItemResponse, error) {

	log.Println("GetItemsByUserID file")

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

func (r *Repository) Ping(ctx context.Context) error {
	return nil
}
