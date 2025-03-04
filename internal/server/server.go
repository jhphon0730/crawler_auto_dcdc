package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func resgisterRoutes(r *gin.Engine) {
	api_v1 := r.Group("/api")
	{
		api_v1.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Hello World",
			})
		})
	}
}

func InitialServer() {
	r := gin.Default()
	// cors
	cors_cfg := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}
	r.Use(cors.New(cors_cfg))

	r.Run("0.0.0.0:8080")
}

