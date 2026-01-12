package api

import (
	"github.com/gin-gonic/gin"
	"github.com/spotify-genre-organizer/backend/internal/api/handlers"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/health", handlers.HealthCheck)

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.GET("/login", handlers.NotImplemented)
			auth.GET("/callback", handlers.NotImplemented)
			auth.GET("/me", handlers.NotImplemented)
			auth.POST("/logout", handlers.NotImplemented)
		}

		api.POST("/organize", handlers.NotImplemented)
		api.GET("/organize/:id", handlers.NotImplemented)
	}
}
