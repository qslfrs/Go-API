package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	database "github.com/qslfrs/Go-API/db"
	"github.com/qslfrs/Go-API/models"
)

type TransactionResponse struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	StartDate   string `json:"startdate"`
	EndDate     string `json:"enddate"`
}

type DateResponse struct {
	Date        string                `json:"date"`
	Transaction []TransactionResponse `json:"transaction"`
}

type FinalResponse struct {
	Status  int            `json:"status"`
	Data    []DateResponse `json:"data"`
	Message string         `json:"message"`
}

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
	// startStr := startDate.Format("2006-01-02")
	// endStr := endDate.Format("2006-01-02")

	// result := database.DB.Where(
	// 	"(DATE(startdate) BETWEEN ? AND ?) OR (DATE(enddate) BETWEEN ? AND ?) OR (? BETWEEN DATE(startdate) AND DATE(enddate))",
	// 	startStr, endStr, startStr, endStr, startStr,
	// ).Find(&transaction)
	result := database.DB.Where("startdate <= ? AND enddate >= ?", endDate, startDate).Find(&transaction)

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
	// responseData := make(map[string][]map[string]interface{})

	// for _, t := range transaction {
	// 	// Rentang tanggal transaksi
	// 	for d := t.StartDate; !d.After(t.EndDate); d = d.AddDate(0, 0, 1) {
	// 		if d.Month() == startDate.Month() && d.Year() == startDate.Year() {
	// 			dateStr := d.Format("2006-01-02")
	// 			typeText := map[string]string{"A": "Attendance", "L": "Leave"}[t.Type]

	// 			responseData[dateStr] = append(responseData[dateStr], map[string]interface{}{
	// 				"name":        t.Name,
	// 				"type":        typeText,
	// 				"description": t.Description,
	// 				"startdate":   t.StartDate.Format("2006-01-02"),
	// 				"enddate":     t.EndDate.Format("2006-01-02"),
	// 			})
	// 		}
	// 	}
	// }

	//----------------------PAKE STRUCT-------------------------------- Lebih Rapi, Urutan JSONnya sesuai yang kita mao
	var responseData []DateResponse
	dateIndex := make(map[string]int)
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		responseData = append(responseData, DateResponse{
			Date:        dateStr,
			Transaction: []TransactionResponse{},
		})
		dateIndex[dateStr] = len(responseData) - 1
	}

	for _, t := range transaction {
		for d := t.StartDate; !d.After(t.EndDate); d = d.AddDate(0, 0, 1) {
			if d.Month() == startDate.Month() && d.Year() == startDate.Year() {
				dateStr := d.Format("2006-01-02")
				typeText := map[string]string{"A": "Attendance", "L": "Leave"}[t.Type]

				tr := TransactionResponse{
					Name:        t.Name,
					Type:        typeText,
					Description: t.Description,
					StartDate:   t.StartDate.Format("2006-01-02"),
					EndDate:     t.EndDate.Format("2006-01-02"),
				}

				idx := dateIndex[dateStr]
				responseData[idx].Transaction = append(responseData[idx].Transaction, tr)
			}
		}
	}

	finalResponse := FinalResponse{
		Status:  200,
		Data:    responseData,
		Message: "success",
	}

	c.JSON(http.StatusOK, finalResponse)

	// ---------------PAKE MAP--------- tidak dijamin sesuai dengan urutan, karena menggunakan urutan key dalam JSON bawaan Go (package encoding/json)
	// responseData := []map[string]interface{}{}
	// for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
	// 	dateStr := d.Format("2006-01-02")
	// 	responseData = append(responseData, map[string]interface{}{
	// 		"date":        dateStr,
	// 		"transaction": []map[string]interface{}{},
	// 	})
	// }

	// dateIndex := map[string]int{}
	// for i, item := range responseData {
	// 	dateIndex[item["date"].(string)] = i
	// }

	// for _, t := range transaction {
	// 	for d := t.StartDate; !d.After(t.EndDate); d = d.AddDate(0, 0, 1) {
	// 		if d.Month() == startDate.Month() && d.Year() == startDate.Year() {
	// 			dateStr := d.Format("2006-01-02")
	// 			typeText := map[string]string{"A": "Attendance", "L": "Leave"}[t.Type]

	// 			idx := dateIndex[dateStr]
	// 			transactions := responseData[idx]["transaction"].([]map[string]interface{})
	// 			transactions = append(transactions, map[string]interface{}{
	// 				"name":        t.Name,
	// 				"type":        typeText,
	// 				"description": t.Description,
	// 				"startdate":   t.StartDate.Format("2006-01-02"),
	// 				"enddate":     t.EndDate.Format("2006-01-02"),
	// 			})
	// 			responseData[idx]["transaction"] = transactions
	// 		}
	// 	}
	// }

	// c.JSON(http.StatusOK, gin.H{
	// 	"status":  200,
	// 	"data":    responseData,
	// 	"message": "success",
	// })
}
