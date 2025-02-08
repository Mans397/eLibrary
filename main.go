package main

import (
	// ... ваш код ...
	db "github.com/Mans397/eLibrary/Database"
	sc "github.com/Mans397/eLibrary/serverConnection"

	// ВАЖНО: импортируем пакет с микросервисом
	"github.com/Mans397/eLibrary/paymentMicroservice"

	"fmt"
	"log"
	"os"
)

func main() {
	// Ваши старые строчки, не трогаем
	err := db.Init()
	if !CheckDBConnection(err) {
		return
	}

	if err := db.MigrateBooks(); err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}
	if err := db.MigrateUser(); err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}
	if err := db.MigrateOTP(); err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}
	if err := db.MigrateEmailConfirmations(); err != nil {
		log.Fatalf("Ошибка миграции EmailConfirmations: %v", err)
	}

	// Если у вас есть миграция для корзины/транзакций
	// db.DB.AutoMigrate(&db.Cart{}, &db.CartItem{}, &db.Transaction{})

	db.FetchAndSaveBooks()

	// ➜ Новая строка: запускаем микросервис на :8081 в отдельной горутине
	go paymentMicroservice.StartMicroservice()

	// ➜ Ваш основной сервер (на :8080)
	sc.ConnectToServer()
}

func CheckDBConnection(err error) bool {
	if err != nil {
		fmt.Println("Error:", err.Error())
		fmt.Println("Do you wish to continue server initialization? (y/n)")
		var answer string
		_, err = fmt.Scan(&answer)
		if err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		}
		if answer == "y" {
			return true
		} else if answer == "n" {
			fmt.Println("Exiting...")
			return false
		}
	}
	return true
}
