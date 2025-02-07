package paymentMicroservice

import (
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
)

func GeneratePDFReceipt(mtx MicroTransaction, items []struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Fiscal Receipt")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Transaction #: %d", mtx.TransactionID))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Date: %s", time.Now().Format("2006-01-02 15:04:05")))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Customer: %s", mtx.CustomerName))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Status: %s", mtx.Status))
	pdf.Ln(8)

	pdf.Ln(8)
	pdf.Cell(40, 10, "Items:")
	pdf.Ln(8)

	var total float64
	for _, it := range items {
		line := fmt.Sprintf("%s (ID=%s) — %.2f", it.Name, it.ID, it.Price)
		pdf.Cell(40, 10, line)
		pdf.Ln(6)
		total += it.Price
	}
	pdf.Ln(6)
	pdf.Cell(40, 10, fmt.Sprintf("Total: %.2f", total))

	filename := fmt.Sprintf("receipt_%d.pdf", mtx.TransactionID)
	if err := pdf.OutputFileAndClose(filename); err != nil {
		return "", err
	}
	return filename, nil
}
