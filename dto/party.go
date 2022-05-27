package dto

import (
	"time"

	"github.com/twpayne/go-geom"
)

type Party struct {
	ID            string
	UserId        string
	Title         string
	IsPublic      bool
	Location      geom.Point
	StreetAddress string
	PostalCode    string
	State         string
	Country       string
	StartDate     time.Time
	EndDate       time.Time
}
