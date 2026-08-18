package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"clean-architecture/Delivery/controllers"
	domain "clean-architecture/Domain"
	infrastructure "clean-architecture/Infrastructure"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================================
// MOCK JWT SERVICE
// ============================================================

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateAccessToken(
	userID string,
	role string,
) (string, error) {

	args := m.Called(userID, role)

	return args.String(0), args.Error(1)
}

func (m *MockJWTService) GenerateRefreshToken(
	userID string,
) (string, error) {

	args := m.Called(userID)

	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateAccessToken(
	tokenString string,
) (string, string, error) {

	args := m.Called(tokenString)

	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockJWTService) ValidateRefreshToken(
	tokenString string,
) (string, error) {

	args := m.Called(tokenString)

	return args.String(0), args.Error(1)
}

// Make sure the mock implements the interface.
var _ infrastructure.JWTService = (*MockJWTService)(nil)

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
// SETUP
// ============================================================

func createTestRouter() (
	*gin.Engine,
	*MockJWTService,
) {

	mockTask := new(MockTaskUsecase)
	mockUser := new(MockUserUsecase)
	mockJWT := new(MockJWTService)

	controller := controllers.NewController(
		mockTask,
		mockUser,
	)

	router := SetupRouter(
		controller,
		mockJWT,
	)

	return router, mockJWT
}

// ============================================================
// PUBLIC ROUTES
// ============================================================

func TestSetupRouter_RegisterRoute(t *testing.T) {

	router, _ := createTestRouter()

	request, _ := http.NewRequest(
		http.MethodPost,
		"/register",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	// The route exists.
	//
	// We don't expect 404.
	assert.NotEqual(
		t,
		http.StatusNotFound,
		recorder.Code,
	)
}

func TestSetupRouter_LoginRoute(t *testing.T) {

	router, _ := createTestRouter()

	request, _ := http.NewRequest(
		http.MethodPost,
		"/login",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.NotEqual(
		t,
		http.StatusNotFound,
		recorder.Code,
	)
}

// ============================================================
// PROTECTED ROUTES
// ============================================================

func TestSetupRouter_ProtectedRouteRequiresAuth(t *testing.T) {

	router, _ := createTestRouter()

	request, _ := http.NewRequest(
		http.MethodGet,
		"/tasks/tasks",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(
		t,
		http.StatusUnauthorized,
		recorder.Code,
	)
}

// ============================================================
// ADMIN ROUTES
// ============================================================

func TestSetupRouter_AdminRouteRequiresAuth(t *testing.T) {

	router, _ := createTestRouter()

	request, _ := http.NewRequest(
		http.MethodGet,
		"/admin/users",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(
		t,
		http.StatusUnauthorized,
		recorder.Code,
	)
}