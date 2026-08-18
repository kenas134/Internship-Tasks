package infrastructure

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateAccessToken(userID string, role string) (string, error)
	GenerateRefreshToken(userID string) (string, error)

	ValidateAccessToken(tokenString string) (userID string, role string, err error)
	ValidateRefreshToken(tokenString string) (userID string, err error)
}

type JWTServiceImpl struct {
	secretKey []byte
}

func NewJWTService() *JWTServiceImpl {

	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		panic("JWT_SECRET is not set")
	}

	return &JWTServiceImpl{
		secretKey: []byte(secret),
	}
}

type AccessClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`

	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID string `json:"user_id"`

	jwt.RegisteredClaims
}

func (service *JWTServiceImpl) GenerateAccessToken(
	userID string,
	role string,
) (string, error) {

	claims := AccessClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(15 * time.Minute),
			),
			IssuedAt: jwt.NewNumericDate(
				time.Now(),
			),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	tokenString, err := token.SignedString(
		service.secretKey,
	)

	if err != nil {
		return "", fmt.Errorf(
			"failed to sign access token: %w",
			err,
		)
	}

	return tokenString, nil
}

func (service *JWTServiceImpl) GenerateRefreshToken(
	userID string,
) (string, error) {

	claims := RefreshClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(7 * 24 * time.Hour),
			),
			IssuedAt: jwt.NewNumericDate(
				time.Now(),
			),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	tokenString, err := token.SignedString(
		service.secretKey,
	)

	if err != nil {
		return "", fmt.Errorf(
			"failed to sign refresh token: %w",
			err,
		)
	}

	return tokenString, nil
}

func (service *JWTServiceImpl) ValidateAccessToken(
	tokenString string,
) (string, string, error) {

	claims := &AccessClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New(
					"unexpected signing method",
				)
			}

			return service.secretKey, nil
		},
	)

	if err != nil {
		return "", "", fmt.Errorf(
			"invalid access token: %w",
			err,
		)
	}

	if !token.Valid {
		return "", "", errors.New(
		"invalid access token",
		)
	}

	return claims.UserID, claims.Role, nil
}

func (service *JWTServiceImpl) ValidateRefreshToken(
	tokenString string,
) (string, error) {

	claims := &RefreshClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New(
					"unexpected signing method",
				)
			}

			return service.secretKey, nil
		},
	)

	if err != nil {
		return "", fmt.Errorf(
			"invalid refresh token: %w",
			err,
		)
	}

	if !token.Valid {
		return "", errors.New(
			"invalid refresh token",
		)
	}

	return claims.UserID, nil
}

