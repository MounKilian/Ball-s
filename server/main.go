package main

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
		"balls/dbp"

)

func main() {
    router := gin.Default()
    
    router.LoadHTMLGlob("web/*.html")
    
    router.Static("/static", "./static")
    
    router.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "home.html", nil)
    })
    
    router.GET("/login", func(c *gin.Context) {
        c.HTML(http.StatusOK, "login.html", nil)
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

        isAuthenticated, err := dbp.AuthenticateUser(user.Username, user.Password)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'authentification de l'utilisateur"})
            return
        }

        if !isAuthenticated {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Nom d'utilisateur ou mot de passe incorrect"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Authentification réussie"})
    })
    
    if err := router.Run(":8080"); err != nil {
        log.Fatal(err)
    }
}
