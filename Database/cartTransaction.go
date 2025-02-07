package Database

import (
	"gorm.io/gorm"
	"log"
)

// Transaction — таблица для платежных транзакций.
type Transaction struct {
	gorm.Model
	CartID uint   // Или любая другая логика, если есть "корзина"
	Status string // "pending", "paid", "declined" и т.п.
}

// MigrateCartAndTransaction — метод для миграции таблиц корзины/транзакций (опционально).
func MigrateCartAndTransaction() error {
	err := DB.AutoMigrate(&Transaction{})
	if err != nil {
		log.Println("Ошибка миграции Transaction:", err)
		return err
	}
	return nil
}
