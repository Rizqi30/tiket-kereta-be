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

		// API List Stations (PUBLIC)
		api.GET("/stations", controllers.GetStations)
		api.GET("/stations/:id", controllers.GetStationByID)


		// ROLE ADMIN
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
		{
			// API CREATE, UPDATE, DELETE STATIONS
			admin.POST("/stations", controllers.CreateStation)
			admin.PUT("/stations/:id", controllers.UpdateStation)
			admin.DELETE("/stations/:id", controllers.DeleteStation)
		}
	}

	return r
}