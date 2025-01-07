package Database

import (
	"fmt"
	"github.com/Mans397/eLibrary/Api"
	"gorm.io/gorm"
	"log"
)

// Book структура для базы данных
type Book struct {
	gorm.Model
	Title       string `gorm:"unique;not null"` // Название книги должно быть уникальным
	Description string
	Price       string
	Attributes  string
	Date        string
	ImageURL    string
}

// MigrateBooks создаёт таблицу для книг в базе данных
func MigrateBooks() error {
	// Автоматическая миграция для создания таблицы Book
	if err := DB.AutoMigrate(&Book{}); err != nil {
		log.Println("Ошибка при миграции:", err)
		return fmt.Errorf("ошибка миграции для Book: %v", err)
	}
	return nil
}

// SaveBooks сохраняет список книг в базу данных
func SaveBooks(books []Api.Book) error {
	for _, apiBook := range books {
		// Преобразуем Api.Book в Database.Book
		book := Book{
			Title:       apiBook.Title,
			Description: apiBook.Description,
			Price:       apiBook.Price,
			Attributes:  apiBook.Attributes,
			Date:        apiBook.Date,
			ImageURL:    apiBook.ImageURL,
		}

		// Сохраняем книгу, пропуская дубликаты
		if err := DB.Where("title = ?", book.Title).FirstOrCreate(&book).Error; err != nil {
			log.Printf("Ошибка сохранения книги '%s': %v\n", book.Title, err)
		} else {
			log.Printf("Книга '%s' успешно сохранена.\n", book.Title)
		}
	}
	return nil
}

func FetchAndSaveBooks() {
	// Получаем книги из API
	books, err := Api.FetchBooks()
	if err != nil {
		log.Fatalf("Ошибка при получении данных из API: %v", err)
	}

	// Проходим по каждой книге и сохраняем, если её ещё нет в базе
	for _, apiBook := range books {
		// Проверяем, существует ли книга с таким же названием в базе данных
		var count int64
		if err := DB.Model(&Book{}).Where("title = ?", apiBook.Title).Count(&count).Error; err != nil {
			log.Fatalf("Ошибка при проверке наличия книги в базе: %v", err)
		}

		// Если книга уже есть в базе данных, пропускаем её
		if count == 0 {
			// Преобразуем Api.Book в Database.Book
			book := Book{
				Title:       apiBook.Title,
				Description: apiBook.Description,
				Price:       apiBook.Price,
				Attributes:  apiBook.Attributes,
				Date:        apiBook.Date,
				ImageURL:    apiBook.ImageURL,
			}

			// Сохраняем книгу, пропуская дубликаты
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
