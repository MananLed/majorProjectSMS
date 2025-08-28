package model

import (
	"time"

	"github.com/google/uuid"
)

type Invoice struct {
	ID     uuid.UUID   `json:"id"`
	Amount float64     `json:"amount"`
	Month  time.Month  `json:"month"` 
	Year   int         `json:"year"`  
}
