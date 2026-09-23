package main

import (
	"database/sql"
	"fmt"
	"strings"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db: db,
	}
}

func (r *SQLiteRepository) Init() error {
	_, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			customer TEXT NOT NULL,
			products TEXT NOT NULL,
			total REAL NOT NULL,
			status TEXT NOT NULL
		)
	`)

	return err
}

func (r *SQLiteRepository) Save(order Order) error {
	products := strings.Join(order.Products, ", ")

	_, err := r.db.Exec(
		`INSERT INTO orders
			(customer, products, total, status)
		 VALUES (?, ?, ?, ?)`,
		order.Customer,
		products,
		order.Total,
		order.Status,
	)

	return err
}

func (r *SQLiteRepository) String() string {
	return fmt.Sprintf("SQLiteRepository(%p)", r.db)
}
