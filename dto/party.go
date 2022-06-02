package dto

import (
	"time"
)

type Party struct {
	ID            string
	UserId        string
	Title         string
	IsPublic      bool
	Lat           float32
	Long          float32
	StreetAddress string
	PostalCode    string
	State         string
	Country       string
	StartDate     time.Time
	EndDate       time.Time
}
