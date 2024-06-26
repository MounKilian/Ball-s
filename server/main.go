package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"balls/dbp"

	"github.com/gin-gonic/gin"
)

// func main() {
// 	addSport()
// }

func main() {
	// addUsers()
	// addSwipe()
	// db := dbp.DB
	// users := []dbp.User{}
	// db.Find(&users, &dbp.User{})
	// usertosortfrom := rand.Intn(len(users))
	// sort(users[usertosortfrom])
	router := gin.Default()

	router.LoadHTMLGlob("pages/*.html")

	router.Static("/static", "./static")

	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	router.GET("/", func(c *gin.Context) {

		c.HTML(http.StatusOK, "home.html", nil)
	})

	router.GET("/profilUser", func(c *gin.Context) {
		c.HTML(http.StatusOK, "profilUser.html", nil)
	})

	router.GET("/profilOther", func(c *gin.Context) {
		c.HTML(http.StatusOK, "profilOther.html", nil)
	})

	router.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	})

	router.GET("/message", func(c *gin.Context) {
		c.HTML(http.StatusOK, "message.html", nil)
	})

	router.GET("/discussion", func(c *gin.Context) {
		c.HTML(http.StatusOK, "discussion.html", nil)
	})

	router.POST("/register", func(c *gin.Context) {
		var user struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&user); err != nil {
			log.Println("Erreur de liaison des données :", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Println("Données utilisateur reçues :", user)

		err := dbp.RegisterUser(user.Username, user.Email, user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'enregistrement de l'utilisateur"})
			return
		}

		userID, _, _ := dbp.AuthenticateUser(user.Username, user.Password)

		cookie := &http.Cookie{
			Name:     "user_id",
			Value:    userID,
			SameSite: http.SameSiteStrictMode,
			// 	HttpOnly: true,
			MaxAge: 3600 * 24 * 30,
		}
		http.SetCookie(c.Writer, cookie)

		// http.Redirect(c.Writer, c.Request, "/form?id="+userID, http.StatusFound)
		c.JSON(http.StatusOK, gin.H{"message": userID})
	})

	router.POST("/login", func(c *gin.Context) {
		var user struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, isAuthenticated, err := dbp.AuthenticateUser(user.Username, user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if !isAuthenticated {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Nom d'utilisateur ou mot de passe incorrect"})
			return
		}

		// Création d'un cookie avec l'ID utilisateur
		cookie := &http.Cookie{
			Name:     "user_id",
			Value:    userID,
			SameSite: http.SameSiteStrictMode,
			// 	HttpOnly: true,
			MaxAge: 3600 * 24 * 30,
		}
		http.SetCookie(c.Writer, cookie)

		c.JSON(http.StatusOK, gin.H{"message": "Authentification réussie"})
	})

	router.GET("/form", func(c *gin.Context) {
		c.HTML(http.StatusOK, "welcomePage.html", nil)
	})

	router.GET("/account", func(c *gin.Context) {
		c.HTML(http.StatusOK, "accountPage.html", nil)
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
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
	for _, v := range []string{"Football", "Basketball", "Tennis", "Baseball", "Surf", "Volley", "Pingpong", "Golf", "Natation", "Rugby", "Bowling", "Handball", "Escalade", "Cyclisme", "Sauts", "Plongée", "Acrobranche", "Tyroliènne", "Course", "Musculation", "Randonnée", "Paddle", "Acrobranche", "Ski", "Boxe", "MMA", "Kapoera", "Pétanque", "Gymnastique", "Danse", "Karting", "Paintball", "Judo", "Karaté", "Escrime", "Ultimate", "LaserGame", "Je ne fait pas que du sport"} {

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
		sport := &dbp.Stat{ID: uint(rand.Intn(17))}
		db.First(&sport, sport.ID)
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
			Sport:       sport.Name,
			City:        cityname,
		}
		db.Create(&user)
	}
}

func profileUser(c *gin.Context) {
	user := dbp.User{
		Username:    "username",
		Email:       "email@email.com",
		Password:    "Password",
		Image:       "uglyprofilpic.png",
		Biography:   "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum",
		Sport:       "Football",
		DateOfBirth: time.Now(),
		City:        "69420",
	}

	c.HTML(http.StatusOK, "profilUser.html", user)

}
