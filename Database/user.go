package Database

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"regexp"
	"time"
)

// Структура пользователя в базе данных
type User struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `json:"name"`
	Email       string `json:"email" gorm:"unique"`
	Password    string `json:"-"`            // Храним пароль в зашифрованном виде
	IsConfirmed bool   `json:"is_confirmed"` // Подтвержден ли email
	ConfirmCode string `json:"-"`            // Код подтверждения (не отправляется клиенту)
	Role        string `json:"role"`         // user или admin
}

// Генерация случайного 6-значного кода подтверждения
func GenerateConfirmCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000)) // 6-значный код
}

// Хеширование пароля перед сохранением
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// Проверка пароля
func CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// Создание пользователя (с проверкой email и хешированием пароля)
func CreateUser(user User) error {
	if IsUserExistEmail(user.Email) {
		return errors.New("User already exists")
	}
	if !IsValidEmail(user.Email) {
		return errors.New("Invalid email format")
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}

	// Генерация кода подтверждения
	confirmCode := GenerateConfirmCode()

	user.Password = hashedPassword
	user.ConfirmCode = confirmCode
	user.IsConfirmed = false
	user.Role = "user" // По умолчанию обычный пользователь

	if err := DB.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

// Чтение пользователя по email
func (u *User) ReadUserEmail(email string) error {
	var user User
	log.Println("Reading user", email)

	if !IsUserExistEmail(email) {
		return errors.New("User not found")
	}

	result := DB.Select("name, email, role, is_confirmed").Where("email = ?", email).First(&user)
	if result.Error != nil {
		return result.Error
	}

	u.CopyUser(user)
	return nil
}

// Чтение всех пользователей
func ReadUserAll() ([]User, error) {
	var users []User
	log.Println("Reading all users")

	result := DB.Select("name, email, role, is_confirmed").Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

// Чтение пользователя по имени
func (u *User) ReadUserName(name string) error {
	var user User
	log.Println("Reading user", name)

	if !IsUserExistName(name) {
		return errors.New("User not found")
	}

	result := DB.Select("name, email, role, is_confirmed").Where("name = ?", name).First(&user)
	if result.Error != nil {
		return result.Error
	}

	u.CopyUser(user)
	return nil
}

// Копирование данных пользователя
func (u *User) CopyUser(user User) {
	u.ID = user.ID
	u.Email = user.Email
	u.Name = user.Name
	u.IsConfirmed = user.IsConfirmed
	u.Role = user.Role
}

// Обновление данных пользователя
func UpdateUser(email, newName string) error {
	var user User
	log.Println("Updating user:", email)

	result := DB.First(&user, "email = ?", email)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.New("User not found")
		}
		return result.Error
	}

	user.Name = newName
	if err := DB.Save(&user).Error; err != nil {
		return err
	}

	log.Println("User updated successfully:", user)
	return nil
}

// Удаление пользователя
func DeleteUser(email string) error {
	var user User
	log.Println("Deleting user:", email)

	result := DB.First(&user, "email = ?", email)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.New("User not found")
		}
		return result.Error
	}

	if err := DB.Delete(&user).Error; err != nil {
		return err
	}

	log.Println("User deleted successfully:", user)
	return nil
}

// Проверка валидности email (теперь через регулярные выражения)
func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// Проверка, существует ли пользователь с данным именем
func IsUserExistName(name string) bool {
	var count int64
	DB.Model(&User{}).Where("name = ?", name).Count(&count)
	return count > 0
}

// Проверка, существует ли пользователь с данным email
func IsUserExistEmail(email string) bool {
	var count int64
	DB.Model(&User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

// Форматированный вывод пользователя
func (u User) Stringer() string {
	return fmt.Sprintf("Name: %s, Email: %s, Role: %s, Confirmed: %v", u.Name, u.Email, u.Role, u.IsConfirmed)
}
