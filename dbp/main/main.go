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
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var Letters = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "T", "U", "V", "W", "X", "Y", "Z"}

func selectRandomLetter() string {
	randomIndex := rand.Intn(len(Letters) - 1)
	return Letters[randomIndex]
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

var clients = make(map[*Client]bool)

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
	router.POST("/strikeOrMiss", StrikeOrMiss)
	router.GET("/matchs", getAllMatches)

	router.GET("/ws", handleWebSocket)
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

func StrikeOrMiss(c *gin.Context) {
	db := dbp.DB
	userId, _ := strconv.Atoi(c.PostForm("id"))
	otherUserId, _ := strconv.Atoi(c.PostForm("otherUserId"))
	decision := c.PostForm("decision")

	if decision == "miss" {
		miss := dbp.Miss{UserAID: userId, UserBID: otherUserId}
		db.Create(&miss)
	} else if decision == "strike" {
		strike := dbp.Strike{UserAID: userId, UserBID: otherUserId}
		db.Create(&strike)
	}

	c.JSON(http.StatusOK, gin.H{"message": "User information updated successfully"})
}

func handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error during connection upgrade:", err)
		return
	}
	defer conn.Close()

	client := &Client{conn: conn, send: make(chan []byte)}
	clients[client] = true

	go client.writePump()

	client.readPump()
}

func (client *Client) readPump() {
	defer func() {
		delete(clients, client)
		client.conn.Close()
	}()

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		for otherClient := range clients {
			if otherClient != client {
				otherClient.send <- message
			}
		}
	}
}

func (client *Client) writePump() {
	defer func() {
		client.conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
			if !ok {
				return
			}
			err := client.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return
			}
		}
	}
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
		c.JSON(http.StatusNotFound, "User not found")
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

func getAllMatches(c *gin.Context) {
	db := dbp.DB
	matches := []dbp.Match{}
	db.Preload("UserA").Preload("UserB").Find(&matches)
	result, _ := json.Marshal(&matches)
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
		sportsList = selectedSports[0]
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
	log.Printf("User information updated successfully: %d", userId)
	// c.JSON(http.StatusOK, gin.H{"data": "User information updated successfully"})
}

func sort(startUser dbp.User) []dbp.User {
	db := dbp.DB
	users := []dbp.User{}
	misses := []int{}
	strikes := []int{}
	db.Model(&dbp.Miss{}).Where(&dbp.Miss{UserAID: int(startUser.ID)}).Pluck("user_b_id", &misses)
	db.Model(&dbp.Strike{}).Where(&dbp.Strike{UserAID: int(startUser.ID)}).Pluck("user_b_id", &strikes)

	db.Select("id", "username", "biography", "gender", "sport", "secondary_sport", "image", "city", "date_of_birth").Not(misses).Not(strikes).Find(&users)

	fmt.Println(misses)

	var potential []dbp.User

	for i := 0; i < len(users); i++ {
		if startUser.City == users[i].City && users[i].ID != startUser.ID {
			if startUser.Sport == users[i].Sport && startUser.DesiredGender == users[i].Gender {
				potential = append(potential, users[i])
			}
		}
	}

	for i := 0; i < 21; i++ {
		var usernumber int64
		listuser := []dbp.User{}
		db.Model(&dbp.User{}).Where("gender = ?", startUser.DesiredGender).Where("city = ?", startUser.City).Not("id = ?", startUser.ID).Not("sport = ?", startUser.Sport).Not(misses).Not(strikes).Count(&usernumber)
		db.Model(&dbp.User{}).Where("gender = ?", startUser.DesiredGender).Where("city = ?", startUser.City).Not("id = ?", startUser.ID).Not(misses).Not(strikes).Find(&listuser)
		// db.Model(&dbp.User{}).Where("gender = ?", startUser.DesiredGender).Where("city = ?", startUser.City).Not("id = ?", startUser.ID).Find(users)
		usertoadd := rand.Int63n(usernumber)
		potential = append(potential, users[usertoadd])
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
