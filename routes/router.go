package routes

import (
	"tiket-kereta-be/controllers"
	"tiket-kereta-be/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) *gin.Engine {
	// grouping API route
	api := r.Group("/api")
	{
		// sub group for auth
		auth := api.Group("/auth")
		{
			// API LOGIN
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
			auth.GET("/me", middleware.AuthMiddleware(), controllers.GetMe)
		}
	}

	return r
}