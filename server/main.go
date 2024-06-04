package main

import (
	"log"
	"net/http"

	"balls/dbp"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("pages/*.html")

	router.Static("/static", "./static")

	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", nil)
	})

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

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
