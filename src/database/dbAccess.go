package database

import (
	"fmt"
	"os"

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

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	row := conn.QueryRow("SELECT userID FROM user")
	if err != nil {
		fmt.Println("Query failed:", err)
		return
	}
	fmt.Println(row)
}
