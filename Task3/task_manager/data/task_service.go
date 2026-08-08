package data

import (
	"task_manager/models"
	"time"
	"errors"
)

var tasks = []models.Task{

	{
		ID:          "1",
		Title:       "Clean the room",
		Description: "Organize books, clothes, and clean the floor",
		DueDate:     time.Date(2026, time.August, 10, 8, 0, 0, 0, time.UTC),
		Status:      "Pending",
	},

	{
		ID:          "2",
		Title:       "Wash dishes",
		Description: "Clean all dishes after breakfast and dinner",
		DueDate:     time.Date(2026, time.August, 10, 19, 0, 0, 0, time.UTC),
		Status:      "Completed",
	},

	{
		ID:          "3",
		Title:       "Buy groceries",
		Description: "Buy vegetables, fruits, milk, and other household items",
		DueDate:     time.Date(2026, time.August, 12, 15, 30, 0, 0, time.UTC),
		Status:      "Pending",
	},

	{
		ID:          "4",
		Title:       "Do laundry",
		Description: "Wash and fold dirty clothes",
		DueDate:     time.Date(2026, time.August, 14, 9, 0, 0, 0, time.UTC),
		Status:      "In Progress",
	},

	{
		ID:          "5",
		Title:       "Water plants",
		Description: "Water all plants in the house and garden",
		DueDate:     time.Date(2026, time.August, 16, 7, 30, 0, 0, time.UTC),
		Status:      "Pending",
	},
}


func GetTasks()[]models.Task{
	return tasks
}

func GetTaskByID(id string) (models.Task, error) {

    for _, task := range tasks {
        if task.ID == id {
            return task, nil
        }
    }

    return models.Task{}, errors.New("task not found")
}

func CreateTask(task models.Task) {
    tasks = append(tasks, task)
}

func DeleteTaskByID(id string) error {

    for i, task := range tasks {

        if task.ID == id {

            tasks = append(tasks[:i], tasks[i+1:]...)

            return nil
        }
    }

    return errors.New("task not found")
}

func UpdateTaskByID(id string, updatedTask models.Task) error {

    for i, task := range tasks {

        if task.ID == id {

            if updatedTask.Title != "" {
                tasks[i].Title = updatedTask.Title
            }

            if updatedTask.Description != "" {
                tasks[i].Description = updatedTask.Description
            }

            if updatedTask.Status != "" {
                tasks[i].Status = updatedTask.Status
            }

            if !updatedTask.DueDate.IsZero() {
                tasks[i].DueDate = updatedTask.DueDate
            }

            return nil
        }
    }

    return errors.New("task not found")
}

