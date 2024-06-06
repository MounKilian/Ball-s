package dbp

import (
	"errors"
	"log"
	"strconv"

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
		return errors.New("erreur lors du hachage du mot de passe")
	}

	user := User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}
	err = DB.Create(&user).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		log.Println("Erreur de contrainte de vérification")
		return errors.New("un utilisateur avec le même nom d'utilisateur ou la même adresse e-mail existe déjà")
	} else if err != nil {
		return errors.New("erreur lors de l'enregistrement de l'utilisateur dans la base de données")
	}

	return nil
}

func AuthenticateUser(username, password string) (string, bool, error) {
	var user User
	err := DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return "", false, errors.New("Nom d'utilisateur incorrect")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", false, errors.New("Mot de passe incorrect")
	}

	return strconv.Itoa(int(user.ID)), true, nil
}
