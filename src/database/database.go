package database

import (
	"fmt"

	"github.com/jackc/pgx"
)

// Create target connection for the database
func GetConn() *pgx.Conn {
	connConfig := pgx.ConnConfig{
		Host:     "34.170.5.185",
		Port:     5432,
		Database: "rapidtransfer",
		User:     "postgres",
		Password: "postgres",
	}
	conn, err := pgx.Connect(connConfig)
	if err != nil {
		fmt.Print("Failed to connect")
	}
	return conn
}

// Creates a public connection that all functions can use
var conn *pgx.Conn = GetConn()

// Inits all of the tables for the database
func InitializeDatabase() {

	conn.Exec(`
	CREATE TABLE IF NOT EXISTS transfer (id SERIAL PRIMARY KEY, userFrom INT NOT NULL, userTo INT NOT NULL, keyword VARCHAR(100), address VARCHAR(100), filename VARCHAR(100));
	CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL UNIQUE, keyword VARCHAR(100), macaddr VARCHAR(100));
	CREATE TABLE IF NOT EXISTS friends (orig_user INT NOT NULL, friend_id INT NOT NULL, total_transfers INT NOT NULL DEFAULT 0);
	`)

}
