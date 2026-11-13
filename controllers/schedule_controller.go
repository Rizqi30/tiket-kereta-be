package controllers

import (
	"net/http"
	"tiket-kereta-be/config"
	"tiket-kereta-be/models"

	"github.com/gin-gonic/gin"
)

// SearchSchedules - search the schedule based on departure station, destination, and date (PUBLIC)
func SearchSchedules(c *gin.Context) {
	departureStationID := c.Query("departure_station_id")
	arrivalStationID := c.Query("arrival_station_id")
	date := c.Query("date")

	if departureStationID == "" || arrivalStationID == "" || date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter departure_station_id, arrival_station_id, and date is required"})
		return
	}

	var schedules []models.Schedule
	config.DB.
		Preload("Train").
		Preload("DepartureStation").
		Preload("ArrivalStation").
		Where("departure_station_id = ? AND arrival_station_id = ? AND date = ?",	
			departureStationID, arrivalStationID, date).
		Find(&schedules)
	
	c.JSON(http.StatusOK, schedules)
}

// GetScheduleByID - Detail 1 schedule with the train info & station (PUBLIC)
func GetScheduleByID(c *gin.Context) {
	var schedule models.Schedule
	id := c.Param("id")

	if err := config.DB.
		Preload("Train.Carriages").
		Preload("DepartureStation").
		Preload("ArrivalStation").
		First(&schedule, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
			return
		}

	c.JSON(http.StatusOK, schedule)
}

// GetScheduleSeats - Viewing the carriage seating plan for a spesific schedule (PUBLIC)
func GetScheduleSeats(c *gin.Context) {
	scheduleID := c.Param("id")
	carriageID := c.Query("carriage_id")

	if carriageID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter carriage_id is required"})
		return
	}

	// Check is the schedule available
	var schedule models.Schedule
	if err := config.DB.First(&schedule, scheduleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule is not found"})
		return
	}

	// Take all the seats in that carriage
	var seats []models.Seat
	config.DB.Where("carriage_id = ?", carriageID).Find(&seats)

	// Check the seats already booked

	// All seats is available
	type SeatStatus struct {
		ID uint `json:"id"`
		SeatNumber string `json:"seat_number"`
		Status string `json:"status"` // available or booked
	}

	var seatStatuses []SeatStatus
	for _, seat := range seats {
		seatStatuses = append(seatStatuses, SeatStatus{
			ID: seat.ID,
			SeatNumber: seat.SeatNumber,
			Status: "AVAILABLE",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"schedule_id": schedule.ID,
		"carriage_id": carriageID,
		"seats": seatStatuses,
	})
}

// input struct for make the schedule
type ScheduleInput struct {
	TrainID uint `json:"train_id" binding:"required"`
	DepartureStationID uint `json:"departure_station_id" binding:"required"`
	ArrivalStationID uint `json:"arrival_station_id" binding:"required"`
	DepartureTime string `json:"departure_time" binding:"required"`
	ArrivalTime string `json:"arrival_time" binding:"required"`
	Date string `json:"date" binding:"required"`
	PriceEkonomi int `json:"price_ekonomi"`
	PriceEksekutif int `json:"price_eksekutif"` 
}

// CreateSchedule - Make new train schedule (ADMIN)
func CreateSchedule(c *gin.Context) {
	var input ScheduleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// validation: departure station and destination can't same
	if input.DepartureStationID == input.ArrivalStationID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Station departure and arrival can't same"})
		return
	}

	// check is the train available
	var train models.Train
	if err := config.DB.First(&train, input.TrainID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Train not found"})
		return
	}

	// check is departure station available
	var depStation models.Station
	if err := config.DB.First(&depStation, input.DepartureStationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Departure station not found"})
		return
	}

	// check is arrival station available
	var arrStation models.Station
	if err := config.DB.First(&arrStation, input.ArrivalStationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Arrival station not found"})
		return
	}

	schedule := models.Schedule{
		TrainID: input.TrainID,
		DepartureStationID: input.DepartureStationID,
		ArrivalStationID: input.ArrivalStationID,
		DepartureTime: input.DepartureTime,
		ArrivalTime: input.ArrivalTime,
		Date: input.Date,
		PriceEkonomi: input.PriceEkonomi,
		PriceEksekutif: input.PriceEksekutif,
	}

	// check error when insert to DB
	if err := config.DB.Create(&schedule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create schedule: " + err.Error()})
		return
	}

	// Preload relasi for response
	if err := config.DB.
		Preload("Train").
		Preload("DepartureStation").
		Preload("ArrivalStation").
		First(&schedule, schedule.ID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load schedule relations: " + err.Error()})
			return
		}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Schedule successfully created",
		"data": schedule,
	})
}

// UpdateSchedule - Edit schedule (ADMIN)
func UpdateSchedule(c *gin.Context) {
	var schedule models.Schedule
	id := c.Param("id")

	if err := config.DB.First(&schedule, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
		return
	}

	var input ScheduleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&schedule).Updates(models.Schedule{
		TrainID: input.TrainID,
		DepartureStationID: input.DepartureStationID,
		ArrivalStationID: input.ArrivalStationID,
		DepartureTime: input.DepartureTime,
		ArrivalTime: input.ArrivalTime,
		Date: input.Date,
		PriceEkonomi: input.PriceEkonomi,
		PriceEksekutif: input.PriceEksekutif,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Schedule successfully updated",
		"data": schedule,
	})
}

// DeleteSchedule - Delete schedule (ADMIN)
func DeleteSchedule(c *gin.Context) {
	var schedule models.Schedule
	id := c.Param("id")

	if err := config.DB.First(&schedule, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
		return
	}

	config.DB.Delete(&schedule)

	c.JSON(http.StatusOK, gin.H{
		"message": "Schedule successfully deleted",
	})
}