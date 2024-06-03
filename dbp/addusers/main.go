package main

import (
	"balls/dbp"
	"fmt"
)

func main() {
	last := &dbp.User{}
	tx := dbp.DB.Last(last)
	if tx.RowsAffected > 0 {
		fmt.Println("last ID :", last.ID)
	} else {
		last.ID = 0
	}

	user := dbp.User{
		Username: "Test : " + fmt.Sprint(last.ID+1),
		Password: "test1",
		Email:    "test" + fmt.Sprint(last.ID+1) + "@test.com",
	}

	dbp.DB.Create(&user)
	userID := user.ID

	fmt.Println("userID :", userID)
}
