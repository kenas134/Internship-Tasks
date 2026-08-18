package usecases

import (
	"errors"
	"testing"

	domain "clean-architecture/Domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================================
// MOCK USER REPOSITORY
// ============================================================

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserByUsername(
	username string,
) (domain.User, error) {

	args := m.Called(username)

	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepository) CreateUser(
	user *domain.User,
) (domain.User, error) {

	args := m.Called(user)

	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(
	id string,
) (domain.User, error) {

	args := m.Called(id)

	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepository) GetUsers() ([]domain.User, error) {

	args := m.Called()

	return args.Get(0).([]domain.User), args.Error(1)
}

// ============================================================
// MOCK PASSWORD SERVICE
// ============================================================

type MockPasswordService struct {
	mock.Mock
}

func (m *MockPasswordService) Hash(
	password string,
) (string, error) {

	args := m.Called(password)

	return args.String(0), args.Error(1)
}

func (m *MockPasswordService) ComparePassword(
	hashedPassword string,
	password string,
) error {

	args := m.Called(hashedPassword, password)

	return args.Error(0)
}

// ============================================================
// MOCK JWT SERVICE
// ============================================================

type MockJWTService struct {
	mock.Mock
}

// GenerateAccessToken is used by Login.
func (m *MockJWTService) GenerateAccessToken(
	userID string,
	role string,
) (string, error) {

	args := m.Called(userID, role)

	return args.String(0), args.Error(1)
}

// GenerateRefreshToken is used by Login.
func (m *MockJWTService) GenerateRefreshToken(
	userID string,
) (string, error) {

	args := m.Called(userID)

	return args.String(0), args.Error(1)
}

// These two are required because they are part
// of the JWTService interface, even though the
// current UserUsecase does not use them.

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

// ============================================================
// REGISTER
// ============================================================

func TestRegister(t *testing.T) {

	t.Run("successful registration", func(t *testing.T) {

		// -------------------------
		// Arrange
		// -------------------------

		mockRepo := new(MockUserRepository)
		mockPassword := new(MockPasswordService)
		mockJWT := new(MockJWTService)

		usecase := NewUserUsecase(
			mockRepo,
			mockPassword,
			mockJWT,
		)

		user := domain.User{
			Username: "john",
			Password: "password123",
			Role:     "user",
		}

		expectedUser := domain.User{
			ID:       "123",
			Username: "john",
			Password: "",
			Role:     "user",
		}

		// User does not exist.
		mockRepo.
			On("GetUserByUsername", "john").
			Return(
				domain.User{},
				errors.New("user not found"),
			)

		// Password is successfully hashed.
		mockPassword.
			On("Hash", "password123").
			Return("hashedPassword", nil)

		// User is successfully created.
		mockRepo.
			On(
				"CreateUser",
				mock.AnythingOfType("*domain.User"),
			).
			Return(expectedUser, nil)

		// -------------------------
		// Act
		// -------------------------

		result, err := usecase.Register(user)

		// -------------------------
		// Assert
		// -------------------------

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)

		mockRepo.AssertExpectations(t)
		mockPassword.AssertExpectations(t)
	})

	t.Run("username is required", func(t *testing.T) {

		mockRepo := new(MockUserRepository)
		mockPassword := new(MockPasswordService)
		mockJWT := new(MockJWTService)

		usecase := NewUserUsecase(
			mockRepo,
			mockPassword,
			mockJWT,
		)

		user := domain.User{
			Username: "",
			Password: "password123",
		}

		result, err := usecase.Register(user)

		assert.Error(t, err)
		assert.Equal(
			t,
			"username is required",
			err.Error(),
		)

		assert.Equal(t, domain.User{}, result)

		// Because validation failed, repository
		// and password service should never be called.
		mockRepo.AssertNotCalled(
			t,
			"GetUserByUsername",
		)

		mockPassword.AssertNotCalled(
			t,
			"Hash",
		)
	})

	t.Run("password is too short", func(t *testing.T) {

		mockRepo := new(MockUserRepository)
		mockPassword := new(MockPasswordService)
		mockJWT := new(MockJWTService)

		usecase := NewUserUsecase(
			mockRepo,
			mockPassword,
			mockJWT,
		)

		user := domain.User{
			Username: "john",
			Password: "123",
		}

		result, err := usecase.Register(user)

		assert.Error(t, err)
		assert.Equal(
			t,
			"password must be at least 8 characters",
			err.Error(),
		)

		assert.Equal(t, domain.User{}, result)

		mockRepo.AssertNotCalled(
			t,
			"GetUserByUsername",
		)

		mockPassword.AssertNotCalled(
			t,
			"Hash",
		)
	})

	t.Run("username already exists", func(t *testing.T) {

		mockRepo := new(MockUserRepository)
		mockPassword := new(MockPasswordService)
		mockJWT := new(MockJWTService)

		usecase := NewUserUsecase(
			mockRepo,
			mockPassword,
			mockJWT,
		)

		existingUser := domain.User{
			ID:       "123",
			Username: "john",
			Role:     "user",
		}

		user := domain.User{
			Username: "john",
			Password: "password123",
		}

		mockRepo.
			On("GetUserByUsername", "john").
			Return(existingUser, nil)

		result, err := usecase.Register(user)

		assert.Error(t, err)
		assert.Equal(
			t,
			"username already exists",
			err.Error(),
		)

		assert.Equal(t, domain.User{}, result)

		mockPassword.AssertNotCalled(
			t,
			"Hash",
		)

		mockRepo.AssertNotCalled(
			t,
			"CreateUser",
		)
	})

	t.Run("password hashing fails", func(t *testing.T) {

		mockRepo := new(MockUserRepository)
		mockPassword := new(MockPasswordService)
		mockJWT := new(MockJWTService)

		usecase := NewUserUsecase(
			mockRepo,
			mockPassword,
			mockJWT,
		)

		user := domain.User{
			Username: "john",
			Password: "password123",
		}

		hashError := errors.New("hashing failed")

		mockRepo.
			On("GetUserByUsername", "john").
			Return(
				domain.User{},
				errors.New("user not found"),
			)

		mockPassword.
			On("Hash", "password123").
			Return("", hashError)

		result, err := usecase.Register(user)

		assert.Error(t, err)
		assert.Equal(t, hashError, err)
		assert.Equal(t, domain.User{}, result)

		mockRepo.AssertNotCalled(
			t,
			"CreateUser",
		)
	})

	t.Run("repository create fails", func(t *testing.T) {

		mockRepo := new(MockUserRepository)
		mockPassword := new(MockPasswordService)
		mockJWT := new(MockJWTService)

		usecase := NewUserUsecase(
			mockRepo,
			mockPassword,
			mockJWT,
		)

		user := domain.User{
			Username: "john",
			Password: "password123",
		}

		createError := errors.New("database error")

		mockRepo.
			On("GetUserByUsername", "john").
			Return(
				domain.User{},
				errors.New("user not found"),
			)

		mockPassword.
			On("Hash", "password123").
			Return("hashedPassword", nil)

		mockRepo.
			On(
				"CreateUser",
				mock.AnythingOfType("*domain.User"),
			).
			Return(domain.User{}, createError)

		result, err := usecase.Register(user)

		assert.Error(t, err)
		assert.Equal(t, createError, err)
		assert.Equal(t, domain.User{}, result)
	})
}

// ============================================================
// LOGIN
// ============================================================

func TestLogin(t *testing.T) {

	t.Run("successful login", func(t *testing.T) {

		mockRepo := new(MockUserRepository)
		mockPassword := new(MockPasswordService)
		mockJWT := new(MockJWTService)

		usecase := NewUserUsecase(
			mockRepo,
			mockPassword,
			mockJWT,
		)

		user := domain.User{
			ID:       "123",
			Username: "john",
			Password: "hashedPassword",
			Role:     "user",
		}

		mockRepo.
			On("GetUserByUsername", "john").
			Return(user, nil)

		mockPassword.
			On(
				"ComparePassword",
				"hashedPassword",
				"password123",
			).
			Return(nil)

		mockJWT.
			On(
				"GenerateAccessToken",
				"123",
				"user",
			).
			Return("access-token", nil)

		mockJWT.
			On(
				"GenerateRefreshToken",
				"123",
			).
			Return("refresh-token", nil)

		accessToken, refreshToken, err :=
			usecase.Login(
				"john",
				"password123",
			)

		assert.NoError(t, err)
		assert.Equal(t, "access-token", accessToken)
		assert.Equal(t, "refresh-token", refreshToken)

		mockRepo.AssertExpectations(t)
		mockPassword.AssertExpectations(t)
		mockJWT.AssertExpectations(t)
	})

	t.Run("user does not exist", func(t *testing.T) {

		mockRepo := new(MockUserRepository)
		mockPassword := new(MockPasswordService)
		mockJWT := new(MockJWTService)

		usecase := NewUserUsecase(
			mockRepo,
			mockPassword,
			mockJWT,
		)

		mockRepo.
			On("GetUserByUsername", "john").
			Return(
				domain.User{},
				errors.New("user not found"),
			)

		accessToken, refreshToken, err :=
			usecase.Login(
				"john",
				"password123",
			)

		assert.Error(t, err)
		assert.Equal(
			t,
			"invalid username or password",
			err.Error(),
		)

		assert.Equal(t, "", accessToken)
		assert.Equal(t, "", refreshToken)

		mockPassword.AssertNotCalled(
			t,
			"ComparePassword",
		)

		mockJWT.AssertNotCalled(
			t,
			"GenerateAccessToken",
		)
	})

	t.Run("wrong password", func(t *testing.T) {

		mockRepo := new(MockUserRepository)
		mockPassword := new(MockPasswordService)
		mockJWT := new(MockJWTService)

		usecase := NewUserUsecase(
			mockRepo,
			mockPassword,
			mockJWT,
		)

		user := domain.User{
			ID:       "123",
			Username: "john",
			Password: "hashedPassword",
			Role:     "user",
		}

		mockRepo.
			On("GetUserByUsername", "john").
			Return(user, nil)

		mockPassword.
			On(
				"ComparePassword",
				"hashedPassword",
				"wrongPassword",
			).
			Return(errors.New("wrong password"))

		accessToken, refreshToken, err :=
			usecase.Login(
				"john",
				"wrongPassword",
			)

		assert.Error(t, err)
		assert.Equal(
			t,
			"invalid username or password",
			err.Error(),
		)

		assert.Equal(t, "", accessToken)
		assert.Equal(t, "", refreshToken)

		mockJWT.AssertNotCalled(
			t,
			"GenerateAccessToken",
		)

		mockJWT.AssertNotCalled(
			t,
			"GenerateRefreshToken",
		)
	})

	t.Run("access token generation fails", func(t *testing.T) {

		mockRepo := new(MockUserRepository)
		mockPassword := new(MockPasswordService)
		mockJWT := new(MockJWTService)

		usecase := NewUserUsecase(
			mockRepo,
			mockPassword,
			mockJWT,
		)

		user := domain.User{
			ID:       "123",
			Username: "john",
			Password: "hashedPassword",
			Role:     "user",
		}

		tokenError := errors.New("access token failed")

		mockRepo.
			On("GetUserByUsername", "john").
			Return(user, nil)

		mockPassword.
			On(
				"ComparePassword",
				"hashedPassword",
				"password123",
			).
			Return(nil)

		mockJWT.
			On(
				"GenerateAccessToken",
				"123",
				"user",
			).
			Return("", tokenError)

		accessToken, refreshToken, err :=
			usecase.Login(
				"john",
				"password123",
			)

		assert.Error(t, err)
		assert.Equal(t, tokenError, err)

		assert.Equal(t, "", accessToken)
		assert.Equal(t, "", refreshToken)

		mockJWT.AssertNotCalled(
			t,
			"GenerateRefreshToken",
		)
	})

	t.Run("refresh token generation fails", func(t *testing.T) {

		mockRepo := new(MockUserRepository)
		mockPassword := new(MockPasswordService)
		mockJWT := new(MockJWTService)

		usecase := NewUserUsecase(
			mockRepo,
			mockPassword,
			mockJWT,
		)

		user := domain.User{
			ID:       "123",
			Username: "john",
			Password: "hashedPassword",
			Role:     "user",
		}

		tokenError := errors.New("refresh token failed")

		mockRepo.
			On("GetUserByUsername", "john").
			Return(user, nil)

		mockPassword.
			On(
				"ComparePassword",
				"hashedPassword",
				"password123",
			).
			Return(nil)

		mockJWT.
			On(
				"GenerateAccessToken",
				"123",
				"user",
			).
			Return("access-token", nil)

		mockJWT.
			On(
				"GenerateRefreshToken",
				"123",
			).
			Return("", tokenError)

		accessToken, refreshToken, err :=
			usecase.Login(
				"john",
				"password123",
			)

		assert.Error(t, err)
		assert.Equal(t, tokenError, err)

		assert.Equal(t, "", accessToken)
		assert.Equal(t, "", refreshToken)
	})
}

// ============================================================
// PROFILE
// ============================================================

func TestProfile(t *testing.T) {

	t.Run("successful profile", func(t *testing.T) {

		mockRepo := new(MockUserRepository)
		mockPassword := new(MockPasswordService)
		mockJWT := new(MockJWTService)

		usecase := NewUserUsecase(
			mockRepo,
			mockPassword,
			mockJWT,
		)

		expectedUser := domain.User{
			ID:       "123",
			Username: "john",
			Role:     "user",
		}

		mockRepo.
			On("GetUserByID", "123").
			Return(expectedUser, nil)

		result, err := usecase.Profile("123")

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("repository returns error", func(t *testing.T) {

		mockRepo := new(MockUserRepository)
		mockPassword := new(MockPasswordService)
		mockJWT := new(MockJWTService)

		usecase := NewUserUsecase(
			mockRepo,
			mockPassword,
			mockJWT,
		)

		repositoryError := errors.New("user not found")

		mockRepo.
			On("GetUserByID", "123").
			Return(domain.User{}, repositoryError)

		result, err := usecase.Profile("123")

		assert.Error(t, err)
		assert.Equal(t, repositoryError, err)
		assert.Equal(t, domain.User{}, result)

		mockRepo.AssertExpectations(t)
	})
}

// ============================================================
// LIST USERS
// ============================================================

func TestListUsers(t *testing.T) {

	t.Run("successful list", func(t *testing.T) {

		mockRepo := new(MockUserRepository)
		mockPassword := new(MockPasswordService)
		mockJWT := new(MockJWTService)

		usecase := NewUserUsecase(
			mockRepo,
			mockPassword,
			mockJWT,
		)

		expectedUsers := []domain.User{
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

		mockRepo.
			On("GetUsers").
			Return(expectedUsers, nil)

		result, err := usecase.ListUsers()

		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("repository returns error", func(t *testing.T) {

		mockRepo := new(MockUserRepository)
		mockPassword := new(MockPasswordService)
		mockJWT := new(MockJWTService)

		usecase := NewUserUsecase(
			mockRepo,
			mockPassword,
			mockJWT,
		)

		repositoryError := errors.New("database error")

		mockRepo.
			On("GetUsers").
			Return([]domain.User{}, repositoryError)

		result, err := usecase.ListUsers()

		assert.Error(t, err)
		assert.Equal(t, repositoryError, err)
		assert.Empty(t, result)

		mockRepo.AssertExpectations(t)
	})
}