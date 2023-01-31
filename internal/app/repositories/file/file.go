package file

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/AntonNikol/go-shortener/internal/app/models"
	"github.com/AntonNikol/go-shortener/internal/app/repositories"
	"github.com/AntonNikol/go-shortener/internal/app/repositories/inmemory"
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
	repositories.Repository
	//items    map[string]models.Item
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
		internal.AddItem(item)

		log.Printf("Построчное чтение, item : %s", item)
	}

	return &Repository{
		filename:   filename,
		file:       file,
		Repository: internal,
	}
}

func (r *Repository) Ping(ctx context.Context) error {
	return nil
}

func (r *Repository) AddItem(item models.Item) (models.Item, error) {
	data, err := json.Marshal(item)
	if err != nil {
		return models.Item{}, fmt.Errorf("unable serialise item %w", err)
	}

	// Тут добавляем item в Repository который положили в структуру (в этом случае inmemory),
	// а в файл просто пишем новые строки чтобы при инициализации репозитория после перезапуска сервера
	// заполнить мапу данными
	_, err = r.Repository.AddItem(item)
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
