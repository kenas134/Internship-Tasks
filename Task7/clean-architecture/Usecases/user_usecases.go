package usecases

import (
	domain "clean-architecture/Domain"
	infrastructure "clean-architecture/Infrastructure"
	repositories "clean-architecture/Repositories"
	"errors"
)

type UserUsecase struct {
	userRepository  repositories.UserRepository
	passwordService infrastructure.PasswordService
	jwtService      infrastructure.JWTService
}

func NewUserUsecase(
	userRepository repositories.UserRepository,
	passwordService infrastructure.PasswordService,
	jwtService infrastructure.JWTService,

) *UserUsecase {
	return &UserUsecase{
		userRepository:  userRepository,
		passwordService: passwordService,
		jwtService:      jwtService,
	}
}

func (u *UserUsecase) Register(user domain.User) (domain.User, error) {

	if user.Username == "" {
		return domain.User{}, errors.New("username is required")
	}

	if len(user.Password) < 8 {
		return domain.User{}, errors.New("password must be at least 8 characters")
	}

	// Check if username already exists
	_, err := u.userRepository.GetUserByUsername(user.Username)

	if err == nil {
		return domain.User{}, errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := u.passwordService.Hash(user.Password)
	if err != nil {
		return domain.User{}, err
	}

	user.Password = hashedPassword

	newUser, err := u.userRepository.CreateUser(&user)
	if err != nil {
		return domain.User{}, err
	}

	return newUser, nil
}

func (u *UserUsecase) Login(username string, password string) (string, string, error) {
	user, err := u.userRepository.GetUserByUsername(username)

	if err != nil {
		return "", "", errors.New(
			"invalid username or password",
		)
	}

	err = u.passwordService.ComparePassword(
		user.Password,
		password,
	)

	if err != nil {
		return "", "", errors.New(
			"invalid username or password",
		)
	}

	accessToken, err := u.jwtService.GenerateAccessToken(
		user.ID,
		user.Role,
	)

	if err != nil {
		return "", "", err
	}

	refreshToken, err := u.jwtService.GenerateRefreshToken(
		user.ID,
	)

	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (u *UserUsecase) Profile(id string) (domain.User, error) {
	return u.userRepository.GetUserByID(id)
}

func (u *UserUsecase) ListUsers() ([]domain.User, error) {
	return u.userRepository.GetUsers()
}
