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
	// dbp.DB.Last(&last)
	// // i := last.ID
	// for i := last.ID + 1; i <= 5000; i++ {
	// 	// time.Sleep(1 * time.Second)

	// 	dbp.RegisterUser("Test : "+fmt.Sprint(i), "test"+fmt.Sprint(i)+"@test.com", "test")
	// 	fmt.Println("ID :", i)
	// }

	// CloneDb("test")
	AddUsers()
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
		// db2.Create(stats)
	}
}

func AddSport() {
	db := dbp.DB
	for _, v := range []string{"Football", "Basketball", "Tennis", "Baseball", "Surf", "Volley", "Pingpong", "Golf", "Natation", "Rugby", "Bowling", "Handball", "Escalade", "Cyclisme", "Sauts", "Plongée", "Acrobranche", "Tyroliènne", "Course", "Musculation", "Randonnée", "Paddle", "Acrobranche", "Ski", "Boxe", "MMA", "Kapoera", "Pétanque", "Gymnastique", "Danse", "Karting", "Paintball", "Judo", "Karaté", "Escrime", "Ultimate", "LaserGame", "Je ne fait pas que du sport"} {

		db.Create(&dbp.Stat{Name: v})
	}
}

func AddUsers() {
	last := &dbp.User{}
	tx := dbp.DB.Last(last)
	if tx.RowsAffected > 0 {
		fmt.Println("last ID :", last.ID)
	} else {
		last.ID = 0
	}

	for i := last.ID; i < 3000; i++ {
		db := dbp.DB

		max := int64(0)
		sport := dbp.Stat{}

		db.Table("stats").Count(&max)
		db.First(&sport, rand.Intn(int((max-1)/3))+1)

		cityrand := rand.Intn(2)
		genderand := rand.Intn(2)
		genderprefrand := rand.Intn(2)
		var cityname string
		var gender string
		var genderpref string
		if cityrand == 0 {
			cityname = "paris"
		} else {
			cityname = "Lyon"
		}

		if genderand == 0 {
			gender = "men"
		} else {
			gender = "women"
		}

		if genderprefrand == 0 {
			genderpref = "men"
		} else {
			genderpref = "women"
		}
		user := dbp.User{
			Username:      "User",
			DateOfBirth:   time.Now().Add(time.Duration(rand.Intn(10000)) * time.Hour * 24 * -1),
			Sport:         sport.Name,
			Gender:        gender,
			DesiredGender: genderpref,
			City:          cityname,
		}

		db.Create(&user)

	}
}

func AddSwipe() {

	db := dbp.DB
	users := []dbp.User{}
	swipes := []int{}

	db.Find(&swipes, &dbp.Swipe{})
	db.Find(&users)

	for f := 0; f < len(users); f++ {

		var potential []dbp.User
		potential = nil
		for i := 0; i < len(users); i++ {
			if users[f].City == users[i].City && users[i].ID != users[f].ID {
				if users[f].Sport == users[i].Sport && users[f].DesiredGender == users[i].Gender {
					potential = append(potential, users[i])
				}
			}
		}

		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(potential), func(i, j int) { potential[i], potential[j] = potential[j], potential[i] })
		// result, _ := json.Marshal(&users)
		// resultpot, _ := json.Marshal((&potential))

		swipe := dbp.Swipe{
			UserAID: int(users[f].ID),
			UserA:   users[f],
			UserBID: int(potential[0].ID),
			UserB:   potential[0],
		}
		db.Create(&swipe)

	}
}
