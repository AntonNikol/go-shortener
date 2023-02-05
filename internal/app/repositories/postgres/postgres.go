package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/AntonNikol/go-shortener/internal/app/models"
	"github.com/AntonNikol/go-shortener/internal/app/repositories"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"log"
)

var err error

type Postgres struct {
	DB *sql.DB
}

func (p Postgres) AddItem(item models.Item) (models.Item, error) {

	var id string
	err := p.DB.QueryRow("INSERT INTO short_links (full_url, short_url, user_id) values ($1, $2, $3) RETURNING id",
		item.FullURL, item.ShortURL, item.UserID).Scan(&id)
	if err != nil {
		log.Printf("postgres AddItem ошибка id: %s", id)

		return models.Item{}, err
	}
	log.Printf("postgres AddItem успешно id: %s", id)

	item.ID = id
	return item, nil
}

func (p Postgres) GetItemByID(id string) (models.Item, error) {
	//TODO: откуда тут правильно тянуть контекст?!
	row := p.DB.QueryRowContext(context.Background(),
		"SELECT id,full_url, short_url FROM short_links where id=$1", id)
	// обязательно закрываем перед возвратом функции
	//defer rows.Close()

	var i models.Item

	//// пробегаем по всем записям
	//for rows.Next() {
	//	err = rows.Scan(&i.ID, &i.FullURL, &i.ShortURL)
	//	if err != nil {
	//		log.Printf("postgres GetItemByID Scan ошибка: %v", err)
	//
	//		return models.Item{}, err
	//	}
	//}
	err = row.Scan(&i.ID, &i.FullURL, &i.ShortURL)
	if err != nil {
		log.Printf("postgres GetItemByID Scan ошибка: %v", err)

		return models.Item{}, repositories.ErrNotFound
	}
	if i.ShortURL == "" {
		return models.Item{}, repositories.ErrNotFound
	}

	log.Printf("postgres GetItemByID Scan успех: %v", i)

	return i, nil
}

func (p Postgres) GetItemsByUserID(userID string) ([]models.ItemResponse, error) {
	var res []models.ItemResponse

	//TODO: откуда тут правильно тянуть контекст?!
	rows, err := p.DB.QueryContext(context.Background(),
		"SELECT id,full_url, short_url, user_id FROM short_links where user_id=$1", userID)

	if err != nil {
		return nil, err
	}
	// обязательно закрываем перед возвратом функции
	defer rows.Close()

	// пробегаем по всем записям

	for rows.Next() {
		i := models.ItemResponse{}
		err = rows.Scan(&i.ID, &i.FullURL, &i.ShortURL, &i.UserID)
		if err != nil {
			return nil, err
		}

		res = append(res, i)
	}

	log.Printf("postgres GetItemsByUserID Scan успех: %v", res)

	// проверяем на ошибки
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	if len(res) > 0 {
		return res, nil
	}
	return res, repositories.ErrNotFound

}

func New(ctx context.Context, DSN string) *Postgres {
	db, err := sql.Open("postgres",
		DSN)
	if err != nil {
		panic(err)
	}

	// накатка миграций
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./internal/migrations",
		"postgres", driver)
	if err != nil {
		panic(err)
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(err)
	}

	return &Postgres{DB: db}
}

func (p Postgres) Ping(ctx context.Context) error {
	return p.DB.Ping()
}

func (p Postgres) AddItemsList(items map[string]models.Item) (map[string]models.Item, error) {

	result := map[string]models.Item{}
	ctx := context.Background()

	// шаг 1 — объявляем транзакцию
	tx, err := p.DB.Begin()
	if err != nil {
		return nil, err
	}
	// шаг 1.1 — если возникает ошибка, откатываем изменения
	defer tx.Rollback()

	// шаг 2 — готовим инструкцию
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO short_links(full_url) VALUES($1) RETURNING id")
	if err != nil {
		return nil, err
	}
	// шаг 2.1 — не забываем закрыть инструкцию, когда она больше не нужна
	defer stmt.Close()

	var id string
	for k, v := range items {
		// шаг 3 — указываем, что каждое видео будет добавлено в транзакцию
		err := stmt.QueryRowContext(ctx, v.FullURL).Scan(&id)
		if err != nil {
			return nil, err
		}
		result[k] = models.Item{ID: id}
	}
	// шаг 4 — сохраняем изменения
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return result, nil
}
