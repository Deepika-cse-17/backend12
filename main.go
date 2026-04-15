package main

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

var firestoreClient *firestore.Client

func main() {
	ctx := context.Background()

	opt := option.WithCredentialsFile("serviceAccountKey.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("Failed to init Firebase: %v", err)
	}

	firestoreClient, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Failed to init Firestore: %v", err)
	}
	defer firestoreClient.Close()

	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/signup", func(c *gin.Context) {
		var user struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if user already exists
		docs, err := firestoreClient.Collection("users").
			Where("email", "==", user.Email).
			Documents(ctx).GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if len(docs) > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already in use"})
			return
		}

		_, _, err = firestoreClient.Collection("users").Add(ctx, map[string]interface{}{
			"email":    user.Email,
			"password": user.Password,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User created"})
	})

	r.POST("/login", func(c *gin.Context) {
		var user struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		docs, err := firestoreClient.Collection("users").
			Where("email", "==", user.Email).
			Documents(ctx).GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if len(docs) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		data := docs[0].Data()
		storedPass, _ := data["password"].(string)
		if storedPass != user.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
	})

	r.GET("/users", func(c *gin.Context) {
		docs, err := firestoreClient.Collection("users").Documents(ctx).GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		users := make([]map[string]string, 0, len(docs))
		for _, doc := range docs {
			data := doc.Data()
			email, _ := data["email"].(string)
			users = append(users, map[string]string{"email": email})
		}

		c.JSON(http.StatusOK, users)
	})

	log.Println("Server starting on :8081")
	r.Run(":8081")
}
