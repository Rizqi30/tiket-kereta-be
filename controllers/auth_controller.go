package controllers

import (
	"net/http"
	"tiket-kereta-be/config"
	"tiket-kereta-be/models"
	"tiket-kereta-be/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// struct to catch the JSON from user when Register
type RegisterInput struct {
	Name 			string `json:"name" binding:"required"`
	Email 			string `json:"email" binding:"required,email"`
	Phone 			string `json:"phone" binding:"required"`
	Password 		string `json:"password" binding:"required,min=6"`
	Role 			string `json:"role"`
}

func Register(c *gin.Context) {
	var input RegisterInput

	// valdasi input JSON from request
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check if email already exist
	var user models.User
	if config.DB.Where("email = ?", input.Email).First(&user).Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is already exists"})
		return
	}

	// check if phone already exist
	if config.DB.Where("phone = ?", input.Phone).First(&user).Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone is already exists"})
		return
	}

	// hashing password use bcrypt 
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed generate password"})
		return
	}

	// check role
	userRole := "PENUMPANG"
	if input.Role == "ADMIN" {
		userRole = "ADMIN"
	}

	// saving data to database via GORM
	user = models.User{
		Name:     input.Name,
		Email:    input.Email,
		Phone:    input.Phone,
		Password: string(hashedPassword),
		Role: userRole,
	}

	// check error insert user
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Success Register",
		"user": gin.H{
			"id": user.ID,
			"name": user.Name,
			"email": user.Email,
			"phone": user.Phone,
			"role": user.Role,
		},
	})
}

// struct to catch the JSON from user when Login
type LoginInput struct {
	Email 		string `json:"email" binding:"required,email"`
	Password 	string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var input LoginInput
	var user models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// search user by email
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email or password incorrect"})
		return
	}

	// verify password by bcrypt
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email or password incorrect"})
		return
	}

	// generate JWT token after success login
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success Login",
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"phone": user.Phone,
			"role":  user.Role,
		},
	})
}

// GetMe - take user profile who login based on JWT token
func GetMe(c *gin.Context) {
	// take user_id from context (already di-set by AuthMiddleware)
	userID, _ := c.Get("user_id")

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id": user.ID,
			"name": user.Name,
			"email": user.Email,
			"phone": user.Phone,
			"role": user.Role,
		},
	})
}