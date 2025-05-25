package main

import (
	"github.com/gin-gonic/gin"
	database "github.com/qslfrs/Go-API/db"
	"github.com/qslfrs/Go-API/models"
	"github.com/qslfrs/Go-API/routes"
)

func main() {
	// 1. Koneksi DB & auto-migrate
	database.Connect()
	database.DB.AutoMigrate(&models.User{})

	// 2. Setup router
	r := gin.Default()
	routes.UserRoutes(r)

	// 3. Run server
	r.Run(":8080") // listen on localhost:8080
}
