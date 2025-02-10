package paymentMicroservice

import (
	"log"

	gomail "gopkg.in/gomail.v2"
)

func SendReceiptByEmail(pdfPath string, toEmail string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "elibrarysender@gmail.com")
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Your Receipt")
	m.SetBody("text/plain", "Thank you for your purchase! See attached receipt.")

	if pdfPath != "" {
		m.Attach(pdfPath)
	}

	d := gomail.NewDialer("smtp.gmail.com", 587, "elibrarysender@gmail.com", "ocxwblzcockfwcud")

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send PDF email to %s: %v\n", toEmail, err)
		return err
	}
	log.Printf("PDF receipt emailed to %s\n", toEmail)
	return nil
}
