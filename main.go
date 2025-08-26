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
	database.DB.AutoMigrate(&models.Transaction{})

	// 2. Setup router
	r := gin.Default()
	routes.Routes(r)

	// 3. Run server
	r.Run(":8080") // listen on localhost:8080
}
