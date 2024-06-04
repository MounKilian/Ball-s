package dbp

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", HomePage)
	r.GET("/users", getAllUsers)
	r.GET("/user", getUserByID)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}

func HomePage(c *gin.Context) {

}

func getUserByID(c *gin.Context) {
	db := DB
	id, err := c.GetQuery("id")
	if !err {
		fmt.Fprintln(c.Writer, "Id not provided in get request")
		return
	}
	user := &User{}
	tx := db.First(user, id)
	if tx.RowsAffected > 0 {
		result, _ := json.Marshal(user)
		fmt.Fprintln(c.Writer, string(result))
	} else {
		fmt.Fprintln(c.Writer, "User not found")
	}
}

func getAllUsers(c *gin.Context) {
	db := DB
	users := []User{}
	db.Find(&users, &User{})
	result, _ := json.Marshal(&users)
	fmt.Fprintln(c.Writer, string(result))
}
