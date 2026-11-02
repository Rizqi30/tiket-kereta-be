package routes

import (
	"tiket-kereta-be/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) *gin.Engine {
	// grouping API route
	api := r.Group("/api")
	{
		// sub group for auth
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
		}
	}

	return r
}