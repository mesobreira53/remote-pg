package persistence

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
)

type sql_storage struct {
	db_name string
}

type Storage interface {
	Init_instance_db(table_name string) error
	Save_pg_intance(pg_name string, pg_port int) error
	Read_all_pg_instance() (map[string]int, error)
	Check_pg_instance(pg_name string) (bool, error)
}

func NewSQLStorage(db_name string) sql_storage {
	return sql_storage{db_name: db_name}
}

func (s sql_storage) Init_instance_db(table_name string) error {
	// Open SQLite database (creates if not exists)
	db, err := sql.Open("sqlite3", s.db_name)
	if err != nil {
		return err
	}
	defer db.Close()

	// Create table
	sqlStmt := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		pg_id INTEGER PRIMARY KEY AUTOINCREMENT,
		pg_name TEXT NOT NULL,
		pg_port INT NOT NULL,
	);
	`, table_name)
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	return nil
}

func (s sql_storage) Save_pg_intance(pg_name string, pg_port int) error {
	// Open SQLite database (creates if not exists)
	db, err := sql.Open("sqlite3", s.db_name)
	if err != nil {
		return err
	}
	defer db.Close()

	// Insert data
	sqlStmt := fmt.Sprintf(`
	INSERT INTO pg_instance (pg_name, pg_port)
	VALUES ('%s', %d);
	`, pg_name, pg_port)
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	return nil
}

func (s sql_storage) Read_all_pg_instance() (map[string]int, error) {
	instances := make(map[string]int)
	// Open SQLite database (creates if not exists)
	db, err := sql.Open("sqlite3", s.db_name)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Read data
	rows, err := db.Query("SELECT pg_name, pg_port FROM pg_instance")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var pg_name string
		var pg_port int
		err = rows.Scan(&pg_name, &pg_port)
		if err != nil {
			return nil, err
		}
		instances[pg_name] = pg_port
	}
	return instances, nil
}

func (s sql_storage) Check_pg_instance(pg_name string) (bool, error) {
	// Open SQLite database (creates if not exists)
	db, err := sql.Open("sqlite3", s.db_name)
	if err != nil {
		return false, err
	}
	defer db.Close()

	// Read data
	rows, err := db.Query("SELECT pg_name FROM pg_instance WHERE pg_name = ?", pg_name)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	if rows.Next() {
		return true, nil
	}
	return false, nil
}
