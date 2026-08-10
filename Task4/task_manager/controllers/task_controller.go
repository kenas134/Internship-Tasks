package controllers

import (
	"net/http"

	"task_manager/data"
	"task_manager/models"

	"github.com/gin-gonic/gin"
)

// GET /tasks
func GetTasks(ctx *gin.Context) {
	tasks, err := data.GetTasks()

	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{"tasks": tasks},
	)
}

// GET /tasks/:id
func GetTaskByID(ctx *gin.Context) {
	id := ctx.Param("id")

	task, err := data.GetTaskByID(id)

	if err != nil {
		ctx.JSON(
			http.StatusNotFound,
			gin.H{"error": err.Error()},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		task,
	)
}

// POST /tasks
func CreateTask(ctx *gin.Context) {
	var newTask models.Task

	if err := ctx.ShouldBindJSON(&newTask); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	createdTask, err := data.CreateTask(newTask)

	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	ctx.JSON(
		http.StatusCreated,
		createdTask,
	)
}

// PUT /tasks/:id
func UpdateTaskByID(ctx *gin.Context) {
	id := ctx.Param("id")

	var updatedTask models.Task

	if err := ctx.ShouldBindJSON(&updatedTask); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	err := data.UpdateTaskByID(id, updatedTask)

	if err != nil {
		ctx.JSON(
			http.StatusNotFound,
			gin.H{"error": err.Error()},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{"message": "Task updated"},
	)
}

// DELETE /tasks/:id
func DeleteTaskByID(ctx *gin.Context) {
	id := ctx.Param("id")

	err := data.DeleteTaskByID(id)

	if err != nil {
		ctx.JSON(
			http.StatusNotFound,
			gin.H{"error": err.Error()},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{"message": "Task removed"},
	)
}
