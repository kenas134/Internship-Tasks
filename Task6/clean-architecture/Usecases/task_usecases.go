package usecases

import (
	"errors"

	domain "clean-architecture/Domain"
	"clean-architecture/Repositories"
)

type TaskUsecase struct {
	repository repositories.TaskRepository
}

func NewTaskUsecase(
	repository repositories.TaskRepository,
) *TaskUsecase {
	return &TaskUsecase{
		repository: repository,
	}
}

func (u *TaskUsecase) CreateTask(task domain.Task) (domain.Task,error) {

	if task.Title == "" {
		return domain.Task{},errors.New("title is required")
	}

	if task.Status == "" {
		task.Status = "pending"
	}

	newTask,err := u.repository.Create(&task)
	return *newTask,err
}



func (u *TaskUsecase) GetTasks() ([]domain.Task, error) {

	return u.repository.GetTasks()
}

func (u *TaskUsecase) GetTaskByID(
	id string,
) (domain.Task, error) {

	if id == "" {
		return domain.Task{}, errors.New("task ID is required")
	}

	return u.repository.GetTaskByID(id)
}

func (u *TaskUsecase) UpdateTask(
	id string,
	task domain.Task,
) error {

	if id == "" {
		return errors.New("task ID is required")
	}

	if task.Title == "" {
		return errors.New("title is required")
	}

	task.ID = id

	return u.repository.Update(task)
}

func (u *TaskUsecase) DeleteTask(id string) error {

	if id == "" {
		return errors.New("task ID is required")
	}

	return u.repository.Delete(id)
}