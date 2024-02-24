package database

import (
	"github.com/jackc/pgx"
)

func getName() {
	connConfig := pgx.ConnConfig{
		Host:     "34.170.5.185",
		Port:     5432,
		Database: "rapidtransfer",
		User:     "postgres",
		Password: "postgres",
	}
	conn, err := pgx.Connect(connConfig)
	// context.Background(), os.Getenv(dbAddress)
}
