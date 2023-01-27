package file

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/AntonNikol/go-shortener/internal/app/models"
	"github.com/AntonNikol/go-shortener/internal/app/repositories"
	"github.com/AntonNikol/go-shortener/internal/app/repositories/inmemory"
	"log"
	"os"
)

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

	scanner := bufio.NewScanner(file)

	//items := map[string]models.Item{}
	//
	//for scanner.Scan() {
	//	item := models.Item{}
	//
	//	if err := json.Unmarshal(scanner.Bytes(), &item); err != nil {
	//		panic(err)
	//	}
	//	items[item.ID] = item
	//
	//	log.Printf("Построчное чтение, item : %s", item)
	//}

	internal := inmemory.New()

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

func (r *Repository) AddItem(item models.Item) (models.Item, error) {
	data, err := json.Marshal(item)
	if err != nil {
		panic(err)
	}

	// добавляем в мапу items чтобы можно было получить данные без запроса файла
	//r.items[item.ID] = item

	_, err = r.Repository.AddItem(item)
	if err != nil {
		return models.Item{}, fmt.Errorf("unable to add item to internal repository: %w", err)
	}

	// пишем в буфер
	writer := bufio.NewWriter(r.file)
	_, err = writer.Write(data)
	if err != nil {
		panic(err)
	}

	// добавляем перенос строки
	_, err = writer.Write([]byte("\n"))
	if err != nil {
		//panic(err)

		return models.Item{}, fmt.Errorf("ошибка записи в файл: %w", err)
	}

	log.Printf("Запись в файл произведена")

	// записываем буфер в файл
	writer.Flush()
	return item, nil
}

//func (r *Repository) GetItemByID(id string) (models.Item, error) {
//	log.Println("GetItemById file")
//
//	// проверяем мапу на наличие там айтема по ключу
//	if res, ok := r.items[id]; ok {
//		log.Printf("Результат найден в мапе")
//		return res, nil
//	}
//
//	return models.Item{}, repositories.ErrNotFound
//}
