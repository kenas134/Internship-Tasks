package main

import (
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    "net/http"
)

type RegisterRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

func registerHandler(c * gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }


	hashedPassword,err := bcrypt.GenerateFromPassword([]byte(req.Password),12)
	if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to process password",
        })
        return
    }

	user := &User{
        Email:    req.Email,
        Password: string(hashedPassword),
        Role:     "user", // Default role
    }
	if err := userStore.Create(user); err != nil {
        c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
        return
    }

	user.Password = ""
    c.JSON(http.StatusCreated, gin.H{
        "message": "User registered successfully",
        "user":    user,
    })
}