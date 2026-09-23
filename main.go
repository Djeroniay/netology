package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "orders.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqliteRepository := NewSQLiteRepository(db)

	var initializer RepositoryInitializer = sqliteRepository
	if err := initializer.Init(); err != nil {
		log.Fatal(err)
	}

	var writer RepositoryWriter = sqliteRepository

	emailSender := NewEmailSender()
	emailService := NewOrderService(writer, emailSender)

	if err := emailService.CreateOrder(
		"Иван",
		[]string{"apple", "banana"},
		10.5,
	); err != nil {
		log.Fatal(err)
	}

	smsSender := NewSMSSender()
	smsService := NewOrderService(writer, smsSender)

	if err := smsService.CreateOrder(
		"Павел",
		[]string{"book", "notebook"},
		25.0,
	); err != nil {
		log.Fatal(err)
	}

	memoryRepository := NewMemoryRepository()

	if err := memoryRepository.Init(); err != nil {
		log.Fatal(err)
	}

	memoryService := NewOrderService(
		memoryRepository,
		emailSender,
	)

	if err := memoryService.CreateOrder(
		"Мария",
		[]string{"pen"},
		3.5,
	); err != nil {
		log.Fatal(err)
	}
}
