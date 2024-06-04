package main

import (
	"balls/dbp"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// addUsers()
	// addSport()
	// http.HandleFunc("/", HomePage)
	// fs := http.FileServer(http.Dir("static"))
	// http.Handle("/static/", http.StripPrefix("/static/", fs))
	// http.ListenAndServe(":8080", nil)
	sort()
}

func addSport() {
	db := dbp.DB
	// stat := dbp.Stat{
	// 	ID: 1,
	// 	Name: "",
	// }
	// user2 := dbp.User{
	// 	Username:    "user2",
	// 	DateOfBirth: time.Now(),
	// }
	// db.Create(&user1)
	for _, v := range []string{"Football", "Basketball", "Tenis", "baseball", "Volley", "Pingpong", "Golf", "Natation", "Bowling", "Escalade", "Cyclisme", "Sauts", "Plongée", "Acrobranches", "tyroliènne", "Course", "Musculation"} {

		db.Create(&dbp.Stat{Name: v})
	}
	// db.Create(&stat)
}

func addUsers() {

	last := &dbp.User{}
	tx := dbp.DB.Last(last)
	if tx.RowsAffected > 0 {
		fmt.Println("last ID :", last.ID)
	} else {
		last.ID = 0
	}

	for i := 0; i < 50; i++ {
		db := dbp.DB
		sportid := rand.Intn(17)
		cityrand := rand.Intn(2)
		var cityname string
		if cityrand == 0 {
			cityname = "paris"
		} else {
			cityname = "Lyon"
		}
		user := dbp.User{
			Username:    "User",
			DateOfBirth: time.Now(),
			SportID:     sportid,
			City:        cityname,
		}

		db.Create(&user)

	}
}

func sort() {
	db := dbp.DB
	users := []dbp.User{}
	db.Find(&users, &dbp.User{})
	startUser := rand.Intn(50)
	var potential []dbp.User
	for i := 0; i < len(users); i++ {
		if users[startUser].City == users[i].City && i != startUser {
			if users[startUser].SportID == users[i].SportID {
				potential = append(potential, users[i])
			}
		}
	}
	// result, _ := json.Marshal(&users)
	// resultpot, _ := json.Marshal((&potential))
	fmt.Println(users[startUser].ID)
	fmt.Println("potential")
	for i := 0; i < len(potential); i++ {
		fmt.Println(potential[i].ID)
	}
}
