package usecase

import (
	"fmt"
	"github.com/LandGAA/authh2/internal/entity"
	"github.com/LandGAA/authh2/internal/repository"
	"github.com/LandGAA/authh2/pkg/jwt"
	"github.com/LandGAA/authh2/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var ErrorWrongPassword = fmt.Errorf("Неправильный пароль")

type UseCase interface {
	GetAllUsers() ([]entity.User, error)
	GetUserByID(id int) (entity.User, error)
	GetUserByEmail(email string) (entity.User, error)
	DeleteUser(id int) error
	CreateUser(user entity.User) error
	UpdatePassword(user entity.User) error
	HashPassword(user entity.User) (entity.User, error)
	CheckHashPassword(password string, user entity.User) bool
	ToDTO(users []entity.User) []entity.DTOUser
	Authenticate(email string, password string) (string, string, int64, error)
}

type UserUseCase struct {
	repo repository.Repository
}

func NewUserUseCase(repo repository.Repository) UseCase {
	return &UserUseCase{repo: repo}
}

func (u *UserUseCase) GetAllUsers() ([]entity.User, error) {
	return u.repo.GetAll()
}

func (u *UserUseCase) GetUserByID(id int) (entity.User, error) {
	return u.repo.GetByID(id)
}

func (u *UserUseCase) GetUserByEmail(email string) (entity.User, error) {
	return u.repo.GetByEmail(email)
}

func (u *UserUseCase) DeleteUser(id int) error {
	return u.repo.Delete(id)
}

func (u *UserUseCase) CreateUser(user entity.User) error {
	userWithHashPassword, err := u.HashPassword(user)
	if err != nil {
		logger.Logger.Error("Ошибка при хешировании пароля",
			zap.Error(err),
			zap.String("usecase", "CreateUser"))
		return err
	}
	userWithHashPassword.CreateAt = time.Now().String()
	return u.repo.Create(userWithHashPassword)
}

func (u *UserUseCase) UpdatePassword(user entity.User) error {
	userH, err := u.HashPassword(user)
	if err != nil {
		return err
	}
	return u.repo.UpdatePassword(userH)
}

func (u *UserUseCase) HashPassword(user entity.User) (entity.User, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return entity.User{}, err
	}
	user.Password = string(password)
	return user, nil
}

func (u *UserUseCase) CheckHashPassword(password string, user entity.User) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func (u *UserUseCase) ToDTO(users []entity.User) []entity.DTOUser {
	var dtoUsers = make([]entity.DTOUser, len(users))
	for i, user := range users {
		dtoUsers[i] = entity.DTOUser{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Role:     user.Role,
			CreateAt: user.CreateAt,
		}
	}
	return dtoUsers
}

func (uc *UserUseCase) Authenticate(email string, password string) (string, string, int64, error) {
	user, err := uc.repo.GetByEmail(email)
	if err != nil {
		return "", "", 0, err
	}

	if !uc.CheckHashPassword(password, user) {
		return "", "", 0, ErrorWrongPassword
	}

	accessToken, expiresIn, err := jwt.GenerateAccessToken(user)
	if err != nil {
		return "", "", 0, fmt.Errorf("ошибка генерации access токена: %w", err)
	}

	refreshToken, err := jwt.GenerateRefreshToken(user)
	if err != nil {
		return "", "", 0, fmt.Errorf("ошибка генерации refresh токена: %w", err)
	}

	return accessToken, refreshToken, expiresIn, nil
}
