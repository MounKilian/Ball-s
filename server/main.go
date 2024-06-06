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

func main() {
	// addUsers()
	// addSwipe()
	// db := dbp.DB
	// users := []dbp.User{}
	// db.Find(&users, &dbp.User{})
	// usertosortfrom := rand.Intn(len(users))
	// sort(users[usertosortfrom])
	router := gin.Default()
	// sort()
	router.LoadHTMLGlob("pages/*.html")
	router.Static("/static", "./static")

	// Routes pour les pages HTML
	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", nil)
	})

	router.GET("/profilUser", func(c *gin.Context) {
		c.HTML(http.StatusOK, "profilUser.html", nil)
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

	// Routes pour les API
	router.POST("/register", handleRegister)
	router.POST("/login", handleLogin)

	// Route de test pour la fonction sort
	router.GET("/form", func(c *gin.Context) {
		c.HTML(http.StatusOK, "welcomePage.html", nil)
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}

	// Pour le test: affiche les utilisateurs triés
	fmt.Println("Utilisateur de départ:", users[usertosortfrom].ID)
	fmt.Println("Utilisateurs triés:")
	for _, user := range sortedUsers {
		fmt.Println(user.ID)
	}
}

func handleRegister(c *gin.Context) {
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
}

func handleLogin(c *gin.Context) {
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
			Name:   "user_id",
			Value:  userID,
			MaxAge: 3600,
		}
		http.SetCookie(c.Writer, cookie)
		// Après la définition de handleLogin

		c.JSON(http.StatusOK, gin.H{"message": "Authentification réussie"})
		}

		// Définir les routes en dehors de la fonction main
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

	for i := 0; i < 300; i++ {
		db := dbp.DB
		sportid := rand.Intn(17)
		cityrand := rand.Intn(2)
		genderand := rand.Intn(2)
		genderprefrand := rand.Intn(2)
		var cityname string
		var gender string
		var genderpref string
		if cityrand == 0 {
			cityname = "paris"
		} else {
			cityname = "Lyon"
		}

		if genderand == 0 {
			gender = "men"
		} else {
			gender = "women"
		}

		if genderprefrand == 0 {
			genderpref = "men"
		} else {
			genderpref = "women"
		}
		user := dbp.User{
			Username:      "User",
			DateOfBirth:   time.Now(),
			Sport:         sport.Name,
			Gender:        gender,
			DesiredGender: genderpref,
			City:          cityname,
		}

		db.Create(&user)

	}
}

func addSwipe() {

	db := dbp.DB
	users := []dbp.User{}
	swiped := []int64{}

	db.Not(&swiped).Find(&users)

	if len(users) == 0 {
		fmt.Println("No users found")
		return
	}

	rand.Seed(time.Now().UnixNano())
	startUser := rand.Intn(len(users))
	var potential []dbp.User

	for i := 0; i < len(users); i++ {
		if users[startUser].City == users[i].City && i != startUser {
			if users[startUser].Sport == users[i].Sport {
				potential = append(potential, users[i])
			}
		}
	}

	rand.Seed(time.Now().UnixNano())
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
