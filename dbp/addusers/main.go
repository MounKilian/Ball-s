package main

import (
	"balls/dbp"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// last := &dbp.User{}
	// tx := dbp.DB.Last(last)
	// if tx.RowsAffected > 0 {
	// 	fmt.Println("last ID :", last.ID)
	// } else {
	// 	last.ID = 0
	// }

	// user := dbp.User{
	// 	Username: "Test : " + fmt.Sprint(last.ID+1),
	// 	Password: "test1",
	// 	Email:    "test" + fmt.Sprint(last.ID+1) + "@test.com",
	// }

	// dbp.DB.Create(&user)
	// userID := user.ID

	// fmt.Println("userID :", userID)
	addMatches()
}

func Marshal(v any) string {
	re, err := json.Marshal(v)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(re)
}

func CloneDb() {
	db1 := dbp.DB
	users := []dbp.User{}
	db1.Find(&users)
	stats := []dbp.Stat{}
	db1.Find(&stats)
	fmt.Println(Marshal(users[0]), Marshal(stats[0]))
	db2, err := gorm.Open(sqlite.Open("balls2.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if true && len(users) > 0 && len(stats) > 0 {
		db2.Create(users)
		db2.Create(stats)
	}
}

func addMatches() {
	db := dbp.DB

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 10; i++ {
		match := dbp.Match{
			UserAID: rand.Intn(5) + 1,
			UserBID: rand.Intn(5) + 1,
		}
		for match.UserAID == match.UserBID {
			match.UserBID = rand.Intn(5) + 1
		}

		if err := db.Create(&match).Error; err != nil {
			log.Printf("Error inserting match: %v", err)
			return
		}
	}
}
