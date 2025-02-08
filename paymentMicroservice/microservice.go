package paymentMicroservice

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

// StartMicroservice – точка входа «микросервиса», которую мы вызовем из main.go
func StartMicroservice() {
	initMicroserviceDB()

	http.HandleFunc("/payment", handlePayment)

	log.Println("Microservice running on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Println("Error starting microservice:", err)
	}
}

// handlePayment обрабатывает POST-запрос с данными транзакции и корзины
func handlePayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		TransactionID uint `json:"transaction_id"`
		CartItems     []struct {
			ID    string  `json:"id"`
			Name  string  `json:"name"`
			Price float64 `json:"price"`
		} `json:"cartItems"`
		Customer struct {
			ID    uint   `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"customer"`
	}

	// Парсим JSON
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Создаём запись в таблице нашего микросервиса
	mtx := MicroTransaction{
		TransactionID: input.TransactionID,
		CustomerID:    input.Customer.ID,
		CustomerName:  input.Customer.Name,
		CustomerEmail: input.Customer.Email,
		Status:        "pending",
		CreatedAt:     time.Now().Format("2006-01-02 15:04:05"),
	}
	if err := microDB.Create(&mtx).Error; err != nil {
		http.Error(w, "Failed to create MicroTransaction", http.StatusInternalServerError)
		return
	}

	// Допустим, чётный ID = успех, нечётный = отказ (как пример)
	paymentSuccess := (mtx.TransactionID%2 == 0)

	if paymentSuccess {
		mtx.Status = "paid"
		microDB.Save(&mtx)

		// Генерируем PDF (см. pdfReceipt.go)
		pdfPath, err := GeneratePDFReceipt(mtx, input.CartItems)
		if err != nil {
			log.Println("Failed to generate PDF:", err)
		} else {
			// Отправляем чек на email (см. emailReceipt.go)
			if sendErr := SendReceiptByEmail(pdfPath, mtx.CustomerEmail); sendErr != nil {
				log.Println("Failed to send receipt:", sendErr)
			}
			// Удаляем временный файл
			os.Remove(pdfPath)
		}

		respondJSON(w, map[string]interface{}{
			"success": true,
		})
	} else {
		mtx.Status = "declined"
		microDB.Save(&mtx)

		respondJSON(w, map[string]interface{}{
			"success":       false,
			"error_message": "Payment declined by test logic",
		})
	}
}

func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
