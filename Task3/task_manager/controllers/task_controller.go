package controllers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "task_manager/data"
    "task_manager/models"
)

func GetTasks(ctx *gin.Context) {

    tasks := data.GetTasks()

    ctx.JSON(
        http.StatusOK,
        gin.H{"tasks": tasks},
    )
}

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

func CreateTask(ctx *gin.Context) {

    var newTask models.Task

    if err := ctx.ShouldBindJSON(&newTask); err != nil {

        ctx.JSON(
            http.StatusBadRequest,
            gin.H{"error": err.Error()},
        )

        return
    }

    data.CreateTask(newTask)

    ctx.JSON(
        http.StatusCreated,
        gin.H{"message": "Task created"},
    )
}

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