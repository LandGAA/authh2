package delivery

import (
	_ "github.com/LandGAA/authh2/docs"
	"github.com/LandGAA/authh2/internal/usecase"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"time"
)

// @title Auth SQLite
// @version 69
// @host localhost:8081
// @BasePath /v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func SetupRouters(u usecase.UseCase) *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	handler := NewUserHandler(u)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("v1")
	{
		api.GET("/users", handler.GetAll)
		api.GET("/users/:id", handler.GetByID)
		api.GET("/users/email/:email", handler.GetByEmail)

		api.POST("/login", handler.Login)
		api.POST("/register", handler.Register)

		auth := api.Use(AuthMiddleware())
		{
			auth.DELETE("/users/:id", handler.Delete)
			auth.PATCH("/users/:id", handler.UpdatePassword)
			auth.POST("/refresh", handler.Refresh)
		}
	}
	return r
}
