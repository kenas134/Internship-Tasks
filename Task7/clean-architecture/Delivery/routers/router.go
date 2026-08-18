package routers

import (
	"clean-architecture/Delivery/controllers"
	infrastructure "clean-architecture/Infrastructure"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	controller *controllers.Controller,
	jwtService infrastructure.JWTService,
) *gin.Engine {
	router := gin.Default()

	router.POST(
		"/register",
		controller.Register,
	)	

	router.POST(
		"/login",
		controller.Login,
	)
	// =========================
	// Authenticated routes
	// =========================

	protected := router.Group("/tasks")

	protected.Use(
		infrastructure.AuthMiddleware(jwtService),
	)

		protected.GET("/profile", controller.Profile)

		protected.GET("/tasks", controller.GetTasks)
		protected.GET("/tasks/:id", controller.GetTaskByID)
		protected.POST("/tasks", controller.CreateTask)
		protected.PUT("/tasks/:id", controller.UpdateTask)

		admin := router.Group("/admin")

	admin.Use(
		infrastructure.AuthMiddleware(jwtService),
		infrastructure.RequiredRole("admin"),
	)

	admin.GET("/users", controller.ListUsers)
	admin.DELETE("/tasks/:id", controller.DeleteTask)

	return router
}