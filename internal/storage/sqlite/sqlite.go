package sqlite

import (
	"database/sql"

	"github.com/Aytaditya/students-api-golang/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

// we will call this inside main
func New(cfg *config.Config) (*Sqlite, error) {

	db, err := sql.Open("sqlite3", cfg.StoragePath) // open sqlite database connection first argument is type of database and second is path
	if err != nil {
		return nil, err // returning 2 things nil and error because there is error
	}

	// we need to install the driver manually
	_, er := db.Exec(`CREATE TABLE IF NOT EXISTS students(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		age INTEGER NOT NULL
	)`)
	if er != nil {
		return nil, er
	}

	return &Sqlite{
		Db: db}, nil

}
