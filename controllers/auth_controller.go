package controllers

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"go-supabase-api/config"
	"go-supabase-api/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Confirm  string `json:"confirm_password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"message": "Format JSON salah"})
		return
	}

	if req.Password != req.Confirm {
		c.JSON(400, gin.H{"message": "Password tidak sama"})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 10)

	config.SupabaseRequest("POST", "users", gin.H{
		"name":     req.Name,
		"email":    req.Email,
		"password": string(hash),
		"role":     "user",
	})

	c.JSON(201, gin.H{"message": "Register berhasil"})
}

func Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"message": "Format JSON salah"})
		return
	}

	resp, _ := config.SupabaseRequest(
		"GET",
		fmt.Sprintf("users?email=eq.%s&select=*", req.Email),
		nil,
	)

	var users []models.User
	json.NewDecoder(resp.Body).Decode(&users)

	if len(users) == 0 {
		c.JSON(401, gin.H{"message": "Email tidak ditemukan"})
		return
	}

	user := users[0]

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		c.JSON(401, gin.H{"message": "Password salah"})
		return
	}

	token, _ := generateJWT(user)

	// ⬇️ RESPONSE BARU
	c.JSON(200, gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func generateJWT(user models.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key" // fallback, should use env variable
	}

	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // 24 hours expiration
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
