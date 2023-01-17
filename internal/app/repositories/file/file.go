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

	// пишем в буфер
	writer := bufio.NewWriter(file)
	_, err = writer.Write(data)
	if err != nil {
		panic(err)
	}

	// добавляем перенос строки
	_, err = writer.Write([]byte("\n"))
	if err != nil {
		panic(err)
	}

	log.Printf("Запись в файл произведена")

	// записываем буфер в файл
	writer.Flush()
	return item, nil
}

func (r *Repository) GetItemByID(id string) (models.Item, error) {
	log.Println("GetItemById file")

	// проверяем мапу на наличие там айтема по ключу
	if res, ok := r.items[id]; ok {
		log.Printf("Результат найден в мапе")
		return res, nil
	}

	file, err := os.OpenFile(r.filename, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		item := models.Item{}

		if err = json.Unmarshal(scanner.Bytes(), &item); err != nil {
			panic(err)
		}
		log.Printf("Построчное чтение, item : %s", item)

		if item.ID == id {
			log.Println("результат найден в файле")
			return item, nil
		}
	}

	return models.Item{}, errors.New("not found")
}
