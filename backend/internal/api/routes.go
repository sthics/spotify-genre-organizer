package api

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spotify-genre-organizer/backend/internal/api/handlers"
)

func SetupRoutes(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("FRONTEND_URL")},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.GET("/health", handlers.HealthCheck)

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.GET("/login", handlers.Login)
			auth.GET("/callback", handlers.Callback)
			auth.GET("/me", handlers.Me)
			auth.POST("/logout", handlers.Logout)
		}

		api.POST("/organize", handlers.StartOrganize)
		api.GET("/organize/:id", handlers.GetOrganizeStatus)
	}
}
