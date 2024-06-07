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

	// CloneDb("balls", "test")
	// AddUsers()
	// AddImages()
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

func AddImages() {
	images := []string{
		"/static/img/_08095b50-f04c-4ff5-9750-26debab1ba96.jpg",
		"/static/img/_a63562ee-b105-4eb8-bb4f-1b877f274e36.jpg",
		"/static/img/_afa838f5-cb23-492e-abdf-ee8778076d94.jpg",
		"/static/img/angle-droit (1).png",
		"/static/img/angle-droit.png",
		"/static/img/_b4a4f3a2-ab2b-4a31-9b62-9a601e8375e5.jpg",
		"/static/img/_e1bfede1-8d6f-4344-9469-0e0238991e26.jpg",
		"/static/img/image.png",
		// "/static/img/IMG_8302.jpg",
		"/static/img/IMG_9082.jpg",
		"/static/img/IMG_9083.jpg",
		"/static/img/lapin.jpg",
		"/static/img/IMG_1362.jpg",
	}
	_ = images
	db := dbp.DB
	users := []dbp.User{}
	db.Not(4, 5).Find(&users, map[string]any{"image": ""})
	fmt.Println("len :", len(users))
	for _, u := range users {
		if u.ID == 4 || u.ID == 5 {
			fmt.Println("error: killian found")
			fmt.Println(u)
		}
		// fmt.Println(u.ID)
		image := images[rand.Intn(12)]
		u.Image = image
		db.Save(&u)
	}
}

func CloneDb(dbName1, dbName2 string) {
	db1, err := gorm.Open(sqlite.Open(dbName1+".db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	users := []dbp.User{}
	db1.Find(&users)
	fmt.Println("len :", len(users))
	stats := []dbp.Stat{}
	db1.Find(&stats)
	// fmt.Println(Marshal(users[0]), Marshal(stats[0]))
	db2, err := gorm.Open(sqlite.Open(dbName2+".db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if len(users) > 0 && len(stats) > 0 {
		db2.Save(stats)
		// db2.Create(users)
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

	for i := last.ID; i < 1500; i++ {
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
			cityname = "Paris"
		} else {
			cityname = "Lyon"
		}

		if genderand == 0 {
			gender = "Homme"
		} else {
			gender = "Femme"
		}

		if genderprefrand == 0 {
			genderpref = "Homme"
		} else {
			genderpref = "Femme"
		}
		user := dbp.User{
			Username:      "User" + fmt.Sprint(i),
			Email:         "user" + fmt.Sprint(i) + "@test.com",
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

	db.Find(&swipes, &dbp.Miss{})
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

		swipe := dbp.Miss{
			UserAID: int(users[f].ID),
			UserA:   users[f],
			UserBID: int(potential[0].ID),
			UserB:   potential[0],
		}
		db.Create(&swipe)

	}
}

func addMatches() {
	db := dbp.DB
	rand.Seed(time.Now().UnixNano())

	existingMatches := make(map[[2]int]bool)

	for i := 0; i < 2000; {
		userA := rand.Intn(1511) + 1
		userB := rand.Intn(1511) + 1
		if userA == userB {
			continue
		}

		if existingMatches[[2]int{userA, userB}] || existingMatches[[2]int{userB, userA}] {
			continue
		}
		roomName := fmt.Sprint(userA) + "-" + fmt.Sprint(userB)
		match1 := dbp.Match{
			UserAID:  userA,
			UserBID:  userB,
			RoomName: roomName,
		}

		match2 := dbp.Match{
			UserAID:  userB,
			UserBID:  userA,
			RoomName: roomName,
		}

		if err := db.Create(&match1).Error; err != nil {
			log.Printf("Erreur lors de l'insertion de la correspondance : %v", err)
			return
		}

		if err := db.Create(&match2).Error; err != nil {
			log.Printf("Erreur lors de l'insertion de la correspondance inverse : %v", err)
			return
		}

		existingMatches[[2]int{userA, userB}] = true
		existingMatches[[2]int{userB, userA}] = true
		i++
	}
}
