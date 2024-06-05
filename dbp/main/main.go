package main

import (
	"balls/dbp"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("pages/*.html")

	router.Static("/static", "./static")
	router.GET("/users", getAllUsers)
	router.GET("/user", getUserByID)
	router.GET("/sports", getAllSports)

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	router.Use(cors.New(config))

	log.Fatal(router.Run(":8081"))
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

func getAllSports(c *gin.Context) {
	db := dbp.DB
	sports := []dbp.Stat{}
	db.Find(&sports, []dbp.Stat{})
	result, _ := json.Marshal(&sports)
	fmt.Fprintln(c.Writer, string(result))
}
