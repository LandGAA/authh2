package delivery

import (
	"fmt"
	"github.com/LandGAA/authh2/internal/entity"
	"github.com/LandGAA/authh2/internal/usecase"
	"github.com/LandGAA/authh2/pkg/jwt"
	"github.com/LandGAA/authh2/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserHandler struct {
	u usecase.UseCase
}

func NewUserHandler(usecase usecase.UseCase) *UserHandler {
	return &UserHandler{u: usecase}
}

// @Summary Получить всех пользователей
// @Description Получение всех пользователей (Без паролей)
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} entity.User
// @Failure 404 {string} string "Пустая база данных"
// @Router /users [get]
func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.u.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Таблица с пользователями пуста :("})
	}

	c.JSON(http.StatusOK, h.u.ToDTO(users))
}

// @Summary Получить пользователя по ID
// @Description Получишь пользователя по ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "Неправильный параметр"
// @Failure 404 {string} string "Пользователя нет"
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неправильный параметр"})
		return
	}

	user, err := h.u.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, fmt.Sprintf("Пользователь не найден: %w", err))
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Получить пользователя по Email
// @Description Получить пользователя по его Email
// @Accept json
// @Produce json
// @Param email path string true "Email пользователя"
// @Success 200 {object} entity.User
// @Failure 404 {sting} string "Пользователь не найден"
// @Router /users/email/{email} [get]
func (h *UserHandler) GetByEmail(c *gin.Context) {
	email := c.Param("email")
	user, err := h.u.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// @Summary Удалить пользователя по ID
// @Description Удаление пользователя по его ID
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {string} string "Пользователь удален"
// @Failure 400 {string} string "Неправильные параметры"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неправильный параметр"})
		return
	}

	err = h.u.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Пользователь с ID = %d удален", id))
}

// @Summary Регистрация пользователя
// @Description Введите данные пользователя
// @Accept json
// @Produce json
// @Success 200 {object} entity.User
// @Failure 400 {string} string "Некоректные данные"
// @Failure 500 {string} string "Пользователь уже создан или ошибка сервера"
// @Router /register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   fmt.Sprintf("Введены некоректные данные: %w", err),
			"details": user,
		})
		return
	}

	if err := h.u.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	logger.Logger.Info("Пользователь зарегистрировался")
	c.JSON(http.StatusOK, user)
}

// @Summary Обновление пароля
// @Description Введите новый пароль
// @Accept json
// @Produce json
// @Success 200 {object} entity.User
// @Failure 400 {string} string "Некоректные данные"
// @Failure 404 {string} string "Пользователь не найден"
// @Failure 500 {string} string "Пользователь уже создан или ошибка сервера"
// @Router /register [post]
func (h *UserHandler) UpdatePassword(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ResponsePassword struct {
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&ResponsePassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Неправильные данные %v", err)})
		return
	}

	user, err := h.u.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	user.Password = ResponsePassword.Password
	err = h.u.UpdatePassword(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Пароль успешно изменен")
}

// @Summary Вход
// @Description Вход email + пароль
// @Accept json
// @Produce json
// @Success 200 {object} entity.User
// @Failure 400 {string} string "Неправильные данные"
// @Failure 401 {string} string "Ошибка аутентификации, неверный email или пароль"
// @Failure 500 {string} string "Ошибка создания токенов"
// @Router /login [get]
func (h *UserHandler) Login(c *gin.Context) {
	var req jwt.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(":)))")
		return
	}

	user, err := h.u.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(h.u.CheckHashPassword(req.Password, user))

	if !h.u.CheckHashPassword(req.Password, user) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неправильный пароль"})
		return
	}

	accessToken, refreshToken, expiresIn, err := h.u.Authenticate(req.Email, req.Password)
	if err != nil {
		if err == usecase.ErrorWrongPassword {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Ошибка аутентификации",
				"details": "Неверный email или пароль",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, jwt.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		UserID:       user.ID,
		Role:         user.Role,
	})
}

func (h *UserHandler) Refresh(c *gin.Context) {
	var req jwt.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := jwt.ValidateToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Невалидный refresh токен"})
		return
	}

	user, err := h.u.GetUserByEmail(claims.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
		return
	}

	accessToken, expiresIn, err := jwt.GenerateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
		return
	}

	refreshToken, err := jwt.GenerateRefreshToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации refresh токена"})
		return
	}

	c.JSON(http.StatusOK, jwt.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	})
}
