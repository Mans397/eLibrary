package main

import (
	"fmt"
	db "github.com/Mans397/eLibrary/Database"
	sc "github.com/Mans397/eLibrary/serverConnection"
	"log"
	"os"

	// ➜ Добавляем импорт вашего "микросервиса"
	"github.com/Mans397/eLibrary/paymentMicroservice"
)

func main() {

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

	if err := db.MigrateCartAndTransaction(); err != nil {
		log.Fatalf("Ошибка миграции Transaction: %v", err)
	}

	db.FetchAndSaveBooks()

	// -- СТАРЫЕ СТРОКИ НЕ ТРОГАЕМ, ТОЛЬКО ДОБАВЛЯЕМ НОВУЮ --

	// ➜ Запуск "второго" сервера (микросервиса) в отдельной горутине на порту :8081
	go paymentMicroservice.StartMicroservice()

	// ➜ Основной сервер (как и прежде) стартует на :8080
	sc.ConnectToServer()
}

// Старую функцию не трогаем
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
