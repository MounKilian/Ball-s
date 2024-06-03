package dbp

import (
	"errors"
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
	var existingUser User
	err := DB.Where("username = ? OR email = ?", username, email).First(&existingUser).Error
	if err == nil {
		return errors.New("un utilisateur avec le même nom d'utilisateur ou la même adresse e-mail existe déjà")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("erreur lors de la recherche de l'utilisateur dans la base de données")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("erreur lors du hachage du mot de passe")
	}

	user := User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}
	err = DB.Create(&user).Error
	if err != nil {
		return errors.New("erreur lors de l'enregistrement de l'utilisateur dans la base de données")
	}

	return nil
}

func AuthenticateUser(username, password string) (bool, error) {
	var user User
	err := DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return false, errors.New("nom d'utilisateur incorrect")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return false, errors.New("mot de passe incorrect")
	}

	return true, nil
}
