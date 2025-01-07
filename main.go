package main

import (
	"fmt"
	_ "github.com/Mans397/eLibrary/Api"
	db "github.com/Mans397/eLibrary/Database"
	sc "github.com/Mans397/eLibrary/serverConnection"
	"log"
	"os"
)

func main() {
	err := db.Init()
	if !CheckDBConnection(err) {
		return
	}

	// Миграция для таблицы Book
	if err := db.MigrateBooks(); err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}

	// Загрузка и сохранение книг
	db.FetchAndSaveBooks()
	sc.ConnectToServer()
}

func CheckDBConnection(err error) bool {
	if err != nil {
		fmt.Println("Error: " + err.Error())
		fmt.Println("Do you still wish to initialise the server(y/n)?")
		var answer string
		_, err = fmt.Scan(&answer)
		if err != nil {
			fmt.Println("Error:" + err.Error())
			os.Exit(1)
		}
		if answer == "y" {
			return true
		} else if answer == "n" {
			fmt.Println("Exiting...")
			return false
		} else {
			fmt.Println("Wrong answer")
			return false
		}
	}
	return true
}
