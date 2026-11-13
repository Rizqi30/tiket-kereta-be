package controllers

import (
	"fmt"
	"net/http"
	"tiket-kereta-be/config"
	"tiket-kereta-be/models"

	"github.com/gin-gonic/gin"
)

// GetTrains - Take all the train (public)
func GetTrains(c *gin.Context) {
	var trains []models.Train
	config.DB.Find(&trains)

	c.JSON(http.StatusOK, trains)
}

// GetTrainByID - Detail 1 kereta with the carriage (public)
func GetTrainByID(c *gin.Context) {
	var train models.Train
	id := c.Param("id")

	// Preload: when take the train data also take the carriages data & seat
	if err := config.DB.Preload("Carriages").Preload("Carriages.Seats").First(&train, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Train is not found"})
		return
	}

	c.JSON(http.StatusOK, train)
}

// Input struct for create/edit train
type TrainInput struct {
	Name string `json:"name" binding:"required"`
	Code string `json:"code" binding:"required"`
}

// CreateTrain - add new train (ADMIN)
func CreateTrain(c *gin.Context) {
	var input TrainInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check duplicate train code
	var existing models.Train
	if config.DB.Where("code = ?", input.Code).First(&existing).Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Train code already used"})
		return
	}

	train := models.Train{
		Name: input.Name,
		Code: input.Code,
	}

	config.DB.Create(&train)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Train successfully added",
		"data": train,
	})
}

// UpdateTrain - Edit kereta (ADMIN)
func UpdateTrain(c *gin.Context) {
	var train models.Train
	id := c.Param("id")

	if err := config.DB.First(&train, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Train not found"})
		return
	}

	var input TrainInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&train).Updates(models.Train{
		Name: input.Name,
		Code: input.Code,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Train successfully updated",
		"data": train,
	})
}

// DeleteTrain - delete train with the carriage & seats (ADMIN)
func DeleteTrain(c *gin.Context) {
	var train models.Train
	id := c.Param("id")

	if err := config.DB.First(&train, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Train not found"})
		return
	}

	// delete all seats from carriages this train
	var carriages []models.Carriage
	config.DB.Where("train_id = ?", train.ID).Find(&carriages)
	for _, carriage := range carriages {
		config.DB.Where("carriage_id = ?", carriage.ID).Delete(&models.Seat{})
	}

	// delete all carriages belonging this train
	config.DB.Where("train_id = ?", train.ID).Delete(&models.Carriage{})

	// delete the train
	config.DB.Delete(&train)

	c.JSON(http.StatusOK, gin.H{
		"message": "Train with all carriage and seat successfully deleted",
	})
}


// BAGIAN GERBONG (CARRIAGE) + AUTO GENERATE KURSI
// input struct for the carriage
type CarriageInput struct {
	Name string `json:"name" binding:"required"`
	Class string `json:"class" binding:"required"`
	TotalSeats int `json:"total_seats" binding:"required"`
}

// CreateCarriage - add carriage to train + AUTO GENERATE SEAT (ADMIN)
func CreateCarriage(c *gin.Context) {
	trainID := c.Param("id")

	// check if train is there
	var train models.Train
	if err := config.DB.First(&train, trainID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Train is not found"})
		return
	}

	var input CarriageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// validation class
	validClasses := map[string]bool{"EKONOMI": true, "EKSEKUTIF": true}
	if !validClasses[input.Class] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Class must be EKONOMI or EKSEKUTIF"})
		return
	}

	// make carriage
	carriage := models.Carriage{
		TrainID: train.ID,
		Name: input.Name,
		Class: input.Class,
		TotalSeats: input.TotalSeats,
	}

	config.DB.Create(&carriage)

	// AUTO GENERATE SEAT
	// example: 50 seat -> 1A, 1B, 1C, 1D, 2A, ... sampai 13B
	columns := []string{"A", "B", "C", "D"}
	seatCount := 0
	row := 1

	for seatCount < input.TotalSeats {
		for _, col := range columns {
			if seatCount >= input.TotalSeats {
				break
			}
			seat := models.Seat{
				CarriageID: carriage.ID,
				SeatNumber: fmt.Sprintf("%d%s", row, col),
			}
			config.DB.Create(&seat)
			seatCount++
		}
		row++
	}

	// take the data of carriage with the seat of new generated
	config.DB.Preload("Seats").First(&carriage, carriage.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Carriage and seat successfully added",
		"data": carriage,
	})
}

// UpdateCarriage - edit the carriage (ADMIN)
func UpdateCarriage(c *gin.Context) {
	var carriage models.Carriage
	id := c.Param("id")

	if err := config.DB.First(&carriage, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carriage is not found"})
		return
	}

	var input struct {
		Name string `json:"name"`
		Class string `json:"class"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&carriage).Updates(models.Carriage{
		Name: input.Name,
		Class: input.Class,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Carriage successfully updated",
		"data": carriage,
	})
}

// DeleteCarriage - delete 1 carriage with the seats (ADMIN)
func DeleteCarriage(c *gin.Context) {
	var carriage models.Carriage
	id := c.Param("id")

	if err := config.DB.First(&carriage, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carriage is not found"})
		return
	}

	// delete all seats in this carriage first
	config.DB.Where("carriage_id = ?", carriage.ID).Delete(&models.Seat{})

	// delete the carriage
	config.DB.Delete(&carriage)

	c.JSON(http.StatusOK, gin.H{
		"message": "Carriage and all the seats successfully deleted",
	})
}