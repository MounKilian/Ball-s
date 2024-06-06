package main

import (
	"balls/dbp"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = db.AutoMigrate(&dbp.User{}, &dbp.Strike{}, &dbp.Miss{}, &dbp.Stat{}, &dbp.Match{})
	if err != nil {
		fmt.Println(err)
	}
}
