package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	domain "clean-architecture/Domain"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================================
// MOCK TASK USECASE
// ============================================================

type MockTaskUsecase struct {
	mock.Mock
}

func (m *MockTaskUsecase) GetTasks() ([]domain.Task, error) {
	args := m.Called()

	return args.Get(0).([]domain.Task), args.Error(1)
}

func (m *MockTaskUsecase) GetTaskByID(
	id string,
) (domain.Task, error) {

	args := m.Called(id)

	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *MockTaskUsecase) CreateTask(
	task domain.Task,
) (domain.Task, error) {

	args := m.Called(task)

	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *MockTaskUsecase) UpdateTask(
	id string,
	task domain.Task,
) error {

	args := m.Called(id, task)

	return args.Error(0)
}

func (m *MockTaskUsecase) DeleteTask(
	id string,
) error {

	args := m.Called(id)

	return args.Error(0)
}

// ============================================================
// MOCK USER USECASE
// ============================================================

type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) Register(
	user domain.User,
) (domain.User, error) {

	args := m.Called(user)

	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserUsecase) Login(
	username string,
	password string,
) (string, string, error) {

	args := m.Called(username, password)

	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockUserUsecase) Profile(
	id string,
) (domain.User, error) {

	args := m.Called(id)

	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserUsecase) ListUsers() ([]domain.User, error) {

	args := m.Called()

	return args.Get(0).([]domain.User), args.Error(1)
}

// ============================================================
// HELPER
// ============================================================

func setupController() (
	*Controller,
	*MockTaskUsecase,
	*MockUserUsecase,
) {

	mockTask := new(MockTaskUsecase)
	mockUser := new(MockUserUsecase)

	controller := NewController(
		mockTask,
		mockUser,
	)

	return controller, mockTask, mockUser
}

func setupContext(
	method string,
	path string,
	body interface{},
) (*gin.Context, *httptest.ResponseRecorder) {

	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()

	var requestBody *bytes.Reader

	if body != nil {
		jsonData, _ := json.Marshal(body)
		requestBody = bytes.NewReader(jsonData)
	} else {
		requestBody = bytes.NewReader([]byte{})
	}

	ctx, _ := gin.CreateTestContext(recorder)

	ctx.Request = httptest.NewRequest(
		method,
		path,
		requestBody,
	)

	ctx.Request.Header.Set(
		"Content-Type",
		"application/json",
	)

	return ctx, recorder
}

// ============================================================
// REGISTER
// ============================================================

func TestRegister_Success(t *testing.T) {

	controller, _, mockUser := setupController()

	request := RegisterRequest{
		Firstname: "John",
		Lastname:  "Doe",
		Username:  "john",
		Password:  "password123",
		Role:      "user",
	}

	expectedUser := domain.User{
		ID:        "123",
		Firstname: "John",
		Lastname:  "Doe",
		Username:  "john",
		Password:  "",
		Role:      "user",
	}

	mockUser.
		On("Register", mock.AnythingOfType("domain.User")).
		Return(expectedUser, nil)

	ctx, recorder := setupContext(
		http.MethodPost,
		"/register",
		request,
	)

	controller.Register(ctx)

	assert.Equal(
		t,
		http.StatusCreated,
		recorder.Code,
	)

	mockUser.AssertExpectations(t)
}

func TestRegister_InvalidJSON(t *testing.T) {

	controller, _, _ := setupController()

	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	ctx.Request = httptest.NewRequest(
		http.MethodPost,
		"/register",
		bytes.NewBufferString(`invalid json`),
	)

	ctx.Request.Header.Set(
		"Content-Type",
		"application/json",
	)

	controller.Register(ctx)

	assert.Equal(
		t,
		http.StatusBadRequest,
		recorder.Code,
	)
}

func TestRegister_UsecaseError(t *testing.T) {

	controller, _, mockUser := setupController()

	request := RegisterRequest{
		Firstname: "John",
		Lastname:  "Doe",
		Username:  "john",
		Password:  "password123",
		Role:      "user",
	}

	mockUser.
		On("Register", mock.AnythingOfType("domain.User")).
		Return(
			domain.User{},
			errors.New("username already exists"),
		)

	ctx, recorder := setupContext(
		http.MethodPost,
		"/register",
		request,
	)

	controller.Register(ctx)

	assert.Equal(
		t,
		http.StatusBadRequest,
		recorder.Code,
	)

	mockUser.AssertExpectations(t)
}

// ============================================================
// LOGIN
// ============================================================

func TestLogin_Success(t *testing.T) {

	controller, _, mockUser := setupController()

	request := LoginRequest{
		Username: "john",
		Password: "password123",
	}

	mockUser.
		On(
			"Login",
			"john",
			"password123",
		).
		Return(
			"access-token",
			"refresh-token",
			nil,
		)

	ctx, recorder := setupContext(
		http.MethodPost,
		"/login",
		request,
	)

	controller.Login(ctx)

	assert.Equal(
		t,
		http.StatusOK,
		recorder.Code,
	)

	mockUser.AssertExpectations(t)
}

func TestLogin_InvalidJSON(t *testing.T) {

	controller, _, _ := setupController()

	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	ctx.Request = httptest.NewRequest(
		http.MethodPost,
		"/login",
		bytes.NewBufferString(`invalid json`),
	)

	ctx.Request.Header.Set(
		"Content-Type",
		"application/json",
	)

	controller.Login(ctx)

	assert.Equal(
		t,
		http.StatusBadRequest,
		recorder.Code,
	)
}

func TestLogin_InvalidCredentials(t *testing.T) {

	controller, _, mockUser := setupController()

	request := LoginRequest{
		Username: "john",
		Password: "wrongpassword",
	}

	mockUser.
		On(
			"Login",
			"john",
			"wrongpassword",
		).
		Return(
			"",
			"",
			errors.New("invalid username or password"),
		)

	ctx, recorder := setupContext(
		http.MethodPost,
		"/login",
		request,
	)

	controller.Login(ctx)

	assert.Equal(
		t,
		http.StatusUnauthorized,
		recorder.Code,
	)

	mockUser.AssertExpectations(t)
}

// ============================================================
// PROFILE
// ============================================================

func TestProfile_Success(t *testing.T) {

	controller, _, mockUser := setupController()

	expectedUser := domain.User{
		ID:       "123",
		Username: "john",
		Role:     "user",
	}

	mockUser.
		On("Profile", "123").
		Return(expectedUser, nil)

	ctx, recorder := setupContext(
		http.MethodGet,
		"/profile",
		nil,
	)

	ctx.Set("userID", "123")

	controller.Profile(ctx)

	assert.Equal(
		t,
		http.StatusOK,
		recorder.Code,
	)

	mockUser.AssertExpectations(t)
}

func TestProfile_NoUserID(t *testing.T) {

	controller, _, _ := setupController()

	ctx, recorder := setupContext(
		http.MethodGet,
		"/profile",
		nil,
	)

	controller.Profile(ctx)

	assert.Equal(
		t,
		http.StatusUnauthorized,
		recorder.Code,
	)
}

func TestProfile_InvalidUserIDType(t *testing.T) {

	controller, _, _ := setupController()

	ctx, recorder := setupContext(
		http.MethodGet,
		"/profile",
		nil,
	)

	ctx.Set("userID", 123)

	controller.Profile(ctx)

	assert.Equal(
		t,
		http.StatusUnauthorized,
		recorder.Code,
	)
}

func TestProfile_UserNotFound(t *testing.T) {

	controller, _, mockUser := setupController()

	mockUser.
		On("Profile", "123").
		Return(
			domain.User{},
			errors.New("user not found"),
		)

	ctx, recorder := setupContext(
		http.MethodGet,
		"/profile",
		nil,
	)

	ctx.Set("userID", "123")

	controller.Profile(ctx)

	assert.Equal(
		t,
		http.StatusNotFound,
		recorder.Code,
	)

	mockUser.AssertExpectations(t)
}

// ============================================================
// LIST USERS
// ============================================================

func TestListUsers_Success(t *testing.T) {

	controller, _, mockUser := setupController()

	users := []domain.User{
		{
			ID:       "1",
			Username: "john",
			Role:     "user",
		},
		{
			ID:       "2",
			Username: "admin",
			Role:     "admin",
		},
	}

	mockUser.
		On("ListUsers").
		Return(users, nil)

	ctx, recorder := setupContext(
		http.MethodGet,
		"/users",
		nil,
	)

	controller.ListUsers(ctx)

	assert.Equal(
		t,
		http.StatusOK,
		recorder.Code,
	)

	mockUser.AssertExpectations(t)
}

func TestListUsers_Error(t *testing.T) {

	controller, _, mockUser := setupController()

	mockUser.
		On("ListUsers").
		Return(
			[]domain.User{},
			errors.New("database error"),
		)

	ctx, recorder := setupContext(
		http.MethodGet,
		"/users",
		nil,
	)

	controller.ListUsers(ctx)

	assert.Equal(
		t,
		http.StatusInternalServerError,
		recorder.Code,
	)

	mockUser.AssertExpectations(t)
}

// ============================================================
// GET TASKS
// ============================================================

func TestGetTasks_Success(t *testing.T) {

	controller, mockTask, _ := setupController()

	tasks := []domain.Task{
		{
			ID:          "1",
			Title:       "Learn Go",
			Description: "Learn testing",
			Status:      "pending",
		},
		{
			ID:          "2",
			Title:       "Build API",
			Description: "Build REST API",
			Status:      "completed",
		},
	}

	mockTask.
		On("GetTasks").
		Return(tasks, nil)

	ctx, recorder := setupContext(
		http.MethodGet,
		"/tasks",
		nil,
	)

	controller.GetTasks(ctx)

	assert.Equal(
		t,
		http.StatusOK,
		recorder.Code,
	)

	mockTask.AssertExpectations(t)
}

func TestGetTasks_Error(t *testing.T) {

	controller, mockTask, _ := setupController()

	mockTask.
		On("GetTasks").
		Return(
			[]domain.Task{},
			errors.New("database error"),
		)

	ctx, recorder := setupContext(
		http.MethodGet,
		"/tasks",
		nil,
	)

	controller.GetTasks(ctx)

	assert.Equal(
		t,
		http.StatusInternalServerError,
		recorder.Code,
	)

	mockTask.AssertExpectations(t)
}

// ============================================================
// GET TASK BY ID
// ============================================================

func TestGetTaskByID_Success(t *testing.T) {

	controller, mockTask, _ := setupController()

	task := domain.Task{
		ID:     "123",
		Title:  "Learn Go",
		Status: "pending",
	}

	mockTask.
		On("GetTaskByID", "123").
		Return(task, nil)

	ctx, recorder := setupContext(
		http.MethodGet,
		"/tasks/123",
		nil,
	)

	ctx.Params = gin.Params{
		{
			Key:   "id",
			Value: "123",
		},
	}

	controller.GetTaskByID(ctx)

	assert.Equal(
		t,
		http.StatusOK,
		recorder.Code,
	)

	mockTask.AssertExpectations(t)
}

func TestGetTaskByID_NotFound(t *testing.T) {

	controller, mockTask, _ := setupController()

	mockTask.
		On("GetTaskByID", "123").
		Return(
			domain.Task{},
			errors.New("task not found"),
		)

	ctx, recorder := setupContext(
		http.MethodGet,
		"/tasks/123",
		nil,
	)

	ctx.Params = gin.Params{
		{
			Key:   "id",
			Value: "123",
		},
	}

	controller.GetTaskByID(ctx)

	assert.Equal(
		t,
		http.StatusNotFound,
		recorder.Code,
	)

	mockTask.AssertExpectations(t)
}

// ============================================================
// CREATE TASK
// ============================================================

func TestCreateTask_Success(t *testing.T) {

	controller, mockTask, _ := setupController()

	dueDate := time.Now()

	request := TaskRequest{
		Title:       "Learn Go",
		Description: "Learn testing",
		DueDate:     dueDate,
		Status:      "pending",
	}

	expectedTask := domain.Task{
		ID:          "123",
		Title:       "Learn Go",
		Description: "Learn testing",
		DueDate:     dueDate,
		Status:      "pending",
	}

	mockTask.
		On(
			"CreateTask",
			mock.AnythingOfType("domain.Task"),
		).
		Return(expectedTask, nil)

	ctx, recorder := setupContext(
		http.MethodPost,
		"/tasks",
		request,
	)

	controller.CreateTask(ctx)

	assert.Equal(
		t,
		http.StatusCreated,
		recorder.Code,
	)

	mockTask.AssertExpectations(t)
}

func TestCreateTask_InvalidJSON(t *testing.T) {

	controller, _, _ := setupController()

	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	ctx.Request = httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		bytes.NewBufferString(`invalid json`),
	)

	ctx.Request.Header.Set(
		"Content-Type",
		"application/json",
	)

	controller.CreateTask(ctx)

	assert.Equal(
		t,
		http.StatusBadRequest,
		recorder.Code,
	)
}

func TestCreateTask_UsecaseError(t *testing.T) {

	controller, mockTask, _ := setupController()

	request := TaskRequest{
		Title:  "",
		Status: "pending",
	}

	mockTask.
		On(
			"CreateTask",
			mock.AnythingOfType("domain.Task"),
		).
		Return(
			domain.Task{},
			errors.New("title is required"),
		)

	ctx, recorder := setupContext(
		http.MethodPost,
		"/tasks",
		request,
	)

	controller.CreateTask(ctx)

	assert.Equal(
		t,
		http.StatusBadRequest,
		recorder.Code,
	)

	mockTask.AssertExpectations(t)
}

// ============================================================
// UPDATE TASK
// ============================================================

func TestUpdateTask_Success(t *testing.T) {

	controller, mockTask, _ := setupController()

	request := TaskRequest{
		ID:          "123",
		Title:       "Updated Task",
		Description: "Updated description",
		Status:      "completed",
	}

	mockTask.
		On(
			"UpdateTask",
			"123",
			mock.AnythingOfType("domain.Task"),
		).
		Return(nil)

	ctx, recorder := setupContext(
		http.MethodPut,
		"/tasks/123",
		request,
	)

	controller.UpdateTask(ctx)

	assert.Equal(
		t,
		http.StatusCreated,
		recorder.Code,
	)

	mockTask.AssertExpectations(t)
}

func TestUpdateTask_InvalidJSON(t *testing.T) {

	controller, _, _ := setupController()

	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	ctx.Request = httptest.NewRequest(
		http.MethodPut,
		"/tasks/123",
		bytes.NewBufferString(`invalid json`),
	)

	ctx.Request.Header.Set(
		"Content-Type",
		"application/json",
	)

	controller.UpdateTask(ctx)

	assert.Equal(
		t,
		http.StatusBadRequest,
		recorder.Code,
	)
}

func TestUpdateTask_Error(t *testing.T) {

	controller, mockTask, _ := setupController()

	request := TaskRequest{
		ID:    "123",
		Title: "Updated Task",
	}

	mockTask.
		On(
			"UpdateTask",
			"123",
			mock.AnythingOfType("domain.Task"),
		).
		Return(errors.New("task not found"))

	ctx, recorder := setupContext(
		http.MethodPut,
		"/tasks/123",
		request,
	)

	controller.UpdateTask(ctx)

	assert.Equal(
		t,
		http.StatusBadRequest,
		recorder.Code,
	)

	mockTask.AssertExpectations(t)
}

// ============================================================
// DELETE TASK
// ============================================================

func TestDeleteTask_Success(t *testing.T) {

	controller, mockTask, _ := setupController()

	mockTask.
		On("DeleteTask", "123").
		Return(nil)

	ctx, recorder := setupContext(
		http.MethodDelete,
		"/tasks/123",
		nil,
	)

	ctx.Params = gin.Params{
		{
			Key:   "id",
			Value: "123",
		},
	}

	controller.DeleteTask(ctx)

	assert.Equal(
		t,
		http.StatusOK,
		recorder.Code,
	)

	mockTask.AssertExpectations(t)
}

func TestDeleteTask_NotFound(t *testing.T) {

	controller, mockTask, _ := setupController()

	mockTask.
		On("DeleteTask", "123").
		Return(errors.New("task not found"))

	ctx, recorder := setupContext(
		http.MethodDelete,
		"/tasks/123",
		nil,
	)

	ctx.Params = gin.Params{
		{
			Key:   "id",
			Value: "123",
		},
	}

	controller.DeleteTask(ctx)

	assert.Equal(
		t,
		http.StatusNotFound,
		recorder.Code,
	)

	mockTask.AssertExpectations(t)
}