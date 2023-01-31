package postgres

import (
	"context"
	"github.com/AntonNikol/go-shortener/internal/app/models"
	"github.com/jackc/pgx/v5"
	"log"
)

type Postgres struct {
	DB *pgx.Conn
}

func (p Postgres) AddItem(item models.Item) (models.Item, error) {
	//TODO implement me
	panic("implement me")
}

func (p Postgres) GetItemByID(id string) (models.Item, error) {
	//TODO implement me
	panic("implement me")
}

func (p Postgres) GetItemsByUserID(userID string) ([]models.ItemResponse, error) {
	//TODO implement me
	panic("implement me")
}

func New(ctx context.Context, DSN string) *Postgres {
	//urlExample := "postgres://postgres:qwerty@localhost:5438/postgres"
	conn, err := pgx.Connect(ctx, DSN)
	if err != nil {
		//panic(err)
		log.Fatalf("unable to connect to database: %v\n", err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		log.Println("err ping")
	}

	return &Postgres{DB: conn}
}

func (p Postgres) Ping(ctx context.Context) error {
	return p.DB.Ping(ctx)

}
