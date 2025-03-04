package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/jhphon0730/crawler_auto_dcdc/internal/database"
)

func GetPosts(c *gin.Context) {
	// limit, page
	limit := c.DefaultQuery("limit", "10")
	page := c.DefaultQuery("page", "1")

	posts, err := database.LoadPostsByArray(limit, page)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
	"posts": posts,
	})
}

func resgisterRoutes(r *gin.Engine) {
	api_v1 := r.Group("/api")
	{
		api_v1.GET("/posts", GetPosts)
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

	// routes
	resgisterRoutes(r)

	r.Run("0.0.0.0:8080")
}

