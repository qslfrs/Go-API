package models

import (
	"time"
)

type Transaction struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"startdate" gorm:"column:startdate;type:date"`
	EndDate     time.Time `json:"enddate" gorm:"column:enddate;type:date"`
}

func (Transaction) TableName() string {
	return "transaction"
}
