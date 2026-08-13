package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`

	jwt.RegisteredClaims
}

// Temporary refresh-token storage.
//
// In a real application this should be stored in MongoDB
// or another persistent store.
var refreshTokens = make(map[string]uint)

// --------------------------------------------------
// REGISTER
// --------------------------------------------------


//steps to register
//1.accept req of RegisterRequest(email and password) using shouldbind
//2.check if email already exists
//3.hash pass word using bcrypt.GenerateFromPassword
//
func Register(c *gin.Context) {

	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userMutex.Lock()
	defer userMutex.Unlock()

	// Check whether email already exists.
	for _, user := range users {
		if user.Email == req.Email {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Email already registered",
			})
			return
		}
	}

	// Hash password.
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Create user.
	user := User{
		ID:       nextUserID,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "user",
	}

	nextUserID++

	users = append(users, user)

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}

// --------------------------------------------------
// LOGIN
// --------------------------------------------------
//steps to login
//1.accept the loginRequest
//2.check if user exists
//3.compare the password


func Login(c *gin.Context) {

	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userMutex.Lock()

	var user *User

	for i := range users {
		if users[i].Email == req.Email {
			user = &users[i]
			break
		}
	}

	userMutex.Unlock()

	// Don't reveal whether the email exists.
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Compare entered password with stored hash.
	err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Create access token.
	accessToken, err := createAccessToken(user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create access token",
		})
		return
	}

	// Create refresh token.
	refreshToken, err := createRefreshToken(user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Login successful",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
		"expires_in":    15 * 60,
	})
}

// --------------------------------------------------
// CREATE ACCESS TOKEN
// --------------------------------------------------
//steps to create access token
//1.create  claim



func createAccessToken(user *User) (string, error) {

	secret := os.Getenv("JWT_SECRET")

	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(15 * time.Minute),
			),

			IssuedAt: jwt.NewNumericDate(
				time.Now(),
			),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString([]byte(secret))
}

// --------------------------------------------------
// CREATE REFRESH TOKEN
// --------------------------------------------------

func createRefreshToken(user *User) (string, error) {

	secret := os.Getenv("JWT_REFRESH_SECRET")

	claims := jwt.RegisteredClaims{
		Subject: "refresh",

		ExpiresAt: jwt.NewNumericDate(
			time.Now().Add(7 * 24 * time.Hour),
		),

		IssuedAt: jwt.NewNumericDate(
			time.Now(),
		),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	// Store which user owns this refresh token.
	userMutex.Lock()
	refreshTokens[tokenString] = user.ID
	userMutex.Unlock()

	return tokenString, nil
}

// --------------------------------------------------
// REFRESH
// --------------------------------------------------

func Refresh(c *gin.Context) {

	var req RefreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	refreshToken := req.RefreshToken

	// Check whether we know this refresh token.
	userMutex.Lock()

	userID, exists := refreshTokens[refreshToken]

	userMutex.Unlock()

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid refresh token",
		})
		return
	}

	// Parse refresh token.
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(
		refreshToken,
		claims,
		func(token *jwt.Token) (interface{}, error) {

			if token.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrSignatureInvalid
			}

			return []byte(os.Getenv("JWT_REFRESH_SECRET")), nil
		},
	)

	if err != nil || !token.Valid {
		// Remove invalid/expired token.
		userMutex.Lock()
		delete(refreshTokens, refreshToken)
		userMutex.Unlock()

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid or expired refresh token",
		})
		return
	}

	// Find the user.
	userMutex.Lock()

	var user *User

	for i := range users {
		if users[i].ID == userID {
			user = &users[i]
			break
		}
	}

	userMutex.Unlock()

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})
		return
	}

	// Create a new access token.
	accessToken, err := createAccessToken(user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create access token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"token_type":   "Bearer",
		"expires_in":   15 * 60,
	})
}

// --------------------------------------------------
// LOGOUT
// --------------------------------------------------

func Logout(c *gin.Context) {

	var req RefreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userMutex.Lock()

	delete(refreshTokens, req.RefreshToken)

	userMutex.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}