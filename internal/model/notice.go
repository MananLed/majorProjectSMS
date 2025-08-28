package model

import (
	"time"

	"github.com/google/uuid"
)

type Notice struct {
	ID         uuid.UUID `json:"id"`
	DateIssued time.Time `json:"date_issued"` 
	Content    string    `json:"content"`
	Month      time.Month `json:"month"`     
	Year       int        `json:"year"`     
}
