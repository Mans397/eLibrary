package serverConnection

import (
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// GenerateSingleBookPDFReceipt — генерация PDF для одной книги.
// bookTitle — название книги
// quantity — количество
// price — цена за единицу
// fullName — ФИО клиента
// cardNumber — для вывода "****3456" и т.д.
// transactionID — номер транзакции
func GenerateSingleBookPDFReceipt(
	companyName string,
	transactionID uint,
	bookTitle string,
	price float64,
	quantity int,
	fullName string,
	cardNumber string,
	email string,
) (string, error) {

	// Прячем часть номера карты (например, все кроме последних 4)
	maskedCard := "****"
	if len(cardNumber) > 4 {
		maskedCard = "****" + cardNumber[len(cardNumber)-4:]
	}

	total := price * float64(quantity)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Fiscal Receipt")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, fmt.Sprintf("Company/Project: %s", companyName))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Transaction ID: %d", transactionID))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Date/Time: %s", time.Now().Format("2006-01-02 15:04:05")))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Customer: %s", fullName))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Email: %s", email))
	pdf.Ln(8)

	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Item: %s", bookTitle))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Price per unit: %.2f", price))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Quantity: %d", quantity))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Total sum: %.2f", total))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Payment method (card masked): %s", maskedCard))

	filename := fmt.Sprintf("receipt_%d.pdf", transactionID)
	err := pdf.OutputFileAndClose(filename)
	if err != nil {
		return "", err
	}
	return filename, nil
}
