package main

import (
	"balls/dbp"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	last := &dbp.User{}
	dbp.DB.Last(&last)

	for i := last.ID; i < 5000; i++ {
		time.Sleep(1 * time.Second)

		err := dbp.RegisterUser("Test : "+fmt.Sprint(0), "test"+fmt.Sprint(i)+"@test.com", "test")
		fmt.Println("ID :", i, "error :", err)
	}
	// CloneDb("test")
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
