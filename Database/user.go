package Database

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type User struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func CreateUser(user User) error {
	isExist := IsUserExist(user.Email)
	if isExist {
		return errors.New("User already exists")
	}
	isValid := IsValidEmail(user.Email)
	if !isValid {
		return errors.New("Invalid email")
	}

	DB.Create(&User{Name: user.Name, Email: user.Email})
	return nil

}

func ReadUser(email string, user *User) error {
	isExist := IsUserExist(email)
	if !isExist {
		return errors.New("User not found")
	}

	result := DB.Select("name, email").Where("email = ?", email).First(&user)
	fmt.Println(user)
	if result.Error != nil {
		return errors.New(result.Error.Error())
	}
	return nil
}

func IsValidEmail(email string) bool {
	chars := "@gmail.com"
	if strings.Contains(email, chars) {
		return true
	}
	return false
}

func IsUserExist(email string) bool {
	var user User
	result := DB.First(&user, "email = ?", email)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false
		} else {
			fmt.Println("Произошла ошибка:", result.Error)
		}
	} else {
		return true
	}
	return false
}

func (u User) Stringer() string {
	return fmt.Sprintf("Name: %s Email: %s", u.Name, u.Email)
}
