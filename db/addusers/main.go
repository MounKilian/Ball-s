package main

import (
	"balls/db"
	"fmt"
)

func main() {
	last := &db.User{}
	last.GetLastUser()
	fmt.Println("last ID :", last.ID)

	user := db.User{
		Username: "Test : " + fmt.Sprint(last.ID+1),
		Password: "test1",
		Email:    "test" + fmt.Sprint(last.ID+1) + "@test.com",
	}
	userID := user.CreateUser()
	fmt.Println("userID :", userID)
	//_ = DeleteUser(3)
	//fmt.Println("deleted :", done)
	//_ = DB.First(&User{}, 3)
	//rjson, _ := json.MarshalIndent(struct {
	//	Error        string
	//	RowsAffected int64
	//}{
	//	Error:        fmt.Sprint(tx.Error),
	//	RowsAffected: tx.RowsAffected,
	//}, "", "\t")
	//fmt.Println(string(rjson))
	//DB.Exec(`DELETE FROM users`)
}
