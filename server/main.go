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
	router := gin.Default()
	sort()
	router.LoadHTMLGlob("pages/*.html")

	router.Static("/static", "./static")

	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", nil)
	})

	router.GET("/profilUser", profileUser)

	router.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
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

		c.JSON(http.StatusOK, gin.H{"message": "Utilisateur enregistré avec succès"})
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
			Path:     "/",
			HttpOnly: true,
			MaxAge:   3600, // Durée de vie du cookie en secondes (1 heure ici)
		}
		http.SetCookie(c.Writer, cookie)

		c.JSON(http.StatusOK, gin.H{"message": "Authentification réussie"})
	})

	router.GET("/form", func(c *gin.Context) {
		c.HTML(http.StatusOK, "welcomePage.html", nil)
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

func sort() {
	db := dbp.DB
	users := []dbp.User{}
	swiped := []int64{}
	db.Not(&swiped).Find(&users)
	startUser := rand.Intn(50)
	var potential []dbp.User
	for i := 0; i < len(users); i++ {
		if users[startUser].City == users[i].City && i != startUser {
			if users[startUser].Sport == users[i].Sport {
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

func profileUser(c *gin.Context) {
	user := dbp.User{
		Username:  "username",
		Email:     "email@email.com",
		Password:  "Password",
		Image:     "uglyprofilpic.png",
		Biography: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum",
		SportID:   1,
		Sport: dbp.Stat{
			ID:         2,
			Name:       "pingpong",
			Catégories: "triple champion du monde",
		},
		DateOfBirth: time.Now(),
		City:        "69420",
	}

	c.HTML(http.StatusOK, "profilUser.html", user)

}
