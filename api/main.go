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
	var err error
	db, err = sql.Open("sqlite3", "./base.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Exécuter les migrations
	api.RunMigrations(db)

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Configurer les routes en passant la base de données
	// api.SetupRoutes(db)

	log.Println("API démarrée sur le port 8181")
	/* trunk-ignore(golangci-lint/errcheck) */
	http.ListenAndServe(":8181", nil)
}
