package main

import (
	"fmt"
	"math/rand"
	"time"

	"balls/dbp"
)

// func main() {
// 	addSport()
// }

func main() {
	// addUsers()
	// addSwipe()
	db := dbp.DB
	users := []dbp.User{}
	db.Find(&users, &dbp.User{})
	usertosortfrom := rand.Intn(len(users))
	sort(users[usertosortfrom])
	// router := gin.Default()

	// router.LoadHTMLGlob("pages/*.html")

	// router.Static("/static", "./static")

	// router.GET("/login", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "login.html", nil)
	// })

	// router.GET("/", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "home.html", nil)
	// })

	// router.GET("/register", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "register.html", nil)
	// })

	// router.POST("/register", func(c *gin.Context) {
	// 	var user struct {
	// 		Username string `json:"username"`
	// 		Email    string `json:"email"`
	// 		Password string `json:"password"`
	// 	}
	// 	if err := c.ShouldBindJSON(&user); err != nil {
	// 		log.Println("Erreur de liaison des données :", err)
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	log.Println("Données utilisateur reçues :", user)

	// 	err := dbp.RegisterUser(user.Username, user.Email, user.Password)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'enregistrement de l'utilisateur"})
	// 		return
	// 	}

	// 	c.JSON(http.StatusOK, gin.H{"message": "Utilisateur enregistré avec succès"})
	// })

	// router.POST("/login", func(c *gin.Context) {
	// 	var user struct {
	// 		Username string `json:"username"`
	// 		Password string `json:"password"`
	// 	}
	// 	if err := c.ShouldBindJSON(&user); err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	userID, isAuthenticated, err := dbp.AuthenticateUser(user.Username, user.Password)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	if !isAuthenticated {
	// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Nom d'utilisateur ou mot de passe incorrect"})
	// 		return
	// 	}

	// 	// Création d'un cookie avec l'ID utilisateur
	// 	cookie := &http.Cookie{
	// 		Name:     "user_id",
	// 		Value:    userID,
	// 		SameSite: http.SameSiteStrictMode,
	// 		// 	HttpOnly: true,
	// 		MaxAge: 3600, // Durée de vie sdu cookie en secondes (1 heure ici)
	// 	}
	// 	http.SetCookie(c.Writer, cookie)

	// 	c.JSON(http.StatusOK, gin.H{"message": "Authentification réussie"})
	// })

	// router.GET("/form", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "welcomePage.html", nil)
	// })

	// if err := router.Run(":8080"); err != nil {
	// 	log.Fatal(err)
	// }
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
			SportID:       sportid,
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
	swipes := []int{}

	db.Find(&swipes, &dbp.Swipe{})
	db.Find(&users)

	for f := 0; f < len(users); f++ {

		var potential []dbp.User
		potential = nil
		for i := 0; i < len(users); i++ {
			if users[f].City == users[i].City && users[i].ID != users[f].ID {
				if users[f].SportID == users[i].SportID && users[f].DesiredGender == users[i].Gender {
					potential = append(potential, users[i])
				}
			}
		}

		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(potential), func(i, j int) { potential[i], potential[j] = potential[j], potential[i] })
		// result, _ := json.Marshal(&users)
		// resultpot, _ := json.Marshal((&potential))

		swipe := dbp.Swipe{
			UserAID: int(users[f].ID),
			UserA:   users[f],
			UserBID: int(potential[0].ID),
			UserB:   potential[0],
		}
		db.Create(&swipe)

	}
}

func sort(startUser dbp.User) []dbp.User {
	db := dbp.DB
	users := []dbp.User{}
	swipes := []int{}
	db.Model(&dbp.Swipe{}).Where(&dbp.Swipe{UserAID: int(startUser.ID)}).Pluck("user_b_id", &swipes)

	db.Not(swipes).Find(&users)

	fmt.Println(swipes)

	var potential []dbp.User
	for i := 0; i < len(users); i++ {
		if startUser.City == users[i].City && users[i].ID != startUser.ID {
			if startUser.SportID == users[i].SportID && startUser.DesiredGender == users[i].Gender {
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
