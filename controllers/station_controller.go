package controllers

import (
	"net/http"
	"tiket-kereta-be/config"
	"tiket-kereta-be/models"

	"github.com/gin-gonic/gin"
)

// GetStations - take all data of stations (public)
func GetStations(c *gin.Context) {
	var stations []models.Station
	config.DB.Find(&stations)

	c.JSON(http.StatusOK, stations)
}

// GetStationByID - take detail of station by ID (public)
func GetStationByID(c *gin.Context) {
	var station models.Station
	id := c.Param("id")

	if err := config.DB.First(&station, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Station not found"})
		return
	}

	c.JSON(http.StatusOK, station)
}

// input struct to making or editing station
type StationInput struct {
	Name string `json:"name" binding:"required"`
	Code string `json:"code" binding:"required"`
	City string `json:"city" binding:"required"`
}

// create station - adding new station (ADMIN)
func CreateStation(c *gin.Context) {
	var input StationInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check is station code already used
	var existing models.Station
	if config.DB.Where("code = ?", input.Code).First(&existing).Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Station code already used"})
		return
	}

	station := models.Station{
		Name: input.Name,
		Code: input.Code,
		City: input.City,
	}

	config.DB.Create(&station)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Station successfully added",
		"data": station,
	})
}

// update station - update station data (ADMIN)
func UpdateStation(c *gin.Context) {
	var station models.Station
	id := c.Param("id")

	if err := config.DB.First(&station, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Station not found"})
		return
	}

	var input StationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// update station data
	config.DB.Model(&station).Updates(models.Station{
		Name: input.Name,
		Code: input.Code,
		City: input.City,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Station successfully updated",
		"data": station,
	})
}

// delete station - delete station (ADMIN)
func DeleteStation(c *gin.Context) {
	var station models.Station
	id := c.Param("id")

	if err := config.DB.First(&station, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Station not found"})
		return
	}

	config.DB.Delete(&station)

	c.JSON(http.StatusOK, gin.H{
		"message": "Station successfully deleted",
	})
}