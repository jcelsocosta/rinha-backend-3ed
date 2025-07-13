package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

var dbInstance *sql.DB

func openDB() (*sql.DB, error) {
	if dbInstance != nil {
		return dbInstance, nil
	}

	host := "postgres"
	port := 5432
	user := "admin"
	password := "admin"
	dbname := "postgres"

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Minute * 5)

	dbInstance = db
	return dbInstance, nil
}

func initSql(db *sql.DB) {
	sql := `
		CREATE TABLE IF NOT EXISTS payments (
			correlationid UUID PRIMARY KEY,
			amount NUMERIC NOT NULL,
			requested_at TIMESTAMP NOT NULL,
			origin VARCHAR(50)
		);
		CREATE INDEX IF NOT EXISTS payments_requested_at ON payments USING btree (requested_at);
	`

	_, err := db.Exec(sql)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func RunDB() {
	db, err := openDB()
	if err != nil {
		fmt.Println(err)
		return
	}

	initSql(db)
}
