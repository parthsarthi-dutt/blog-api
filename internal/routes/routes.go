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
blog := handlers.NewBlogHandler()
// Protected routes
r.GET("/blogs/public",blog.ListPublic)
protected := r.Group("/")
protected.Use(middleware.AuthMiddleware())

protected.GET("/test-protected", func(c *gin.Context) {
	email := c.GetString("email")
	c.JSON(200, gin.H{"message": "Hello " + email})
})


protected.POST("/blogs", blog.Create)
protected.GET("/blogs", blog.List)
protected.PUT("/blogs/:id", blog.Update)
protected.DELETE("/blogs/:id", blog.Delete)

}
