package emailSender

import (
	db "github.com/Mans397/eLibrary/Database"
	"gopkg.in/gomail.v2"
	"log"
)

func SendEmailAll(text *string) error {
	var users []db.User

	result := db.DB.Find(&users)
	if result.Error != nil {
		log.Printf("failed to fetch users: %v", result.Error)
		return result.Error
	}

	for _, user := range users {
		err := SendEmail(user.Email, text)
		if err != nil {
			log.Printf("failed to send email: %v", err)
		}
	}

	return nil
}

func SendEmail(email string, text *string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "elibrarysender@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Notification")
	m.SetBody("text/plain", *text)

	d := gomail.NewDialer("smtp.gmail.com", 587, "elibrarysender@gmail.com", "ocxwblzcockfwcud")

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email: %v, User: %v\n", err, email)
		return err
	}
	log.Println("Email sent successfully to user:", email)
	return nil
}
