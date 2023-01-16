package file

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/AntonNikol/go-shortener/internal/app/models"
	"log"
	"os"
)

type Repository struct {
	items    map[string]models.Item
	file     os.File
	filename string
}

func New(filename string) *Repository {
	return &Repository{
		items:    map[string]models.Item{},
		filename: filename,
	}
}

func (r *Repository) AddItem(item models.Item) (models.Item, error) {
	data, err := json.Marshal(item)
	if err != nil {
		panic(err)
	}

	// Открываем файл
	file, err := os.OpenFile(r.filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}

	// добавляем в мапу items чтобы можно было получить данные без запроса файла
	r.items[item.ID] = item

	writer := bufio.NewWriter(file)
	_, err = writer.Write(data)
	log.Printf("Запись в файл добавлена")
	if err != nil {
		panic(err)
	}

	_, err = writer.Write([]byte("\n"))
	if err != nil {
		panic(err)
	}

	writer.Flush()
	return item, nil
}

func (r *Repository) GetItemByID(id string) (models.Item, error) {
	// проверяем мапу на наличие там айтема по ключу
	if res, ok := r.items[id]; ok {
		return res, nil
	}

	// если в мапе не найдено идем в файл
	file, err := os.OpenFile(r.filename, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)

	// если файл пустой
	if !scanner.Scan() {
		return models.Item{}, errors.New("not found")
	}

	// читаем данные из scanner
	data := scanner.Bytes()

	item := models.Item{}
	err = json.Unmarshal(data, &item)
	if err != nil {
		return models.Item{}, err
	}

	return models.Item{}, errors.New("not found")
}
