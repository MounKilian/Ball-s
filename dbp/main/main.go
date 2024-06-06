package main

import (
	"balls/dbp"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	router.POST("/uploadImg", UploadImg)
	router.POST("/welcomeForm", WelcomeForm)
	router.POST("/accountForm", AccountForm)

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

func UploadImg(c *gin.Context) {
	db := dbp.DB

	file, handler, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error getting form file: %v", err)})
		log.Printf("Error getting form file: %v", err)
		return
	}
	defer file.Close()

	dstDir := "static/img/"
	if _, err := os.Stat(dstDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dstDir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error creating directory: %v", err)})
			log.Printf("Error creating directory: %v", err)
			return
		}
	}

	dstPath := fmt.Sprintf("%s%s", dstDir, handler.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error creating file: %v", err)})
		log.Printf("Error creating file: %v", err)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error copying file: %v", err)})
		log.Printf("Error copying file: %v", err)
		return
	}

	userId := c.Query("id")

	_ = db.Exec("UPDATE users SET image = ? WHERE id = ?", dstPath, userId)

	c.JSON(http.StatusCreated, gin.H{"image": dstPath})
	log.Printf("Successfully uploaded image: %s", dstPath)
}

func WelcomeForm(c *gin.Context) {
	db := dbp.DB

	birthday := c.PostForm("birthday")
	genre := c.PostForm("genre")
	sport := c.PostForm("sport")
	profilePicture := c.PostForm("image")
	userId := c.Query("id")

	var selectedSports []string
	if sports, exists := c.Request.PostForm["sports"]; exists {
		selectedSports = sports
	} else {
		selectedSports = []string{}
	}

	sportsList := ""
	if len(selectedSports) > 0 {
		sportsList = fmt.Sprintf("%s", selectedSports[0])
		for i := 1; i < len(selectedSports); i++ {
			sportsList += fmt.Sprintf(",%s", selectedSports[i])
		}
	}

	_ = db.Exec("UPDATE users SET gender = ?, secondary_sport = ? date_of_birth = ?, sport = ?, image = ? WHERE id = ?", genre, sportsList, birthday, sport, profilePicture, userId)

	c.JSON(http.StatusOK, gin.H{"message": "User information updated successfully"})
	log.Printf("User information updated successfully: %s", userId)
}

func AccountForm(c *gin.Context) {
	db := dbp.DB
	userId, _ := strconv.Atoi(c.Query("id"))

	user := &dbp.User{}
	tx := db.First(user, userId)
	if tx.RowsAffected > 0 {
		c.JSON(http.StatusOK, gin.H{"data": user})
	} else {
		fmt.Fprintln(c.Writer, "User not found")
	}

	username := c.PostForm("username")
	if username == "" {
		username = user.Username
	}
	email := c.DefaultPostForm("email", user.Email)
	if email == "" {
		email = user.Email
	}
	biography := c.DefaultPostForm("biography", user.Biography)
	if biography == "" {
		biography = user.Biography
	}
	sport := c.DefaultPostForm("sport", user.Sport)
	if sport == "" {
		sport = user.Sport
	}
	profilePicture := c.DefaultPostForm("image", user.Image)
	if profilePicture == "" {
		profilePicture = user.Image
	}
	city := c.DefaultPostForm("location", user.City)
	if city == "" {
		city = user.City
	}

	var selectedSports []string
	if sports, exists := c.Request.PostForm["sports"]; exists {
		selectedSports = sports
	} else {
		selectedSports = []string{}
	}

	sportsList := strings.Join(selectedSports, ",")
	log.Println(sportsList)

	// result := db.Exec("UPDATE users SET secondary_sport = ?, username = ?, email = ?, biography = ?, sport_id = ?, image = ?, city = ? WHERE id = ?", sportsList, username, email, biography, sport, profilePicture, city, userId)
	result := db.Model(&user).Updates(dbp.User{SecondarySport: sportsList, Username: username, Email: email, Biography: biography, Sport: sport, Image: profilePicture, City: city})
	log.Println(result.RowsAffected)
	log.Printf("User information updated successfully: %s", userId)
	// c.JSON(http.StatusOK, gin.H{"data": "User information updated successfully"})
}
