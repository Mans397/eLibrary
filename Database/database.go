package Database

import (
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

const (
	host     = "localhost"
	port     = 5433
	user     = "postgres"
	password = "rootroot"
	dbname   = "eLibrary"
)

var DB *gorm.DB

func Init() {

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, "disable")
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected!")

}
