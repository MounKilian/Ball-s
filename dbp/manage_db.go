package dbp

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	var err error
	DB, err = gorm.Open(sqlite.Open("balls.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
}

func RegisterUser(username, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}
	err = DB.Create(&user).Error
	if err != nil {
		return err
	}

	return nil
}


func AuthenticateUser(username, password string) (bool, error) {
	var user User
	err := DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return false, nil // Mot de passe incorrect
	}

	return true, nil // Authentification réussie
}
