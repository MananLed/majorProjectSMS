package model

import(
	"github.com/google/uuid"
)

type Feedback struct{
	ID uuid.UUID  `json:"id"`
	ResidentID string `json:"resident_id"`
	Rating int32 `json:"rating"`
	Content string `json:"content"`
}