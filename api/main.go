package main

import (
	api "api/code"
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {

	db = api.Connect()
	defer db.Close()
	// var err error
	// db, err = sql.Open("sqlite3", "./base.db")
	// if err != nil {
	// 	log.Println("nil sql")
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	// err = db.Ping()
	// if err != nil {
	// 	log.Default().Println("erreur ping db")
	// 	log.Fatal(err)
	// }

	// _, err = db.Exec(`
	// 	CREATE TABLE IF NOT EXISTS users (
	// 		id INTEGER PRIMARY KEY AUTOINCREMENT,
	// 		name TEXT NOT NULL,
	// 		email TEXT NOT NULL UNIQUE,
	// 		password TEXT NOT NULL
	// 	)
	// `)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// _, err = db.Exec(`
	// 	CREATE TABLE IF NOT EXISTS posts (
	// 		id INTEGER PRIMARY KEY AUTOINCREMENT,
	// 		user_id INTEGER NOT NULL,
	// 		content TEXT NOT NULL,
	// 		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	// 		FOREIGN KEY(user_id) REFERENCES users(id)
	// 	)
	// `)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	//Configurer les routes
	api.SetupRoutes()

	log.Println("API démarrée sur le port 8181")
	/* trunk-ignore(golangci-lint/errcheck) */
	http.ListenAndServe(":8181", nil)
}
