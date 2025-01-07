package Api

// Book представляет структуру книги из API
type Book struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       string // Цена, добавим вручную
	Attributes  string // Характеристики, добавим вручную
	Date        string `json:"publish_date"`
	ImageURL    string `json:"cover_url"`
}
