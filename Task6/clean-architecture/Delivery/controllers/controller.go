package controllers

import (
	domain "clean-architecture/Domain"
	usecases "clean-architecture/Usecases"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	taskUsecase *usecases.TaskUsecase
	userUsecase *usecases.UserUsecase
}

func NewController(
	taskUsecase *usecases.TaskUsecase,
	userUsecase *usecases.UserUsecase,
) *Controller {
	return &Controller{
		taskUsecase: taskUsecase,
		userUsecase: userUsecase,
	}
}

type RegisterRequest struct {
	ID        string `json:"_id,omitempty"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Username  string `json:"username"`
	Password  string `json:"password,omitempty"`
	Role      string `json:"role"`
}


type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TaskRequest struct {
	ID          string      `json:"_id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	DueDate     time.Time   `json:"duedate"`
	Status      string      `json:"status"`
}





func (c *Controller) Register(ctx *gin.Context) {

	var request RegisterRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	user := domain.User{
		Firstname: request.Firstname,
		Lastname:  request.Lastname,
		Username:  request.Username,
		Password:  request.Password,
		Role:      request.Role,
	}

	newUser,err := c.userUsecase.Register(user)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result := RegisterRequest{
		ID : newUser.ID,
		Firstname: newUser.Firstname,
		Lastname:  newUser.Lastname,
		Username:  newUser.Username,
		Password:  newUser.Password,
		Role:      newUser.Role,
	}


	ctx.JSON(http.StatusCreated,gin.H{
		"message": "user created successfully",
		"user": result,
	})

}



func (c *Controller) Login(ctx *gin.Context) {
	var request LoginRequest

	if err:= ctx.ShouldBindJSON(&request);err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	accessToken, refreshToken, err := c.userUsecase.Login(
		request.Username,
		request.Password,
	)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Logged In successfully",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}


func (c *Controller) Profile(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	id, ok := userID.(string)

	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user ID",
		})
		return
	}

	user, err := c.userUsecase.Profile(id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, user)
}


func (c *Controller) ListUsers(ctx *gin.Context) {

	users, err := c.userUsecase.ListUsers()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get users",
		})
		return
	}

	ctx.JSON(http.StatusOK, users)
}


func (c *Controller) GetTasks(ctx *gin.Context) {

	tasks, err := c.taskUsecase.GetTasks()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get tasks",
		})
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}


func (c *Controller) GetTaskByID(ctx *gin.Context) {
	id := ctx.Param("id")

	task, err := c.taskUsecase.GetTaskByID(id)

	if err != nil {
		ctx.JSON(http.StatusNotFound,gin.H{
			"error": "task not found",
		})
		return
	}

	ctx.JSON(http.StatusOK,task)
}


func (c *Controller) CreateTask(ctx *gin.Context) {
	var request TaskRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	task := domain.Task{  
		Title:          request.Title,
		Description:    request.Description,
		DueDate:        request.DueDate,
		Status:         request.Status,      
		
	}

	newTask,err := c.taskUsecase.CreateTask(task)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result := TaskRequest{  
		ID:             newTask.ID,
		Title:          newTask.Title,
		Description:    newTask.Description,
		DueDate:        newTask.DueDate,
		Status:         newTask.Status,      
	}

		
	ctx.JSON(http.StatusCreated,gin.H{
		"message": "Task created successfully",
		"user": result,
	})
}


func (c *Controller)  UpdateTask(ctx *gin.Context) {
	var request TaskRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	task := domain.Task{  
		ID:             request.ID,
		Title:          request.Title,
		Description:    request.Description,
		DueDate:        request.DueDate,
		Status:         request.Status,      
	}

	err := c.taskUsecase.UpdateTask(task.ID,task)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}


	ctx.JSON(http.StatusCreated,gin.H{
		"message": "Task Updated successfully",
	})
}

func (c *Controller) DeleteTask(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.taskUsecase.DeleteTask(id)

	if err != nil {
		ctx.JSON(http.StatusNotFound,gin.H{
			"error": "task not found",
		})
		return
	}

	ctx.JSON(http.StatusOK,gin.H{
		"message":"Task deleted successfully!",
	})
}
