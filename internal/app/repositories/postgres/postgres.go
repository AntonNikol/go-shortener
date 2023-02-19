package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/AntonNikol/go-shortener/internal/app/models"
	"github.com/AntonNikol/go-shortener/internal/app/repositories"
	"github.com/AntonNikol/go-shortener/pkg/generator"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"sync"
)

var err error

type Postgres struct {
	DB *sql.DB
}

func New(ctx context.Context, DSN string) (*Postgres, error) {
	db, err := sql.Open("pgx",
		DSN)
	if err != nil {
		return nil, err
	}

	// накатка миграций
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./internal/migrations",
		"postgres", driver)
	if err != nil {
		return nil, err
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}

	return &Postgres{DB: db}, nil
}

func (p Postgres) AddItem(ctx context.Context, item models.Item) (*models.Item, error) {
	var id string
	shortURL, _ := generator.GenerateRandomID(3)

	err := p.DB.QueryRowContext(ctx, "INSERT INTO short_links (full_url, user_id, short_url) values ($1, $2, $3) "+
		//"ON CONFLICT (full_url) DO NOTHING"+
		"  RETURNING id ",
		item.FullURL, item.UserID, shortURL).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			log.Printf("да это ошибка pgerr")
			//Получаем запись по full_url
			log.Printf("postgres AddItem получаем запись по полному URL: %v, %v", item, err)
			item, err = p.GetItemByFullURL(ctx, item.FullURL)
			if err != nil {
				return &item, fmt.Errorf("failed to retrieve conflicting row in db: %w", repositories.ErrNotFound)
			}

			return &item, repositories.ErrAlreadyExists
		}
		return nil, fmt.Errorf("unable insert into table %w", err)
	}

	log.Printf("postgres AddItem успешно id: %s", id)

	item.ID = shortURL
	fmt.Printf("возвращаемый item в методе createItem %+v", item)
	return &item, nil
}

func (p Postgres) GetItemByFullURL(ctx context.Context, fullURL string) (models.Item, error) {
	row := p.DB.QueryRowContext(ctx,
		"SELECT short_url,full_url FROM short_links where full_url=$1", fullURL)

	var i models.Item

	err = row.Scan(&i.ID, &i.FullURL)
	if err != nil {
		log.Printf("postgres GetItemByID Scan ошибка: %v", err)

		return models.Item{}, repositories.ErrNotFound
	}
	return i, nil
}

func (p Postgres) GetItemByID(ctx context.Context, id string) (*models.Item, error) {
	row := p.DB.QueryRowContext(ctx,
		"SELECT short_url,full_url,is_deleted FROM short_links where short_url=$1", id)

	var i models.Item

	err = row.Scan(&i.ID, &i.FullURL, &i.IsDeleted)
	if err != nil {
		log.Printf("postgres GetItemByID Scan ошибка: %v", err)

		return nil, repositories.ErrNotFound
	}

	log.Printf("postgres GetItemByID Scan успех: %+v", i)

	return &i, nil
}

func (p Postgres) GetItemsByUserID(ctx context.Context, userID string) ([]models.ItemResponse, error) {
	var res []models.ItemResponse
	rows, err := p.DB.QueryContext(ctx,
		"SELECT id,full_url, user_id FROM short_links where user_id=$1", userID)

	if err != nil {
		return nil, err
	}
	// обязательно закрываем перед возвратом функции
	defer rows.Close()

	// пробегаем по всем записям

	for rows.Next() {
		i := models.ItemResponse{}
		err = rows.Scan(&i.ID, &i.FullURL, &i.UserID)
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

func (p Postgres) Ping(ctx context.Context) error {
	return p.DB.Ping()
}

func (p Postgres) AddItemsList(ctx context.Context, items map[string]models.Item) (map[string]models.Item, error) {

	result := map[string]models.Item{}

	// шаг 1 — объявляем транзакцию
	tx, err := p.DB.Begin()
	if err != nil {
		return nil, err
	}
	// шаг 1.1 — если возникает ошибка, откатываем изменения
	defer tx.Rollback()

	// шаг 2 — готовим инструкцию
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO short_links(full_url,short_url,user_id) VALUES($1, $2, $3) RETURNING id")
	if err != nil {
		return nil, err
	}
	// шаг 2.1 — не забываем закрыть инструкцию, когда она больше не нужна
	defer stmt.Close()

	var id string
	for k, v := range items {
		shortURL, _ := generator.GenerateRandomID(3)
		// шаг 3 — указываем, что каждый item будет добавлен в транзакцию
		err := stmt.QueryRowContext(ctx, v.FullURL, shortURL, v.UserID).Scan(&id)
		if err != nil {
			return nil, err
		}
		result[k] = models.Item{ID: shortURL}
	}
	// шаг 4 — сохраняем изменения
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p Postgres) Delete(ctx context.Context, list []string, userID string) (*int, error) {
	ids := &pgtype.VarcharArray{}
	ids.Set(list)

	var wg sync.WaitGroup

	stmt := "UPDATE short_links SET is_deleted = true WHERE short_url = any ($1) AND user_id= $2"
	wg.Add(1)

	go func(ids *pgtype.VarcharArray) {
		res, err := p.DB.Exec(stmt, ids, userID)
		if err != nil {
			log.Printf("unable update rows: %v", err)
		}
		count, err := res.RowsAffected()
		if err != nil {
			log.Printf("unable check rows affecte: %v", err)
		}
		log.Printf("rows updated: %d", count)

		wg.Done()
	}(ids)

	wg.Wait()

	result := 1
	return &result, nil
}
