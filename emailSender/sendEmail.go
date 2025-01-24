package emailSender

import (
	db "github.com/Mans397/eLibrary/Database"
	"gopkg.in/gomail.v2"
	"log"
	"sync"
)

func SendEmailAll(text *string, imagePath string) error {
	var users []db.User
	wg := new(sync.WaitGroup)

	result := db.DB.Find(&users)
	if result.Error != nil {
		log.Printf("failed to fetch users: %v", result.Error)
		return result.Error
	}
	log.Printf("Starting sending email for %d users", len(users))
	for _, user := range users {
		wg.Add(1)
		go SendEmail(user.Email, text, imagePath, wg)

	}
	wg.Wait()
	log.Printf("Finished sending emails")

	return nil
}

func SendEmail(email string, text *string, imagePath string, wg *sync.WaitGroup) {
	defer wg.Done()

	m := gomail.NewMessage()
	m.SetHeader("From", "elibrarysender@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Notification")
	m.SetBody("text/plain", *text)

	if imagePath != "" {
		m.Attach(imagePath)
	}

	d := gomail.NewDialer("smtp.gmail.com", 587, "elibrarysender@gmail.com", "ocxwblzcockfwcud")

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Ошибка отправки письма для пользователя %s: %v\n", email, err)
		return
	}

	log.Printf("Письмо успешно отправлено пользователю %s\n", email)
}

func SendEmailLogin(email string, text *string, imagePath string) {

	m := gomail.NewMessage()
	m.SetHeader("From", "elibrarysender@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Notification")
	m.SetBody("text/plain", *text)

	if imagePath != "" {
		m.Attach(imagePath)
	}

	d := gomail.NewDialer("smtp.gmail.com", 587, "elibrarysender@gmail.com", "ocxwblzcockfwcud")

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Ошибка отправки письма для пользователя %s: %v\n", email, err)
		return
	}

	log.Printf("Письмо успешно отправлено пользователю %s\n", email)
}
