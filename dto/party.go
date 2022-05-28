package dto

import (
	"time"

	"github.com/cridenour/go-postgis"
)

type Party struct {
	ID            string
	UserId        string
	Title         string
	IsPublic      bool
	Location      postgis.Point
	StreetAddress string
	PostalCode    string
	State         string
	Country       string
	StartDate     time.Time
	EndDate       time.Time
}
