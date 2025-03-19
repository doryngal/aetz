package router

import (
	"binai.net/v2/config"
	"binai.net/v2/internal/handlers"
	"binai.net/v2/internal/repository"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config, db *repository.InitRepository) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			// Разрешить все источники
			return true
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	authHandler := handlers.NewAuthHandler(db.AuthRepo, cfg)
	lotHandler := handlers.NewLotHandler(db.LotRepo)

	auth := router.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/confirmation-code", authHandler.ConfirmationCode)
	}

	lotGroup := router.Group("/api/lots")
	{
		lotGroup.GET("/", lotHandler.GetLotList)
		lotGroup.GET("/:id", lotHandler.GetLotByID)
	}

	return router
}
