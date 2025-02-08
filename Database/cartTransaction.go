package Database

import (
	"gorm.io/gorm"
)

// Transaction — таблица для платежных транзакций.
type Transaction struct {
	gorm.Model
	CartID uint   // Или любая другая логика, если есть "корзина"
	Status string // "pending", "paid", "declined" и т.п.
}
type Cart struct {
	gorm.Model
	UserID uint
	Status string     // "open", "paid", "cancelled"
	Items  []CartItem `gorm:"foreignKey:CartID"`
}

type CartItem struct {
	gorm.Model
	CartID      uint
	ProductID   string
	ProductName string
	Price       float64
	Quantity    int
}

// MigrateCartAndTransaction — метод для миграции таблиц корзины/транзакций (опционально).
func MigrateCartAndTransaction() error {
	// Здесь же мигрируем Cart, CartItem и Transaction
	if err := DB.AutoMigrate(&Cart{}, &CartItem{}, &Transaction{}); err != nil {
		return err
	}
	return nil
}
