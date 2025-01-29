package Database

import (
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

const (
	host     = "localhost"
	port     = 5433
	user     = "postgres"
	password = "rootroot"
	dbname   = "eLibrary"
)

var DB *gorm.DB

func Init() error {
	//host := os.Getenv("DB_HOST")
	//port := os.Getenv("DB_PORT")
	//user := os.Getenv("DB_USER")
	//password := os.Getenv("DB_PASSWORD")
	//dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, "disable")
	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("!!!CHECK IF DOCKER CONTAINER IS STARTED!!!")
		time.Sleep(100 * time.Millisecond)
		return errors.New("Failed to connect to database")
	}

	fmt.Println("Successfully connected!")
	return nil
}
