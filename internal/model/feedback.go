package model

import(
	"github.com/google/uuid"
)

type Feedback struct{
	ID uuid.UUID  `json:"id"`
	ResidentID string `json:"resident_id"`
	Flat string `json:"flat"`
	Rating int32 `json:"rating"`
	Content string `json:"content"`
	ResidentName string `json:"name"`
	RequestID uuid.UUID `json:"request_id"`
	AssignedTo string `json:"assignedto"`
	ServiceType string `json:"servicetype"`
}