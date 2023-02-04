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
)

var err error

type Postgres struct {
	DB *sql.DB
}

func (p Postgres) AddItem(item models.Item) (models.Item, error) {

	_, err := p.DB.Exec("INSERT INTO short_links (full_url, short_url, user_id) values ($1, $2, $3)",
		item.FullURL, item.ShortURL, item.UserID)
	if err != nil {
		return models.Item{}, err
	}
	return item, nil
}

func (p Postgres) GetItemByID(id string) (models.Item, error) {
	//TODO: откуда тут правильно тянуть контекст?!
	rows, err := p.DB.QueryContext(context.Background(),
		"SELECT id,full_url, short_url FROM short_links where short_url=$1", id)
	// обязательно закрываем перед возвратом функции
	defer rows.Close()

	var i models.Item

	// пробегаем по всем записям
	for rows.Next() {
		err = rows.Scan(&i.ID, &i.FullURL, &i.ShortURL)
		if err != nil {
			return models.Item{}, err
		}
	}

	return i, repositories.ErrNotFound
}

func (p Postgres) GetItemsByUserID(userID string) ([]models.ItemResponse, error) {
	var res []models.ItemResponse

	//TODO: откуда тут правильно тянуть контекст?!
	rows, err := p.DB.QueryContext(context.Background(),
		"SELECT id,full_url, short_url, user_id FROM short_links where user_id=$1", userID)
	// обязательно закрываем перед возвратом функции
	defer rows.Close()

	i := models.ItemResponse{}
	// пробегаем по всем записям
	for rows.Next() {
		err = rows.Scan(&i.ID, &i.FullURL, &i.ShortURL, &i.UserID)
		if err != nil {
			return nil, err
		}
	}
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
