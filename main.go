package main

import (
	"balls/dbp"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// db := dbp.DB
	// user1 := dbp.User{
	//     Username:    "user1",
	//     DateOfBirth: time.Now(),
	// }
	// user2 := dbp.User{
	//     Username:    "user2",
	//     DateOfBirth: time.Now(),
	// }
	// db.Create(&user1)
	// db.Create(&user2)

	r.GET("/", HomePage)
	r.GET("/users", getAllUsers)
	r.GET("/user", getUserByID)
	// r.POST("/login", loginUser)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// func loginUser(cgin.Context) {
//     type login struct{
//         NameOrMail string json:"name_or_mail"
//         Password   string json:"password"
//     }
//     loginData := login{}
//     c.Bind(&loginData)
//     db := dbp.DB
//     user := &dbp.User{}
//     db.Where(&dbp.User{Username: loginData.NameOrMail}).Or(&dbp.User{Email: loginData.NameOrMail}).Where(&dbp.User{Password: user.Password}).First(user)
// }

func HomePage(c *gin.Context) {

}

func getUserByID(c *gin.Context) {
	db := dbp.DB
	id, err := c.GetQuery("id")
	if !err {
		fmt.Fprintln(c.Writer, "Id not provided in get request")
		return
	}
	user := &dbp.User{}
	tx := db.First(user, id)
	if tx.RowsAffected > 0 {
		result, _ := json.Marshal(user)
		fmt.Fprintln(c.Writer, string(result))
	} else {
		fmt.Fprintln(c.Writer, "User not found")
	}
}

func getAllUsers(c *gin.Context) {
	db := dbp.DB
	users := []dbp.User{}
	db.Find(&users, &dbp.User{})
	result, _ := json.Marshal(&users)
	fmt.Fprintln(c.Writer, string(result))
}
