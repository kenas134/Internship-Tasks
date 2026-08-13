package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// Load .env
	err := godotenv.Load()

	if err != nil {
		panic("Error loading .env")
	}

	router := gin.Default()

	// ----------------------------------------
	// PUBLIC ROUTES
	// ----------------------------------------

	router.POST("/register", Register)

	router.POST("/login", Login)

	router.POST("/refresh", Refresh)

	router.POST("/logout", Logout)

	// ----------------------------------------
	// AUTHENTICATED ROUTES
	// ----------------------------------------

	auth := router.Group("/")

	auth.Use(AuthMiddleware())

	auth.GET("/profile", Profile)
	auth.GET("/tasks", GetTasks)

	// ----------------------------------------
	// ADMIN ROUTES
	// ----------------------------------------

	admin := router.Group("/admin")

	admin.Use(
		AuthMiddleware(),
		RequireRole("admin"),
	)

	admin.GET("/users", GetUsers)

	admin.DELETE("/users/:id", DeleteUser)

	// ----------------------------------------
	// START SERVER
	// ----------------------------------------

	router.Run(":8080")
}

// ----------------------------------------
// PROFILE
// ----------------------------------------

func Profile(c *gin.Context) {

	userID, _ := c.Get("user_id")
	email, _ := c.Get("email")
	role, _ := c.Get("role")

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"email":   email,
		"role":    role,
	})
}

// ----------------------------------------
// GET TASKS
// ----------------------------------------

func GetTasks(c *gin.Context) {

	userID, _ := c.Get("user_id")

	c.JSON(http.StatusOK, gin.H{
		"message": "You can access tasks",
		"user_id": userID,
		"tasks": []string{
			"Learn Go",
			"Learn JWT",
			"Build API",
		},
	})
}

// ----------------------------------------
// ADMIN: GET USERS
// ----------------------------------------

func GetUsers(c *gin.Context) {

	userMutex.Lock()
	defer userMutex.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

// ----------------------------------------
// ADMIN: DELETE USER
// ----------------------------------------

func DeleteUser(c *gin.Context) {

	id := c.Param("id")

	c.JSON(http.StatusOK, gin.H{
		"message": "Delete user " + id,
	})
}