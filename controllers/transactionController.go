package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	database "github.com/qslfrs/Go-API/db"
	"github.com/qslfrs/Go-API/models"
)

func GetTransaction(c *gin.Context) {
	month := c.Query("month")
	year := c.Query("year")

	if month == "" || year == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "month and year are required",
		})
		return
	}

	// Parsing bulan dan tahun
	monthInt, err := strconv.Atoi(month)
	if err != nil || monthInt < 1 || monthInt > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"status": 400, "message": "invalid month"})
		return
	}
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 400, "message": "invalid year"})
		return
	}

	// Range tanggal awal & akhir bulan
	startDate := time.Date(yearInt, time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	var transaction []models.Transaction
	startStr := startDate.Format("2006-01-02")
	endStr := endDate.Format("2006-01-02")

	result := database.DB.Where(
		"(DATE(startdate) BETWEEN ? AND ?) OR (DATE(enddate) BETWEEN ? AND ?) OR (? BETWEEN DATE(startdate) AND DATE(enddate))",
		startStr, endStr, startStr, endStr, startStr,
	).Find(&transaction)

	// database.DB = database.DB.Debug()

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 500, "message": result.Error.Error()})
		return
	}

	// For Dubugging
	// fmt.Println("StartDate:", startStr)
	// fmt.Println("EndDate:", endStr)
	// fmt.Println("RowsAffected:", result.RowsAffected)
	// fmt.Println("Error:", result.Error)
	// fmt.Println(database.DB.ToSQL(func(tx *gorm.DB) *gorm.DB {
	// 	return tx.Where("(DATE(startdate) BETWEEN ? AND ?) OR (DATE(enddate) BETWEEN ? AND ?) OR (? BETWEEN DATE(startdate) AND DATE(enddate))",
	// 		startStr, endStr, startStr, endStr, startStr)
	// }))

	if len(transaction) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": 404, "message": "no data found"})
		return
	}

	// Struktur data per tanggal
	responseData := make(map[string][]map[string]interface{})

	for _, t := range transaction {
		// Rentang tanggal transaksi
		for d := t.StartDate; !d.After(t.EndDate); d = d.AddDate(0, 0, 1) {
			if d.Month() == startDate.Month() && d.Year() == startDate.Year() {
				dateStr := d.Format("2006-01-02")
				typeText := map[string]string{"A": "Attendance", "L": "Leave"}[t.Type]

				responseData[dateStr] = append(responseData[dateStr], map[string]interface{}{
					"name":        t.Name,
					"type":        typeText,
					"description": t.Description,
					"startdate":   t.StartDate.Format("2006-01-02"),
					"enddate":     t.EndDate.Format("2006-01-02"),
				})
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"data":    responseData,
		"message": "success",
	})
}
