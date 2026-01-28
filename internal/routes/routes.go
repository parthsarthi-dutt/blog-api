package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/parthsarthi-dutt/blog-api/internal/handlers"
	"github.com/parthsarthi-dutt/blog-api/internal/middleware"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

auth := handlers.NewAuthHandler()

r.POST("/register", auth.Register)
r.POST("/login", auth.Login)

// Protected routes
protected := r.Group("/")
protected.Use(middleware.AuthMiddleware())

protected.GET("/test-protected", func(c *gin.Context) {
	email := c.GetString("email")
	c.JSON(200, gin.H{"message": "Hello " + email})
})
blog := handlers.NewTodoHandler()

protected.POST("/todos", blog.Create)
protected.GET("/todos", blog.List)
protected.PUT("/todos/:id", blog.Update)
protected.DELETE("/todos/:id", blog.Delete)

}
