package jwt

import (
	"fmt"
	"github.com/LandGAA/authh2/internal/entity"
	"github.com/LandGAA/authh2/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"os"
	"time"
)

var SECRET_KEY []byte

type Claims struct {
	Email string `json:"email"`
	ID    int    `json:"id"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	UserID       int    `json:"user_id"`
	Role         string `json:"role"`
}

func Init() {
	err := godotenv.Load()
	if err != nil {
		logger.Logger.Fatal(fmt.Sprintf("Ошибка при загрузке файла .env: %v", err))
	}

	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		logger.Logger.Fatal("Переменная окружения JWT_SECRET_KEY не установлена")
	}
	SECRET_KEY = []byte(secret)
}

func GenerateAccessToken(user entity.User) (string, int64, error) {
	expirationTime := time.Now().Add(15 * time.Minute)

	claim := &Claims{
		Email: user.Email,
		ID:    user.ID,
		Role:  user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(SECRET_KEY)
	return tokenString, expirationTime.Unix(), err
}

func GenerateRefreshToken(user entity.User) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	claim := &Claims{
		Email: user.Email,
		ID:    user.ID,
		Role:  user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(SECRET_KEY)
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return SECRET_KEY, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("Невалидный токен")
	}

	return claims, nil
}
