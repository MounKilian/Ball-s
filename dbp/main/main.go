package main

import (
	"balls/dbp"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
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
	router.GET("/sort", func(c *gin.Context) {
		db := dbp.DB
		user := dbp.User{}
		s, ok := c.GetQuery("id")
		if !ok {
			fmt.Fprintln(c.Writer, "Id not provided in get request")
			return
		}
		id, err := strconv.Atoi(s)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error converting id to int: %v", err)})
			return
		}
		db.First(&user, id)
		sortedUsers := sort(user)
		c.JSON(http.StatusFound, gin.H{"message": sortedUsers})
	})

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
	// genre := c.PostForm("genre")
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

	_ = db.Exec("UPDATE users SET date_of_birth = ?, sport_id = ?, image = ? WHERE id = ?", birthday, sport, profilePicture, userId)

	c.JSON(http.StatusOK, gin.H{"message": "User information updated successfully"})
	log.Printf("User information updated successfully: %s", userId)
}

func sort(startUser dbp.User) []dbp.User {
	db := dbp.DB
	users := []dbp.User{}
	swipes := []int{}
	db.Model(&dbp.Swipe{}).Where(&dbp.Swipe{UserAID: int(startUser.ID)}).Pluck("user_b_id", &swipes)

	db.Select("id", "username", "biography", "gender", "sport", "secondary_sport", "image", "city", "date_of_birth").Not(swipes).Find(&users)

	fmt.Println(swipes)

	var potential []dbp.User

	for i := 0; i < len(users); i++ {
		if startUser.City == users[i].City && users[i].ID != startUser.ID {
			if startUser.Sport == users[i].Sport && startUser.DesiredGender == users[i].Gender {
				potential = append(potential, users[i])
			}
		}
	}

	rand.Shuffle(len(potential), func(i, j int) { potential[i], potential[j] = potential[j], potential[i] })
	// result, _ := json.Marshal(&users)
	// resultpot, _ := json.Marshal((&potential))
	fmt.Println(startUser.ID)
	fmt.Println("potential")
	for i := 0; i < len(potential); i++ {
		fmt.Println(potential[i].ID)
	}
	return potential
}
