package Api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

const OpenLibraryAPI = "https://openlibrary.org/search.json?q=programming"

// Пример списка случайных описаний
var descriptions = []string{
	"A captivating story that engages the reader from start to finish.",
	"An insightful exploration of a fascinating subject.",
	"A timeless tale that resonates with readers of all ages.",
	"A must-read for those looking to expand their knowledge.",
	"An enriching experience that broadens the mind and spirit.",
}

// Генерация случайной цены
func generateRandomPrice() string {
	prices := []string{"10.99 USD", "12.50 USD", "8.99 USD", "15.75 USD", "9.99 USD"}
	rand.Seed(time.Now().UnixNano())
	return prices[rand.Intn(len(prices))]
}

// Генерация случайного описания
func generateRandomDescription() string {
	rand.Seed(time.Now().UnixNano())
	return descriptions[rand.Intn(len(descriptions))]
}

// FetchBooks получает данные из API
func FetchBooks() ([]Book, error) {
	resp, err := http.Get(OpenLibraryAPI)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка API: %s", resp.Status)
	}

	// Парсим JSON-ответ
	var data struct {
		Docs []struct {
			Title       string   `json:"title"`
			PublishDate []string `json:"publish_date"`
			CoverID     int      `json:"cover_i"`
		} `json:"docs"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %v", err)
	}

	// Преобразуем в []Book
	books := []Book{}
	for _, doc := range data.Docs {
		var date string
		if len(doc.PublishDate) > 0 {
			date = doc.PublishDate[0] // Берем первую дату, если она есть
		} else {
			date = "Unknown" // Если данных нет, подставляем "Unknown"
		}

		book := Book{
			Title:       doc.Title,
			Date:        date,
			ImageURL:    fmt.Sprintf("https://covers.openlibrary.org/b/id/%d-L.jpg", doc.CoverID),
			Price:       generateRandomPrice(),
			Attributes:  "Softcover, 200 pages",
			Description: generateRandomDescription(),
		}
		books = append(books, book)
		if len(books) >= 30 { // Ограничиваем до 30 записей
			break
		}
	}
	return books, nil
}
