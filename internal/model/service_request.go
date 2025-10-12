package model

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending    Status = "pending"
	StatusApproved   Status = "approved"
	StatusInProgress Status = "inProgress"
	StatusCompleted  Status = "completed"
	StatusCancelled  Status = "cancelled"
)

type ServiceType string

const (
	Electrician ServiceType = "electrician"
	Plumber     ServiceType = "plumber"
)

type ServiceRequest struct {
	RequestID   uuid.UUID   `json:"request_id"`
	ResidentID  string      `json:"resident_id"`
	Flat        string      `json:"flat"`
	Status      Status      `json:"status"`
	TimeSlot    string      `json:"time_slot"`
	StartTime   time.Time   `json:"start_time"`
	EndTime     time.Time   `json:"end_time"`
	ServiceType ServiceType `json:"service_type"`
	Date        string      `json:"date"`
	AssignedTo  string      `json:"assignedto"`
	FeedbackGiven bool      `json:"feedbackgiven"`
}
