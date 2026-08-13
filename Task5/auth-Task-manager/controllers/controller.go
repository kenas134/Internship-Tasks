package controllers

import (
	"net/http"
	"time"

	"auth-Task-Manager/data"
	"auth-Task-Manager/middleware"
	"auth-Task-Manager/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Controller struct {
	UserService *data.UserService
	TaskService *data.TaskService
}

func NewController(
	userService *data.UserService,
	taskService *data.TaskService,
) *Controller {

	return &Controller{
		UserService: userService,
		TaskService: taskService,
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// POST /register
func (controller *Controller) Register(
	ctx *gin.Context,
) {

	var request RegisterRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {

		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	user, err := controller.UserService.CreateUser(
		request.Username,
		request.Password,
	)

	if err != nil {

		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	ctx.JSON(
		http.StatusCreated,
		gin.H{
			"message": "user registered successfully",
			"user":    user,
		},
	)
}

// POST /login
func (controller *Controller) Login(
	ctx *gin.Context,
) {

	var request LoginRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {

		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	user, err := controller.UserService.GetUserByUsername(
		request.Username,
	)

	if err != nil {

		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid username or password",
			},
		)

		return
	}

	err = controller.UserService.CheckPassword(
		user,
		request.Password,
	)

	if err != nil {

		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid username or password",
			},
		)

		return
	}

	expirationTime := time.Now().Add(
		24 * time.Hour,
	)

	claims := middleware.Claims{

		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				expirationTime,
			),

			IssuedAt: jwt.NewNumericDate(
				time.Now(),
			),

			Subject: user.ID,
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	tokenString, err := token.SignedString(
		middleware.JWTSecret(),
	)

	if err != nil {

		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "failed to generate token",
			},
		)

		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message":      "login successful",
			"access_token": tokenString,
			"token_type":   "Bearer",
			"expires_at":   expirationTime,
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"role":     user.Role,
			},
		},
	)
}

// GET /tasks
func (controller *Controller) GetTasks(
	ctx *gin.Context,
) {

	tasks, err := controller.TaskService.GetTasks()

	if err != nil {

		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"tasks": tasks,
		},
	)
}

// GET /tasks/:id
func (controller *Controller) GetTask(
	ctx *gin.Context,
) {

	id := ctx.Param("id")

	task, err := controller.TaskService.GetTaskByID(id)

	if err != nil {

		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	ctx.JSON(
		http.StatusOK,
		task,
	)
}

// POST /tasks
func (controller *Controller) CreateTask(
	ctx *gin.Context,
) {

	var task models.Task

	if err := ctx.ShouldBindJSON(&task); err != nil {

		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	createdTask, err :=
		controller.TaskService.CreateTask(task)

	if err != nil {

		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	ctx.JSON(
		http.StatusCreated,
		createdTask,
	)
}

// PUT /tasks/:id
func (controller *Controller) UpdateTask(
	ctx *gin.Context,
) {

	id := ctx.Param("id")

	var task models.Task

	if err := ctx.ShouldBindJSON(&task); err != nil {

		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	updatedTask, err :=
		controller.TaskService.UpdateTaskByID(
			id,
			task,
		)

	if err != nil {

		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	ctx.JSON(
		http.StatusOK,
		updatedTask,
	)
}

// DELETE /tasks/:id
func (controller *Controller) DeleteTask(
	ctx *gin.Context,
) {

	id := ctx.Param("id")

	err := controller.TaskService.DeleteTaskByID(id)

	if err != nil {

		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "task deleted successfully",
		},
	)
}