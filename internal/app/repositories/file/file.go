package file

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AntonNikol/go-shortener/internal/app/models"
	"github.com/AntonNikol/go-shortener/internal/app/repositories"
	"github.com/AntonNikol/go-shortener/internal/app/repositories/inmemory"
	"log"
	"math/rand"
	"os"
	"strconv"
)

// встраивание реализации inmemory в file
// При создании структуры вкладываем в нее интерфейс,
// далее вызывая в хенделере метод repository.GetItemByID
// будет вызываться метод repositories.inmemory.GetItemByID
// Если реализаций интерфейса будет больше, то в структуру необходимо класть поле repository: inmemory
// и в file.go вызывать метод r.repository.GetItemByID

type Repository struct {
	repositories.Repository
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

	internal := inmemory.New()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		item := models.Item{}

		if err := json.Unmarshal(scanner.Bytes(), &item); err != nil {
			panic(err)
		}
		internal.AddItem(context.Background(), item)

		log.Printf("Построчное чтение, item : %s", item)
	}

	return &Repository{
		filename:   filename,
		file:       file,
		Repository: internal,
	}
}

func (r *Repository) AddItem(ctx context.Context, item models.Item) (models.Item, error) {
	id, err := r.generateUniqueItemID(ctx, "")
	if err != nil {
		return models.Item{}, err
	}
	item.ID = id
	item.ShortURL = item.ShortURL + id

	data, err := json.Marshal(item)
	if err != nil {
		return models.Item{}, fmt.Errorf("unable serialise item %w", err)
	}

	// Тут добавляем item в Repository который положили в структуру (в этом случае inmemory),
	// а в файл просто пишем новые строки чтобы при инициализации репозитория после перезапуска сервера
	// заполнить мапу данными
	_, err = r.Repository.AddItem(ctx, item)
	if err != nil {
		return models.Item{}, fmt.Errorf("unable to add item to internal repository: %w", err)
	}

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

// Получение рандомного id
func (r *Repository) generateUniqueItemID(ctx context.Context, id string) (string, error) {
	randomInt := rand.Intn(999999)
	randomString := strconv.Itoa(randomInt)

	log.Printf("generateUniqueItemID Получение рандомного id: %s", id)
	exist, err := r.checkItemExist(ctx, randomString)
	if err != nil {
		return "", fmt.Errorf("unable to check item exist item by id: %w", err)
	}

	log.Printf("generateUniqueItemID exists id: %v", exist)

	if randomString != id && !exist {
		return randomString, nil
	}

	return r.generateUniqueItemID(ctx, randomString)
}

// Проверка есть ли в файле item с таким id
func (r *Repository) checkItemExist(ctx context.Context, id string) (bool, error) {
	_, err := r.GetItemByID(ctx, id)

	// проверяем что ошибка не пустая и она не нот фаунд
	if err != nil && !errors.Is(err, repositories.ErrNotFound) {
		return false, fmt.Errorf("unable to get item by id: %w", err)
	}
	return !errors.Is(err, repositories.ErrNotFound), nil
}
