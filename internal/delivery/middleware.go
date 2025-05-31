package delivery

import (
	"fmt"
	"github.com/LandGAA/authh2/pkg/jwt"
	"github.com/LandGAA/authh2/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			logger.Logger.Error("Пустой хеддер",
				zap.String("middleware", "пустой хедер"))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Пустой херер"})
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		claims, err := jwt.ValidateToken(token)
		if err != nil {
			logger.Logger.Error(fmt.Sprintf("Невалидный токен %s", token),
				zap.Error(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Невалидный токен"})
			return
		}
		c.Set("email", claims.Email)
		c.Next()
	}
}
