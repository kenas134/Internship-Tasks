package usecases

import (
	"errors"
	"testing"

	domain "clean-architecture/Domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ----------------------------------------------------
// MOCK REPOSITORY
// ----------------------------------------------------

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(task *domain.Task) (*domain.Task, error) {
	args := m.Called(task)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *MockTaskRepository) GetTasks() ([]domain.Task, error) {
	args := m.Called()

	return args.Get(0).([]domain.Task), args.Error(1)
}

func (m *MockTaskRepository) GetTaskByID(id string) (domain.Task, error) {
	args := m.Called(id)

	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *MockTaskRepository) Update(task domain.Task) error {
	args := m.Called(task)

	return args.Error(0)
}

func (m *MockTaskRepository) Delete(id string) error {
	args := m.Called(id)

	return args.Error(0)
}

// ----------------------------------------------------
// CREATE TASK TESTS
// ----------------------------------------------------

func TestCreateTask(t *testing.T) {

	t.Run("successful creation", func(t *testing.T) {

		mockRepo := new(MockTaskRepository)

		usecase := NewTaskUsecase(mockRepo)

		task := domain.Task{
			Title:  "Learn Go",
			Status: "",
		}

		expectedTask := domain.Task{
			ID:     "123",
			Title:  "Learn Go",
			Status: "pending",
		}

		mockRepo.
			On("Create", mock.AnythingOfType("*domain.Task")).
			Return(&expectedTask, nil)

		result, err := usecase.CreateTask(task)

		assert.NoError(t, err)
		assert.Equal(t, expectedTask, result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("title is required", func(t *testing.T) {

		mockRepo := new(MockTaskRepository)

		usecase := NewTaskUsecase(mockRepo)

		task := domain.Task{
			Title: "",
		}

		result, err := usecase.CreateTask(task)

		assert.Error(t, err)
		assert.Equal(t, "title is required", err.Error())
		assert.Equal(t, domain.Task{}, result)

		// Repository should NOT be called
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("repository returns error", func(t *testing.T) {

		mockRepo := new(MockTaskRepository)

		usecase := NewTaskUsecase(mockRepo)

		task := domain.Task{
			Title: "Learn Go",
		}

		repositoryError := errors.New("database error")

		mockRepo.
			On("Create", mock.AnythingOfType("*domain.Task")).
			Return(nil, repositoryError)

		result, err := usecase.CreateTask(task)

		assert.Error(t, err)
		assert.Equal(t, repositoryError, err)
		assert.Equal(t, domain.Task{}, result)

		mockRepo.AssertExpectations(t)
	})
}

// ----------------------------------------------------
// GET TASKS TEST
// ----------------------------------------------------

func TestGetTasks(t *testing.T) {

	mockRepo := new(MockTaskRepository)

	usecase := NewTaskUsecase(mockRepo)

	expectedTasks := []domain.Task{
		{
			ID:     "1",
			Title:  "Learn Go",
			Status: "pending",
		},
		{
			ID:     "2",
			Title:  "Build API",
			Status: "completed",
		},
	}

	mockRepo.
		On("GetTasks").
		Return(expectedTasks, nil)

	result, err := usecase.GetTasks()

	assert.NoError(t, err)
	assert.Equal(t, expectedTasks, result)

	mockRepo.AssertExpectations(t)
}

// ----------------------------------------------------
// GET TASK BY ID TESTS
// ----------------------------------------------------

func TestGetTaskByID(t *testing.T) {

	t.Run("successful get", func(t *testing.T) {

		mockRepo := new(MockTaskRepository)

		usecase := NewTaskUsecase(mockRepo)

		expectedTask := domain.Task{
			ID:     "123",
			Title:  "Learn Go",
			Status: "pending",
		}

		mockRepo.
			On("GetTaskByID", "123").
			Return(expectedTask, nil)

		result, err := usecase.GetTaskByID("123")

		assert.NoError(t, err)
		assert.Equal(t, expectedTask, result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ID is empty", func(t *testing.T) {

		mockRepo := new(MockTaskRepository)

		usecase := NewTaskUsecase(mockRepo)

		result, err := usecase.GetTaskByID("")

		assert.Error(t, err)
		assert.Equal(t, "task ID is required", err.Error())
		assert.Equal(t, domain.Task{}, result)

		mockRepo.AssertNotCalled(t, "GetTaskByID")
	})

	t.Run("repository returns error", func(t *testing.T) {

		mockRepo := new(MockTaskRepository)

		usecase := NewTaskUsecase(mockRepo)

		repositoryError := errors.New("task not found")

		mockRepo.
			On("GetTaskByID", "123").
			Return(domain.Task{}, repositoryError)

		result, err := usecase.GetTaskByID("123")

		assert.Error(t, err)
		assert.Equal(t, repositoryError, err)
		assert.Equal(t, domain.Task{}, result)

		mockRepo.AssertExpectations(t)
	})
}

// ----------------------------------------------------
// UPDATE TASK TESTS
// ----------------------------------------------------

func TestUpdateTask(t *testing.T) {

	t.Run("successful update", func(t *testing.T) {

		mockRepo := new(MockTaskRepository)

		usecase := NewTaskUsecase(mockRepo)

		task := domain.Task{
			Title:  "Updated task",
			Status: "completed",
		}

		expectedTask := task
		expectedTask.ID = "123"

		mockRepo.
			On("Update", expectedTask).
			Return(nil)

		err := usecase.UpdateTask("123", task)

		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ID is empty", func(t *testing.T) {

		mockRepo := new(MockTaskRepository)

		usecase := NewTaskUsecase(mockRepo)

		task := domain.Task{
			Title: "Updated task",
		}

		err := usecase.UpdateTask("", task)

		assert.Error(t, err)
		assert.Equal(t, "task ID is required", err.Error())

		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("title is empty", func(t *testing.T) {

		mockRepo := new(MockTaskRepository)

		usecase := NewTaskUsecase(mockRepo)

		task := domain.Task{
			Title: "",
		}

		err := usecase.UpdateTask("123", task)

		assert.Error(t, err)
		assert.Equal(t, "title is required", err.Error())

		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("repository returns error", func(t *testing.T) {

		mockRepo := new(MockTaskRepository)

		usecase := NewTaskUsecase(mockRepo)

		task := domain.Task{
			Title:  "Updated task",
			Status: "completed",
		}

		expectedTask := task
		expectedTask.ID = "123"

		repositoryError := errors.New("database error")

		mockRepo.
			On("Update", expectedTask).
			Return(repositoryError)

		err := usecase.UpdateTask("123", task)

		assert.Error(t, err)
		assert.Equal(t, repositoryError, err)

		mockRepo.AssertExpectations(t)
	})
}

// ----------------------------------------------------
// DELETE TASK TESTS
// ----------------------------------------------------

func TestDeleteTask(t *testing.T) {

	t.Run("successful delete", func(t *testing.T) {

		mockRepo := new(MockTaskRepository)

		usecase := NewTaskUsecase(mockRepo)

		mockRepo.
			On("Delete", "123").
			Return(nil)

		err := usecase.DeleteTask("123")

		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ID is empty", func(t *testing.T) {

		mockRepo := new(MockTaskRepository)

		usecase := NewTaskUsecase(mockRepo)

		err := usecase.DeleteTask("")

		assert.Error(t, err)
		assert.Equal(t, "task ID is required", err.Error())

		mockRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("repository returns error", func(t *testing.T) {

		mockRepo := new(MockTaskRepository)

		usecase := NewTaskUsecase(mockRepo)

		repositoryError := errors.New("database error")

		mockRepo.
			On("Delete", "123").
			Return(repositoryError)

		err := usecase.DeleteTask("123")

		assert.Error(t, err)
		assert.Equal(t, repositoryError, err)

		mockRepo.AssertExpectations(t)
	})
}
