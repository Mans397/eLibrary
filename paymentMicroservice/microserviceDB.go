// paymentMicroservice/microserviceDB.go
package paymentMicroservice

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var microDB *gorm.DB

type MicroTransaction struct {
	gorm.Model
	TransactionID uint // ID транзакции из основного сервера
	CustomerID    uint
	CustomerName  string
	CustomerEmail string
	Status        string
	CreatedAt     string
}

func initMicroserviceDB() {
	// Можно использовать ту же самую базу eLibrary
	// Важно, чтобы таблица MicroTransaction не имела внешних ключей
	dsn := "host=localhost port=5433 user=postgres password=rootroot dbname=eLibrary sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect microservice DB: %v", err)
	}
	microDB = db

	// Миграция таблицы MicroTransaction
	if err := microDB.AutoMigrate(&MicroTransaction{}); err != nil {
		fmt.Println("Error migrating MicroTransaction table:", err)
	}
}
