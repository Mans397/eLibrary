package Database

import (
	"fmt"
	"github.com/Mans397/eLibrary/Api"
	"gorm.io/gorm"
	"log"
	"time"
)

type Book struct {
	gorm.Model
	Title       string `gorm:"unique;not null"`
	Description string
	Price       string
	Attributes  string
	Date        string
	ImageURL    string
}

func MigrateBooks() error {

	if err := DB.AutoMigrate(&Book{}); err != nil {
		log.Println("Ошибка при миграции:", err)
		return fmt.Errorf("ошибка миграции для Book: %v", err)
	}
	return nil
}

type EmailConfirmation struct {
	gorm.Model
	UserID uint      `gorm:"not null"`        // Связь с пользователем
	Token  string    `gorm:"unique;not null"` // Токен для подтверждения
	Expiry time.Time // Срок действия токена
}

func MigrateEmailConfirmations() error {
	if err := DB.AutoMigrate(&EmailConfirmation{}); err != nil {
		log.Println("Ошибка при миграции EmailConfirmation:", err)
		return fmt.Errorf("ошибка миграции для EmailConfirmation: %v", err)
	}
	return nil
}
func SaveBooks(books []Api.Book) error {
	for _, apiBook := range books {
		book := Book{
			Title:       apiBook.Title,
			Description: apiBook.Description,
			Price:       apiBook.Price,
			Attributes:  apiBook.Attributes,
			Date:        apiBook.Date,
			ImageURL:    apiBook.ImageURL,
		}

		if err := DB.Where("title = ?", book.Title).FirstOrCreate(&book).Error; err != nil {
			log.Printf("Ошибка сохранения книги '%s': %v\n", book.Title, err)
		} else {
			log.Printf("Книга '%s' успешно сохранена.\n", book.Title)
		}
	}
	return nil
}

func FetchAndSaveBooks() {
	fmt.Println("start")
	books, err := Api.FetchBooks()
	fmt.Println("start2")
	if err != nil {
		log.Fatalf("Ошибка при получении данных из API: %v", err)
	}
	fmt.Println("start3")
	for i, apiBook := range books {
		fmt.Println(i)
		var count int64
		if err := DB.Model(&Book{}).Where("title = ?", apiBook.Title).Count(&count).Error; err != nil {
			log.Fatalf("Ошибка при проверке наличия книги в базе: %v", err)
		}

		if count == 0 {

			book := Book{
				Title:       apiBook.Title,
				Description: apiBook.Description,
				Price:       apiBook.Price,
				Attributes:  apiBook.Attributes,
				Date:        apiBook.Date,
				ImageURL:    apiBook.ImageURL,
			}

			if err := DB.Where("title = ?", book.Title).FirstOrCreate(&book).Error; err != nil {
				log.Printf("Ошибка сохранения книги '%s': %v\n", book.Title, err)
			} else {
				log.Printf("Книга '%s' успешно сохранена.\n", book.Title)
			}
		} else {
			log.Printf("Книга '%s' уже существует в базе.\n", apiBook.Title)
		}
	}
}
