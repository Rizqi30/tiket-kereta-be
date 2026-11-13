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

		// Stations (PUBLIC)
		api.GET("/stations", controllers.GetStations)
		api.GET("/stations/:id", controllers.GetStationByID)

		// Train (PUBLIC)
		api.GET("/trains", controllers.GetTrains)
		api.GET("/trains/:id", controllers.GetTrainByID)

		// Schedule (PUBLIC)
		api.GET("/schedules", controllers.SearchSchedules)
		api.GET("/schedules/:id", controllers.GetScheduleByID)
		api.GET("/schedules/:id/seats", controllers.GetScheduleSeats)

		// === ROLE ADMIN ===
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
		{
			// API CREATE, UPDATE, DELETE STATIONS
			admin.POST("/stations", controllers.CreateStation)
			admin.PUT("/stations/:id", controllers.UpdateStation)
			admin.DELETE("/stations/:id", controllers.DeleteStation)

			// TRAIN
			admin.POST("/trains", controllers.CreateTrain)
			admin.PUT("/trains/:id", controllers.UpdateTrain)
			admin.DELETE("/trains/:id", controllers.DeleteTrain)

			// CARRIAGE
			admin.POST("/trains/:id/carriages", controllers.CreateCarriage)
			admin.PUT("/carriages/:id", controllers.UpdateCarriage)
			admin.DELETE("/carriages/:id", controllers.DeleteCarriage)

			// Schedules
			admin.POST("/schedules", controllers.CreateSchedule)
			admin.PUT("/schedules/:id", controllers.UpdateSchedule)
			admin.DELETE("/schedules/:id", controllers.DeleteSchedule)
		}
	}

	return r
}