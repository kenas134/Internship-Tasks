package router

import (
	"auth-Task-Manager/controllers"
	"auth-Task-Manager/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	controller *controllers.Controller,
) *gin.Engine {

	router := gin.Default()

	// =========================
	// Public routes
	// =========================

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
		middleware.AuthMiddleware(),
	)

	{
		protected.GET(
			"",
			controller.GetTasks,
		)

		protected.GET(
			"/:id",
			controller.GetTask,
		)

		protected.POST(
			"",
			controller.CreateTask,
		)

		protected.PUT(
			"/:id",
			controller.UpdateTask,
		)
	}

	// =========================
	// Admin routes
	// =========================

	admin := router.Group("/tasks")

	admin.Use(
		middleware.AuthMiddleware(),
		middleware.RequireRole("admin"),
	)

	{
		admin.DELETE(
			"/:id",
			controller.DeleteTask,
		)
	}

	return router
}