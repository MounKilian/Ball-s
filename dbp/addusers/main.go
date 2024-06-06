package main

import (
	"balls/dbp"
	"encoding/json"
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	last := &dbp.User{}
	tx := dbp.DB.Last(last)
	if tx.RowsAffected > 0 {
		fmt.Println("last ID :", last.ID)
	} else {
		last.ID = 0
	}
	for i := 0; i < 5000; i++ {

		user := dbp.User{
			Username: "Test : " + fmt.Sprint(last.ID+1),
			Password: "test",
			Email:    "test" + fmt.Sprint(last.ID+1) + "@test.com",
			Sport:    "Football",
		}

		dbp.DB.Create(&user)
		userID := user.ID

		fmt.Println("userID :", userID)
		// CloneDb("test")
	}
}

func Marshal(v any) string {
	re, err := json.Marshal(v)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(re)
}

func CloneDb(dbName string) {
	db1 := dbp.DB
	users := []dbp.User{}
	db1.Find(&users)
	stats := []dbp.Stat{}
	db1.Find(&stats)
	fmt.Println(Marshal(users[0]), Marshal(stats[0]))
	db2, err := gorm.Open(sqlite.Open(dbName+".db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if true && len(users) > 0 && len(stats) > 0 {
		db2.Create(users)
		db2.Create(stats)
	}
}
