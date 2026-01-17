package api

import (
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spotify-genre-organizer/backend/internal/api/handlers"
)

// Simple in-memory rate limiter
type rateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     int           // requests per window
	window   time.Duration // time window
}

type visitor struct {
	count    int
	lastSeen time.Time
}

func newRateLimiter(rate int, window time.Duration) *rateLimiter {
	rl := &rateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		window:   window,
	}
	go rl.cleanup()
	return rl
}

func (rl *rateLimiter) cleanup() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.window {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		rl.visitors[ip] = &visitor{count: 1, lastSeen: time.Now()}
		return true
	}

	if time.Since(v.lastSeen) > rl.window {
		v.count = 1
		v.lastSeen = time.Now()
		return true
	}

	if v.count >= rl.rate {
		return false
	}

	v.count++
	v.lastSeen = time.Now()
	return true
}

func rateLimitMiddleware(rl *rateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !rl.allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func SetupRoutes(r *gin.Engine) {
	// Rate limiting: 100 requests per minute per IP
	limiter := newRateLimiter(100, time.Minute)
	r.Use(rateLimitMiddleware(limiter))

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000" // Default to Vite dev server
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{frontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
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

		api.GET("/library/count", handlers.GetLibraryCount)

		api.GET("/settings", handlers.GetSettings)
		api.PUT("/settings", handlers.UpdateSettings)

		api.GET("/playlists", handlers.ListPlaylists)
		api.PATCH("/playlists/:id", handlers.UpdatePlaylist)
		api.DELETE("/playlists/:id", handlers.DeletePlaylist)
		api.POST("/playlists/:id/refresh", handlers.RefreshPlaylist)

		api.GET("/library/sync-status", handlers.GetSyncStatus)
		api.POST("/playlists/sync-all", handlers.SyncAllPlaylists)
	}
}
