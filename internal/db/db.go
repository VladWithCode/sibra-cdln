package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func Connect() (*pgxpool.Pool, error) {
	var err error
	db, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetPool() (*pgxpool.Conn, error) {
	conn, err := db.Acquire(context.Background())

	if err != nil {
		return nil, err
	}

	return conn, nil
}

func GetTxAndPool(ctx context.Context) (pgx.Tx, *pgxpool.Conn, error) {
	conn, err := db.Acquire(context.Background())

	if err != nil {
		return nil, nil, err
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}

	return tx, conn, nil
}
