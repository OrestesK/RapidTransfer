package database

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/jackc/pgx"
)

type Sql struct {
	Host     string
	Port     uint16
	Database string
	User     string
	Password string
}

var conn *pgx.Conn
var connMutex sync.Mutex

// Create target connection for the database
func GetConn() (*pgx.Conn, error) {
	connMutex.Lock()
	defer connMutex.Unlock()

	if conn != nil {
		return conn, nil
	}
	home, _ := os.UserHomeDir()
	directory := filepath.Join(home, "Rapid/.sql.json")

	var sql Sql
	err := parseSqlFile(&sql, directory)
	if err != nil {
		fmt.Println(err)
	}
	connConfig := pgx.ConnConfig{
		Host:     sql.Host,
		Port:     sql.Port,
		Database: sql.Database,
		User:     sql.User,
		Password: sql.Password,
	}

	newConn, err := pgx.Connect(connConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect: %v", err)
	}

	conn = newConn
	return conn, nil
}

func parseSqlFile(sql *Sql, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &sql)
	if err != nil {
		return err
	}
	return nil
}

// Inits all of the tables for the database
func InitializeDatabase() {
	var content embed.FS
	path, _ := content.ReadFile("database.sql")
	conn, err := GetConn()
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		os.Exit(1)
	}

	// Execute the SQL file
	_, err = conn.Exec(string(path))
	if err != nil {
		fmt.Println("Error executing SQL file:", err)
		os.Exit(1)
	}
}
