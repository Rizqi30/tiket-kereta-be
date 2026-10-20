package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Inisialisasi Gin router default (sudah termasuk logger & recovery)
	r := gin.Default()

	// Pasang middleware CORS agar bisa diakses oleh frontend
	r.Use(cors.Default())

	// Endpoint pengecekan kesehatan server (opsional, untuk memastikan backend aktif)
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Tiket Kereta API Service",
			"status":  "active",
		})
	})

	// Jalankan server di port 3001
	r.Run(":3001")
}