package main

import (
	"balls/dbp"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("balls.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = db.AutoMigrate(&dbp.User{}, &dbp.Strike{}, &dbp.Miss{}, &dbp.Stat{}, &dbp.Match{})
	if err != nil {
		fmt.Println(err)
	}
}

/* func main() {
	// Chemin vers la base de données SQLite
	dbPath := "./example.db"

	// Ouvrir la connexion à la base de données
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture de la base de données:", err)
		return
	}
	defer db.Close()

	// Création des tables
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS Users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username VARCHAR NOT NULL UNIQUE,
            password VARCHAR NOT NULL,
            email VARCHAR NOT NULL UNIQUE,
            created_at DATETIME,
            updated_at DATETIME,
            deleted_at DATETIME
        );
        CREATE TABLE IF NOT EXISTS Profiles (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER UNIQUE,
            first_name VARCHAR,
            last_name VARCHAR,
            email VARCHAR,
            phone VARCHAR,
            address TEXT,
            birthday DATE,
            created_at DATETIME,
            updated_at DATETIME,
            deleted_at DATETIME,
            FOREIGN KEY(user_id) REFERENCES Users(id)
        );
        CREATE TABLE IF NOT EXISTS Posts (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title VARCHAR,
            body TEXT,
            user_id INTEGER UNIQUE,
            status VARCHAR,
            created_at DATETIME,
            updated_at DATETIME,
            deleted_at DATETIME,
            FOREIGN KEY(user_id) REFERENCES Users(id)
        );
        CREATE TABLE IF NOT EXISTS Comments (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            body TEXT,
            user_id INTEGER,
            post_id INTEGER,
            created_at DATETIME,
            updated_at DATETIME,
            deleted_at DATETIME,
            FOREIGN KEY(user_id) REFERENCES Users(id),
            FOREIGN KEY(post_id) REFERENCES Posts(id)
        );
        CREATE TABLE IF NOT EXISTS Likes (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER,
            post_id INTEGER NULL,
            comment_id INTEGER NULL,
            created_at DATETIME,
            deleted_at DATETIME,
            UNIQUE(user_id, post_id, comment_id),
            FOREIGN KEY(user_id) REFERENCES Users(id),
            FOREIGN KEY(post_id) REFERENCES Posts(id),
            FOREIGN KEY(comment_id) REFERENCES Comments(id)
        );
    `)
	if err != nil {
		fmt.Println("Erreur lors de la création des tables:", err)
		return
	}

	fmt.Println("Base de données créée avec succès !")
}
*/
